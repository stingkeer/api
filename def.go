package api

import (
	"fmt"
	"gitee.com/fast_api/api/def"
	"gitee.com/fast_api/api/dwarf"
	"gitee.com/fast_api/api/log"

	"math"
	"os"
	"strings"
	"sync"
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
	if maker == nil {
		maker = dwarf.NewDwarfMaker()
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
	for i, fn := range fns {
		findM := maker.LookFun(fn.Fn)
		if findM == nil {
			panic(fmt.Sprintf("not find %s in drawf", fn.Fn))
		}
		var args = make(map[string]dwarf.ArgsMeta)
		for _, arg := range findM.Args {
			args[arg.Name] = arg
		}
		def.GetMethodPools().Set(findM.MethodName, &def.MethodInfo{
			Pkg:        "",
			Receive:    "",
			Method:     fns[i],
			MethodName: findM.MethodName,
			Param:      args,
		})
		log.Infof("[%s] %s(%s) mapping url = %s", fn.Method, trimPrefix(findM.MethodName), printArgs(findM.Args), fn.Url)
	}
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
