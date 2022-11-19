package dwarf

import (
	"debug/gosym"
	"regexp"
)

var includeRegex = []string{
	"^gitee.com/fast_api/.+$",
}

func AddIncludeRegex(pkg string) bool {
	includeRegex = append(includeRegex, pkg)
	return true
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
			return true
		}
	}
	return false
}
