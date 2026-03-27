package dwarf

import (
	"debug/dwarf"
	"debug/elf"
	"debug/macho"
	"debug/pe"
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"runtime"
	"time"
)

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

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
	openData func() *dwarf.Reader
	r        *dwarf.Reader
	debug    map[string]Params
	usedMode FilterMode
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

func (h *DwarfMaker) load(exe *string) {
	path := ""
	if exe != nil {
		path = *exe
	} else {
		if dll, err := os.Executable(); err == nil {
			path = dll
		}
	}
	if path == "" {
		panic("DwarfMaker load path == nil")
	}
	fmt.Printf("system [%s/%s] %s\n", runtime.GOOS, runtime.GOARCH, path)
	switch runtime.GOOS {
	case "windows":
		h.openData = func() *dwarf.Reader {
			f, e := pe.Open(path)
			checkError(e)
			defer f.Close()
			data, e := f.DWARF()
			checkError(e)
			return data.Reader()
		}
	case "linux", "android":
		h.openData = func() *dwarf.Reader {
			f, e := elf.Open(path)
			checkError(e)
			defer f.Close()
			data, e := f.DWARF()
			checkError(e)
			return data.Reader()
		}

	case "darwin":
		h.openData = func() *dwarf.Reader {
			f, e := macho.Open(path)
			checkError(e)
			defer f.Close()
			data, e := f.DWARF()
			checkError(e)
			return data.Reader()
		}
	default:
		panic(fmt.Sprintf("not support %s", runtime.GOOS))
	}
}

func (h *DwarfMaker) Init(exe *string) {
	now := time.Now()
	h.load(exe)
	h.r = h.openData()
	tempName := ""
	for r, _ := h.r.Next(); r != nil; r, _ = h.r.Next() {
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
	log.Printf("DwarfMaker init use %s len %d", time.Since(now), len(h.debug))
}

func (h *DwarfMaker) LookFun(inf interface{}) (*MethodMeta, error) {
	v := reflect.ValueOf(inf)
	fName := runtime.FuncForPC(v.Pointer()).Name()
	if v.Kind() != reflect.Func {
		return nil, errors.New("no func type")
	}
	tp := v.Type()
	//argsNum := tp.NumIn()
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

func (h *DwarfMaker) SetFilterMode(mode FilterMode) {
	h.usedMode = mode
}
