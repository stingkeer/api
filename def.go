package api

import (
	"gitee.com/aifuturewell/methods"
	"gitee.com/fast_api/api/def"
	"gitee.com/fast_api/api/log"

	"math"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

func doMethod(start, end int, fns []*def.Entry) {
	for i := start; i < end; i++ {
		fn := fns[i]
		med := methods.GetHelper().LookFun(fn.Fn)
		var args = make(map[string]methods.ArgsMeta)
		for _, arg := range med.Args {
			args[arg.Name] = arg
		}
		def.GetMethodPools().Set(med.MethodName, &def.MethodInfo{
			Pkg:        "",
			Receive:    "",
			Method:     fns[i],
			MethodName: med.MethodName,
			Param:      args,
		})
		log.Infof("[%s] %s(%s) mapping url = %s", fn.Method, trimPrefix(med.MethodName), printArgs(med.Args), fn.Url)
	}
}

var _prefix string

func trimPrefix(s string) string {
	if s != "" {
		return strings.ReplaceAll(s, _prefix, "")
	}
	return s
}

func SetLogTrimPrefix(prefix string) {
	_prefix = prefix
}

func averageDo(cpu, number int, do func(start, end int, g *sync.WaitGroup)) {
	per := number / cpu
	mod := 0
	if per == 0 {
		per = 1
	} else {
		mod = number % cpu
	}
	maybe := int(math.Min(float64(number), float64(cpu)))
	var wg sync.WaitGroup
	wg.Add(maybe)
	for i := 1; i <= maybe; i++ {
		seg := i * per
		if i == maybe && mod != 0 {
			seg += mod
		}
		go do((i-1)*per, seg, &wg)
	}
	wg.Wait()
}

func PackApi() {
	PackApiWithPath(nil)
}

func SetExecPath(path *string) {
	if path == nil {
		s, e := os.Executable()
		if e == nil {
			methods.Init(s)
		}
	} else {
		methods.Init(*path)
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
	averageDo(runtime.NumCPU(), len(fns), func(start, end int, g *sync.WaitGroup) {
		doMethod(start, end, fns)
		g.Done()
	})
	log.Infof("init use %s", time.Since(start))
}

func printArgs(args []methods.ArgsMeta) string {
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
