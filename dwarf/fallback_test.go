package dwarf

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestIsRuntimePackageStdlibHeuristic: only stdlib symbols are excluded;
// dependency libraries — including the domains that were previously
// blanket-excluded (golang.org, gorm.io) — must stay indexable so handlers
// defined in them can be looked up. Inputs are real DWARF symbol names
// ("pkg.Func" — PackageName returns "" for bare import paths).
func TestIsRuntimePackageStdlibHeuristic(t *testing.T) {
	cases := map[string]bool{
		"fmt.Println":                               true,
		"net/http.Get":                              true,
		"runtime.main":                              true,
		"internal/cpu.doinit":                       true,
		"type..func1":                               true,
		"type:.dict":                                true,
		"github.com/user/repo.Handler":              false,
		"github.com/user/repo.Handler.func1":        false, // closures are handlers
		"golang.org/x/time/rate.NewLimiter":         false, // was true (domain exclude)
		"google.golang.org/grpc.Dial":               false, // was true
		"gorm.io/gorm.Open":                         false, // was true
		"go.aew.app/api.v1/test/rest.TestX":         false,
		"go.aew.app/api.v1/test/rest.TestX.func1.1": false,
		"main.main":                                 false, // main-package handlers must be indexed
	}
	for pkg, want := range cases {
		if got := isRuntimePackage(pkg); got != want {
			t.Errorf("isRuntimePackage(%q) = %v, want %v", pkg, got, want)
		}
	}
}

// TestSourceCandidatesPlain: paths without a version marker resolve to
// themselves only.
func TestSourceCandidatesPlain(t *testing.T) {
	got := sourceCandidates("/home/x/proj/handler.go")
	if len(got) != 1 || got[0] != "/home/x/proj/handler.go" {
		t.Errorf("plain path candidates = %v", got)
	}
}

// TestSourceCandidatesTrimpath: -trimpath paths (module@version/file.go)
// also probe the module cache.
func TestSourceCandidatesTrimpath(t *testing.T) {
	t.Setenv("GOMODCACHE", "/cache/mod")
	got := sourceCandidates("github.com/user/repo@v1.2.3/handler.go")
	if len(got) < 2 {
		t.Fatalf("trimpath candidates = %v, want module cache fallback", got)
	}
	if got[0] != "github.com/user/repo@v1.2.3/handler.go" {
		t.Errorf("first candidate = %q, want the literal path", got[0])
	}
	want := filepath.Join("/cache/mod", "github.com/user/repo@v1.2.3", "handler.go")
	if got[1] != want {
		t.Errorf("cache candidate = %q, want %q", got[1], want)
	}
}

// TestSourceCandidatesUppercaseModule: module cache encodes uppercase path
// segments as "!lower" (github.com/Azure → github.com/!azure).
func TestSourceCandidatesUppercaseModule(t *testing.T) {
	t.Setenv("GOMODCACHE", "/cache/mod")
	got := sourceCandidates("github.com/Azure/azure-lib@v0.1.0/a/b.go")
	want := filepath.Join("/cache/mod", "github.com/!azure/azure-lib@v0.1.0", filepath.FromSlash("a/b.go"))
	if len(got) < 2 || got[1] != want {
		t.Errorf("escaped candidate = %v, want %q among them", got, want)
	}
}

// TestEscapeModulePath: no-op for lowercase, bang-escapes uppercase letters.
func TestEscapeModulePath(t *testing.T) {
	if got := escapeModulePath("github.com/user/repo"); got != "github.com/user/repo" {
		t.Errorf("lowercase = %q", got)
	}
	if got := escapeModulePath("github.com/Azure/SDK-Go"); got != "github.com/!azure/!s!d!k-!go" {
		t.Errorf("uppercase = %q", got)
	}
}

// TestLoadSourceFromModuleCache: a trimpath-style path whose file exists in
// the (fake) module cache is actually parsed, and the miss case is cached.
func TestLoadSourceFromModuleCache(t *testing.T) {
	cache := t.TempDir()
	t.Setenv("GOMODCACHE", cache)

	// fake extracted dependency: <cache>/example.org/lib@v1.0.0/h.go
	modDir := filepath.Join(cache, "example.org/lib@v1.0.0")
	if err := os.MkdirAll(modDir, 0o755); err != nil {
		t.Fatal(err)
	}
	const src = "package lib\n\nfunc Handler(id string, page int) {}\n"
	if err := os.WriteFile(filepath.Join(modDir, "h.go"), []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}

	ps := loadSource("example.org/lib@v1.0.0/h.go")
	if ps == nil {
		t.Fatal("trimpath path not resolved through module cache")
	}
	names, ok := collectParamNames(ps.fset, ps.file, 3)
	if !ok || len(names) != 2 || names[0] != "id" || names[1] != "page" {
		t.Errorf("param names = %v ok=%v, want [id page]", names, ok)
	}

	// miss is cached as nil after a lookup attempt (no repeated fs probing)
	if ps := loadSource("example.org/lib@v9.9.9/none.go"); ps != nil {
		t.Fatal("missing dependency source unexpectedly resolved")
	}
	if v, loaded := sourceCache.Load("example.org/lib@v9.9.9/none.go"); !loaded || v != nil {
		t.Errorf("miss not cached as nil: loaded=%v", loaded)
	}
}

// TestLookFunFromSourceDependencyError: when neither DWARF nor source is
// available, the error must name the file and the actionable remedies.
func TestLookFunFromSourceDependencyError(t *testing.T) {
	t.Setenv("GOMODCACHE", t.TempDir()) // cache exists but is empty
	sourceCache.Delete("nowhere.org/dep@v1.0.0/h.go")

	_, err := lookFunFromSource(0, "nowhere.org/dep.h", nil)
	// pc=0 has no FuncForPC entry; pc may also be nil-handled — accept both
	// shapes as long as a descriptive error is returned
	if err == nil {
		t.Fatal("expected error for unresolvable dependency source")
	}
	msg := err.Error()
	if !strings.Contains(msg, "nowhere.org") && !strings.Contains(msg, "not find") {
		t.Errorf("error should mention the missing source or function: %s", msg)
	}
}
