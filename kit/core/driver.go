package core

import (
	"os"
	"strings"
	"sync"

	"gitee.com/fast_api/api/def"
	"gitee.com/fast_api/api/dwarf"
	"gitee.com/fast_api/api/log"
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

func HttpM(method string, ctx *def.Context) def.HttpMethod {
	//init dwarf
	once.Do(makerInit)

	return func(f interface{}, url string) *def.Option {
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
		if err != nil {
			panic(err)
		}
		log.Infof("[%s] %s(%s) mapping url = %s", entry.HttpMethod, findM.MethodName, printArgs(findM.Args), entry.Url)
		op := def.Option{}
		op.SetMethod(methodInfo)
		op.SetContext(ctx)
		return &op
	}
}
