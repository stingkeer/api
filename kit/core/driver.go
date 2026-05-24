package core

import (
	"os"
	"strings"
	"sync"

	"go.aew.app/api.v1/def"
	"go.aew.app/api.v1/dwarf"
	"go.aew.app/api.v1/log"
)

var (
	maker = dwarf.NewDwarfMaker()
	once  sync.Once
)

func printArgs(args []dwarf.ArgsMeta) string {
	var s strings.Builder
	l := len(args) - 1
	for i, arg := range args {
		s.WriteString(arg.Name)
		if i != l {
			s.WriteString(",")
		}
	}
	return s.String()
}

func makerInit() {
	if dll := os.Getenv("API_DLL"); dll == "" {
		maker.Init(nil)
	} else {
		maker.Init(&dll)
	}
}
func _HttpM(method string, ctx *def.Context) def.HttpMethod {
	return func(f interface{}, url string) def.Option {
		return &option{url: url, method: method, mi: &def.MethodInfo{}}
	}
}

func HttpM(method string, ctx *def.Context) def.HttpMethod {

	if isTestMode() {
		return _HttpM(method, ctx)
	}

	//init dwarf
	once.Do(makerInit)
	return func(f interface{}, url string) def.Option {
		entry := &def.Entry{
			Url:        url,
			HttpMethod: method,
			Fn:         f,
		}
		ctx.Match.Add(url, entry)
		findM, err := maker.LookFun(entry.Fn)
		if err != nil {
			panic(err)
		}
		var args = make(map[string]dwarf.ArgsMeta)
		for _, arg := range findM.Args {
			args[arg.Name] = arg
		}
		methodInfo := &def.MethodInfo{
			Method:     entry,
			MethodName: findM.MethodName,
			Param:      args,
		}
		ctx.Pool.Set(findM.MethodName, methodInfo)
		log.Infof("[%s] %s(%s) mapping url = %s", entry.HttpMethod, findM.MethodName, printArgs(findM.Args), entry.Url)
		op := option{url: url, method: method}
		op.SetMethod(methodInfo)
		op.SetContext(ctx)
		return &op
	}
}

func isTestMode() bool {
	if os.Getenv("API_TEST") == "1" {
		return false
	}
	args := os.Args
	if len(args) == 0 {
		return false
	}
	if strings.HasSuffix(args[0], ".test") {
		for _, _args := range args[1:] {
			if strings.HasPrefix(_args, "-test.") {
				return true
			}
		}
	}
	return false
}
