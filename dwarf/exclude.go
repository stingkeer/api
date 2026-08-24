package dwarf

import (
	"debug/gosym"
	"strings"
)

// exclude lists stdlib top-level package names. A DWARF symbol whose package
// path starts with one of these is skipped when indexing — stdlib functions
// are never framework handlers, and skipping them keeps the index small.
//
// NOTE: whole-domain entries (golang.org, google.golang.org, gorm.io,
// gvisor.dev) were removed — they silently made handlers defined in those
// dependency libraries unfindable. Third-party modules must never be
// excluded: module names always carry a dot in the first path segment, but
// "main" and closure suffixes parse to dotless names too, so a dot-based
// stdlib heuristic would misclassify them.
var exclude = map[string]interface{}{
	"bufio":     nil,
	"fmt":       nil,
	"strconv":   nil,
	"os":        nil,
	"net":       nil,
	"errors":    nil,
	"http":      nil,
	"context":   nil,
	"io":        nil,
	"sync":      nil,
	"testing":   nil,
	"reflect":   nil,
	"regexp":    nil,
	"runtime":   nil,
	"syscall":   nil,
	"sort":      nil,
	"math":      nil,
	"internal":  nil,
	"unicode":   nil,
	"text":      nil,
	"log":       nil,
	"hash":      nil,
	"flag":      nil,
	"html":      nil,
	"heap":      nil,
	"test":      nil,
	"go":        nil,
	"bytes":     nil,
	"time":      nil,
	"strings":   nil,
	"compress":  nil,
	"encoding":  nil,
	"debug":     nil,
	"path":      nil,
	"crypto":    nil,
	"embed":     nil,
	"database":  nil,
	"slices":    nil,
	"maps":      nil,
	"mime":      nil,
	"weak":      nil,
	"cmp":       nil,
	"unique":    nil,
	"container": nil,
	"vendor":    nil,
	"archive":   nil,
}

func isRuntimePackage(pkg string) bool {
	if strings.HasPrefix(pkg, "type..") || strings.HasPrefix(pkg, "type:") {
		return true
	}
	sym := gosym.Sym{
		Name: pkg,
	}
	pkg = sym.PackageName()
	// Empty package names (closures like pkg.F.func1, or symbols gosym
	// cannot decompose) must be indexed: function literals are the most
	// common handler shape.
	if pkg == "" {
		return false
	}
	if index := strings.Index(pkg, "/"); index > 0 {
		if _, b := exclude[pkg[:index]]; b {
			return true
		}
	} else {
		if _, b := exclude[pkg]; b {
			return true
		}
	}
	return false
}
