package api

import (
	"fmt"
	"gitee.com/aifuturewell/methods"
	"gitee.com/fast_api/api/public"
	"github.com/sirupsen/logrus"
	"math"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

func doMethod(start, end int, fns []*public.Entry) {
	for i := start; i < end; i++ {
		fn := fns[i]
		med := methods.GetHelper().LookFun(fn.Fn)
		fmt.Println(med)
		var args = make(map[string]methods.ArgsMeta)
		for _, arg := range med.Args {
			args[arg.Name] = arg
		}
		public.MethodsPools[med.MethodName] = public.MethodInfo{
			Pkg:        "",
			Receive:    "",
			Method:     fns[i],
			MethodName: med.MethodName,
			Param:      args,
		}
		logrus.Infof("[%s] %s(%s) mapping url = %s", fn.Method, med.MethodName, printArgs(med.Args), fn.Url)
	}
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

func PackApiWithPath(exePath func() string) {
	start := time.Now()
	//init methods
	var path string
	if exePath == nil {
		path, _ = os.Executable()
	} else {
		path = exePath()
	}
	methods.Init(path)
	if public.MethodsPools == nil {
		public.MethodsPools = make(public.MetaMethods)
	}
	fns := GetApi().getFnCaches()
	logrus.Debugf("api had caches %d", len(fns))
	averageDo(runtime.NumCPU(), len(fns), func(start, end int, g *sync.WaitGroup) {
		doMethod(start, end, fns)
		g.Done()
	})
	logrus.Infof("init use %s", time.Since(start))
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
