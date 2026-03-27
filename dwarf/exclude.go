package dwarf

import (
	"debug/gosym"
	"strings"
)

var exclude = map[string]interface{}{
	"bufio":    nil,
	"fmt":      nil,
	"strconv":  nil,
	"os":       nil,
	"net":      nil,
	"errors":   nil,
	"http":     nil,
	"context":  nil,
	"io":       nil,
	"sync":     nil,
	"testing":  nil,
	"reflect":  nil,
	"regexp":   nil,
	"runtime":  nil,
	"syscall":  nil,
	"sort":     nil,
	"math":     nil,
	"internal": nil,
	"unicode":  nil,
	"text":     nil,
	"log":      nil,
	"hash":     nil,
	"flag":     nil,
	"html":     nil,
	"heap":     nil,
	"test":     nil,
	"go":       nil,
	"bytes":    nil,
	"time":     nil,
	"strings":  nil,
	"compress": nil,
	"encoding": nil,
	"debug":    nil,
	"path":     nil,
	"crypto":   nil,
	"embed":    nil,
	"database": nil,
	"slices":   nil,
	"maps":     nil,
	"mime":     nil,
	// "github.com":        nil,
	"vendor":            nil,
	"weak":              nil,
	"cmp":               nil,
	"unique":            nil,
	"container":         nil,
	"golang.org":        nil,
	"google.golang.org": nil,
	"gorm.io":           nil,
	"gvisor.dev":        nil,
}

func isRuntimePackage(pkg string) bool {
	if strings.HasPrefix(pkg, "type..") || strings.HasPrefix(pkg, "type:") {
		return true
	}
	sym := gosym.Sym{
		Name: pkg,
	}
	pkg = sym.PackageName()
	if index := strings.IndexAny(pkg, "/"); index > 0 {
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
