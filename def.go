package api

import (
	"fmt"
	"gitee.com/fast_api/api/def"
	"gitee.com/fast_api/api/dwarf"
	"gitee.com/fast_api/api/log"
	"gitee.com/fast_api/api/mg"
	"os"
	"strings"
	"time"
)

var (
	prefix string
	maker  *dwarf.DwarfMaker
)

func trimPrefix(s string) string {
	if s != "" {
		return strings.ReplaceAll(s, prefix, "")
	}
	return s
}

func SetLogTrimPrefix(prefixM string) {
	prefix = prefixM
}

func PackApi() {
	PackApiWithPath(nil)
}

func SetExecPath(path *string) {
	if maker == nil {
		maker = dwarf.NewDwarfMaker()
		mg.Provide(func() *dwarf.DwarfMaker {
			return maker
		})
	}
	if dll, e := os.Executable(); e != nil && path == nil {
		maker.Init(&dll)
	} else {
		maker.Init(path)
	}
}

func PackApiWithPath(exePath func() *string) {
	start := time.Now()
	if exePath == nil {
		SetExecPath(nil)
	} else {
		SetExecPath(exePath())
	}
	fns := getFnCaches()
	log.Debugf("api had caches %d", len(fns))
	mg.Invoke(func(pool *def.MethodsPools) {
		for i, fn := range fns {
			findM := maker.LookFun(fn.Fn)
			if findM == nil {
				panic(fmt.Sprintf("not find %s in drawf", fn.Fn))
			}
			var args = make(map[string]dwarf.ArgsMeta)
			for _, arg := range findM.Args {
				args[arg.Name] = arg
			}
			pool.Set(findM.MethodName, &def.MethodInfo{
				Method:     fns[i],
				MethodName: findM.MethodName,
				Param:      args,
			})
			log.Infof("[%s] %s(%s) mapping url = %s", fn.Method, trimPrefix(findM.MethodName), printArgs(findM.Args), fn.Url)
		}
	})

	log.Infof("init use %s", time.Since(start))
}

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
