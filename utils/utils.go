package utils

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"reflect"
	"regexp"
	"strings"

	"go.aew.app/api.v1/log"
)

// SplitFuncName
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

func Md5String(s string) string {
	h := md5.New()
	io.WriteString(h, s)
	return hex.EncodeToString(h.Sum(nil))
}

// DefaultCallValue
// other param set default value
func DefaultCallValue(typ reflect.Type) reflect.Value {
	switch typ.Kind() {
	case reflect.Bool:
		return reflect.ValueOf(false)
	case reflect.String:
		return reflect.ValueOf("")
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return reflect.ValueOf(0).Convert(typ)
	default:
		log.Errorf("DefaultCallValue error kind %s", typ.Kind())
	}
	return reflect.ValueOf(nil)
}
