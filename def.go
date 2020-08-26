package api

import (
	"gitee.com/aifuturewell/methods"
	"gitee.com/fast_api/api/public"
)

func initDef() {
	if public.MethodsPools == nil {
		public.MethodsPools = make(public.MetaMethods)
	}
	for _, entry := range GetApi().getFnCaches() {
		med := methods.GetHelper().LookFun(entry.Fn)
		var args = make(map[string]methods.ArgsMeta)
		for _, arg := range med.Args {
			args[arg.Name] = arg
		}
		public.MethodsPools[med.MethodName] = public.MethodInfo{
			Pkg:        "",
			Receive:    "",
			Method:     entry.Fn,
			MethodName: med.MethodName,
			Param:      args,
		}
	}
}
