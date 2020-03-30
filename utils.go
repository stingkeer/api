package api

import (
	"regexp"
	"strings"
)

/**
 * return (struct,func)
 */
func SplitFuncName(name string) (string, string) {
	matched, _ := regexp.MatchString(`\(.+\)`, name)
	if matched { //struct func
		sName := name
		var sFunc string
		if strings.ContainsAny(name, "-") {
			last := strings.LastIndex(name, "-")
			sName = name[:last]
			lastF := strings.LastIndex(sName, ".")
			sFunc = sName[lastF+1:]
		}
		re := regexp.MustCompile(`\(.+\)`)
		sStruct := re.FindString(sName)
		return strings.ReplaceAll(sStruct[1:len(sStruct)-1], "*", ""), sFunc
	} else {
		lastF := strings.LastIndex(name, ".")
		return "", name[lastF+1:]
	}

}
