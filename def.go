package api

import (
	"math"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"time"

	"gitee.com/aifuturewell/methods"
	"gitee.com/fast_api/api/public"
	"github.com/sirupsen/logrus"
)

func initDef() {
	start := time.Now()
	if public.MethodsPools == nil {
		public.MethodsPools = make(public.MetaMethods)
	}
	fns := GetApi().getFnCaches()
	num := runtime.NumCPU()
	per := len(fns) / num
	mod := len(fns) % num
	maybe := int(math.Min(float64(len(fns)), float64(num)))
	var wait sync.WaitGroup
	var l sync.Mutex
	wait.Add(maybe)
	for i := 1; i <= maybe; i++ {
		go func(ii int) {
			defer wait.Done()
			rage := ii * per
			if ii == maybe && mod != 0 {
				rage += mod
			}
			for y := (ii - 1) * per; y < rage; y++ {
				v := reflect.ValueOf(fns[y].Fn)
				fName := runtime.FuncForPC(v.Pointer()).Name()
				med := methods.GetHelper().LookFun(fns[y].Fn)
				var args = make(map[string]methods.ArgsMeta)
				for _, arg := range med.Args {
					args[arg.Name] = arg
				}
				l.Lock()
				public.MethodsPools[med.MethodName] = public.MethodInfo{
					Pkg:        "",
					Receive:    "",
					Method:     fns[y],
					MethodName: med.MethodName,
					Param:      args,
				}
				l.Unlock()
				logrus.Infof("[%s] %s(%s) mapping url = %s", fns[y].Method, fName, printArgs(med.Args), fns[y].Url)
			}
		}(i)
	}
	wait.Wait()
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
