package dwarf

import (
	"debug/gosym"
	"regexp"
)

var includeRegex = []string{
	"^go.aew.app/.+$",
}

func isInclude(pkg string) bool {
	sym := gosym.Sym{
		Name: pkg,
	}
	pkg = sym.PackageName()
	for _, regex := range includeRegex {
		match, err := regexp.MatchString(regex, pkg)
		if err != nil {
			panic(err)
		}
		if match {
			return false
		}
	}
	return true
}
