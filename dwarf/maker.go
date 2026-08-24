package dwarf

import (
	"debug/dwarf"
	"debug/elf"
	"debug/macho"
	"debug/pe"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"time"
)

type (
	Mode       int
	Params     []string
	FilterMode struct {
		mode Mode
		fn   func(pkg string) bool
	}
)

func (f *FilterMode) Mode() Mode {
	return f.mode
}

const (
	RuntimePackageMode Mode = 1
	IncludeMode        Mode = 2
)

var (
	RuntimeExclude = FilterMode{fn: isRuntimePackage, mode: RuntimePackageMode}
	SelfInclude    = FilterMode{fn: isInclude, mode: IncludeMode}
)

type DwarfMaker struct {
	openReader func() (*dwarf.Reader, error)
	r          *dwarf.Reader
	debug      map[string]Params
	usedMode   FilterMode
	// noDwarf is set when the executable carries no usable DWARF section
	// (e.g. Go 1.24+ test binaries, which strip it by default). LookFun then
	// falls back to parsing parameter names from the Go source.
	noDwarf bool
}

func (h *DwarfMaker) UsedMode() FilterMode {
	return h.usedMode
}

func NewDwarfMakerWithMode(mode FilterMode) *DwarfMaker {
	return &DwarfMaker{debug: make(map[string]Params, 1000), usedMode: mode}
}

func NewDwarfMaker() *DwarfMaker {
	return NewDwarfMakerWithMode(RuntimeExclude)
}

func (h *DwarfMaker) AddExclude(pkg string) bool {
	if _, b := exclude[pkg]; b {
		return false
	}
	exclude[pkg] = nil
	return true
}

func (h *DwarfMaker) AddIncludeRegex(pkg string) bool {
	includeRegex = append(includeRegex, pkg)
	return true
}

func (h *DwarfMaker) load(exe *string) error {
	path := ""
	if exe != nil {
		path = *exe
	} else {
		if dll, err := os.Executable(); err == nil {
			path = dll
		}
	}
	if path == "" {
		return errors.New("DwarfMaker load path == nil")
	}
	fmt.Printf("system [%s/%s] %s\n", runtime.GOOS, runtime.GOARCH, path)
	openFile := func() (*dwarf.Data, error) {
		switch runtime.GOOS {
		case "windows":
			f, e := pe.Open(path)
			if e != nil {
				return nil, e
			}
			defer f.Close()
			return f.DWARF()
		case "linux", "android":
			f, e := elf.Open(path)
			if e != nil {
				return nil, e
			}
			defer f.Close()
			return f.DWARF()
		case "darwin":
			f, e := macho.Open(path)
			if e != nil {
				return nil, e
			}
			defer f.Close()
			return f.DWARF()
		default:
			return nil, fmt.Errorf("not support %s", runtime.GOOS)
		}
	}
	h.openReader = func() (*dwarf.Reader, error) {
		data, err := openFile()
		if err != nil {
			return nil, err
		}
		return data.Reader(), nil
	}
	return nil
}

func (h *DwarfMaker) Init(exe *string) {
	now := time.Now()
	if err := h.load(exe); err != nil {
		h.noDwarf = true
		log.Printf("DwarfMaker: DWARF unavailable (%v), falling back to source parsing", err)
		return
	}
	reader, err := h.openReader()
	if err != nil {
		h.noDwarf = true
		log.Printf("DwarfMaker: DWARF unavailable (%v), falling back to source parsing", err)
		return
	}
	h.r = reader
	tempName := ""
	for r, rerr := h.r.Next(); r != nil; r, rerr = h.r.Next() {
		if rerr != nil {
			break
		}
		if rName := r.Val(dwarf.AttrName); r.Tag == dwarf.TagSubprogram && rName != nil {
			tempName = rName.(string)
			if h.usedMode.fn(tempName) {
				continue
			}
			h.debug[tempName] = Params{}
		}
		if r.Tag != dwarf.TagFormalParameter {
			continue
		}
		if v, b := r.Val(dwarf.AttrVarParam).(bool); b && v {
			continue
		}
		n := r.Val(dwarf.AttrName)
		if n == nil {
			continue
		}
		if _, b := h.debug[tempName]; b {
			h.debug[tempName] = append(h.debug[tempName], n.(string))
		}
	}
	h.r = nil
	if len(h.debug) == 0 {
		h.noDwarf = true
		log.Printf("DwarfMaker: no usable DWARF entries, falling back to source parsing")
		return
	}
	log.Printf("DwarfMaker init use %s len %d", time.Since(now), len(h.debug))
}

func (h *DwarfMaker) LookFun(inf interface{}) (*MethodMeta, error) {
	v := reflect.ValueOf(inf)
	fName := runtime.FuncForPC(v.Pointer()).Name()
	if v.Kind() != reflect.Func {
		return nil, errors.New("no func type")
	}
	tp := v.Type()
	if h.noDwarf || len(h.debug) == 0 {
		return lookFunFromSource(v.Pointer(), fName, tp)
	}
	if params, b := h.debug[fName]; b {
		var args []ArgsMeta
		for i, vName := range params {
			args = append(args, ArgsMeta{
				Order: i,
				Name:  vName,
				Typ:   tp.In(i),
			})
		}
		return &MethodMeta{
			MethodName: fName,
			Args:       args,
		}, nil
	}
	return nil, fmt.Errorf("not find %s in drawf len %d", fName, len(h.debug))
}

// lookFunFromSource resolves parameter names by parsing the function's
// source file with go/ast. It is the fallback for executables that carry no
// DWARF section — notably test binaries built by Go 1.24+, where DWARF is
// stripped by default.
func lookFunFromSource(pc uintptr, fName string, tp reflect.Type) (*MethodMeta, error) {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return nil, fmt.Errorf("not find %s in runtime", fName)
	}
	file, line := fn.FileLine(pc)
	if file == "" {
		return nil, fmt.Errorf("no source location for %s (binary stripped of pclntab?)", fName)
	}
	src := loadSource(file)
	if src == nil {
		return nil, fmt.Errorf(
			"source %s for %s is not readable on this machine and the binary carries no DWARF "+
				"(common for Go 1.24+ test binaries or -trimpath builds). Handlers defined in "+
				"dependency libraries resolve through the module cache: run where sources "+
				"are available, or rebuild with -ldflags=-w=false, or set API_DLL",
			file, fName)
	}
	names, ok := collectParamNames(src.fset, src.file, line)
	if !ok {
		return nil, fmt.Errorf("not find %s at %s:%d", fName, file, line)
	}
	var args []ArgsMeta
	for i, vName := range names {
		if i >= tp.NumIn() {
			break
		}
		args = append(args, ArgsMeta{
			Order: i,
			Name:  vName,
			Typ:   tp.In(i),
		})
	}
	return &MethodMeta{
		MethodName: fName,
		Args:       args,
	}, nil
}

// parsedSource is a cached parse result for one source file.
type parsedSource struct {
	fset *token.FileSet
	file *ast.File
}

// sourceCache caches parsed files (and misses, as nil) across LookFun calls:
// route registration looks up many handlers in the same files, and module
// cache files are identical for the life of the process.
var sourceCache sync.Map // string -> *parsedSource (nil = known missing)

func loadSource(file string) *parsedSource {
	if v, ok := sourceCache.Load(file); ok {
		src, _ := v.(*parsedSource)
		return src
	}
	src := parseSourceBestEffort(file)
	if src != nil {
		sourceCache.Store(file, src)
	} else {
		// store untyped nil — a typed (*parsedSource)(nil) would compare
		// non-nil when loaded back into an interface
		sourceCache.Store(file, nil)
	}
	return src
}

func parseSourceBestEffort(file string) *parsedSource {
	for _, candidate := range sourceCandidates(file) {
		if _, err := os.Stat(candidate); err != nil {
			continue
		}
		fset := token.NewFileSet()
		astFile, err := parser.ParseFile(fset, candidate, nil, 0)
		if err != nil {
			continue
		}
		return &parsedSource{fset: fset, file: astFile}
	}
	return nil
}

// sourceCandidates returns every on-disk location a recorded source path may
// resolve to. The literal path covers dev machines (absolute paths into the
// module cache: /home/x/go/pkg/mod/github.com/lib@v1.0.0/h.go). For
// -trimpath builds the recorded path is module-relative
// (github.com/lib@v1.0.0/h.go), so the module cache is derived explicitly.
func sourceCandidates(file string) []string {
	out := []string{file}
	i := strings.IndexByte(file, '@')
	if i <= 0 {
		return out
	}
	rest := file[i+1:]
	j := strings.IndexByte(rest, '/')
	if j <= 0 {
		return out
	}
	mod, version, rel := file[:i], rest[:j], rest[j+1:]
	if cache := moduleCacheDir(); cache != "" {
		base := filepath.Join(cache, escapeModulePath(mod)+"@"+version)
		out = append(out,
			filepath.Join(base, filepath.FromSlash(rel)),
			// unescaped variant, in case the cache layout differs
			filepath.Join(cache, mod+"@"+version, filepath.FromSlash(rel)),
		)
	}
	return out
}

// moduleCacheDir locates the Go module cache: GOMODCACHE if set, otherwise
// the default <user cache dir>/go/pkg/mod.
func moduleCacheDir() string {
	if v := os.Getenv("GOMODCACHE"); v != "" {
		return v
	}
	if dir, err := os.UserCacheDir(); err == nil {
		return filepath.Join(dir, "go", "pkg", "mod")
	}
	return ""
}

// escapeModulePath applies the module cache's case-encoding: uppercase
// letters in module paths are stored as "!x" (github.com/Azure →
// github.com/!azure).
func escapeModulePath(mod string) string {
	if !strings.ContainsFunc(mod, func(c rune) bool { return c >= 'A' && c <= 'Z' }) {
		return mod
	}
	var b strings.Builder
	b.Grow(len(mod))
	for i := 0; i < len(mod); i++ {
		if c := mod[i]; c >= 'A' && c <= 'Z' {
			b.WriteByte('!')
			b.WriteByte(c + ('a' - 'A'))
		} else {
			b.WriteByte(c)
		}
	}
	return b.String()
}

// collectParamNames finds the func declaration or literal whose `func` keyword
// starts at the given line and returns its parameter names in order. If no
// exact line match exists, it accepts the innermost function whose span covers
// the line (handles multi-line signatures).
func collectParamNames(fset *token.FileSet, f *ast.File, line int) ([]string, bool) {
	var (
		exact  []string
		inner  []string
		innerW = 1 << 30
		have   bool
	)
	ast.Inspect(f, func(n ast.Node) bool {
		var pl *ast.FieldList
		switch x := n.(type) {
		case *ast.FuncDecl:
			pl = x.Type.Params
		case *ast.FuncLit:
			pl = x.Type.Params
		default:
			return true
		}
		if pl == nil {
			return true
		}
		start := fset.Position(n.Pos()).Line
		end := fset.Position(n.End()).Line
		if start == line {
			if exact == nil {
				exact = fieldNames(pl)
				have = true
			}
			return false
		}
		if start <= line && line <= end {
			if w := end - start; w < innerW {
				innerW = w
				inner = fieldNames(pl)
				have = true
			}
		}
		return true
	})
	if !have {
		return nil, false
	}
	if exact != nil {
		return exact, true
	}
	return inner, true
}

func fieldNames(pl *ast.FieldList) []string {
	var names []string
	for _, field := range pl.List {
		for _, name := range field.Names {
			names = append(names, name.Name)
		}
	}
	return names
}

func (h *DwarfMaker) SetFilterMode(mode FilterMode) {
	h.usedMode = mode
}
