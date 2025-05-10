package call

import (
	"go.aew.app/api.v1/def"
)

var _ def.Caller = (*TraceCaller)(nil)

type TraceCaller struct {
	callerDefault
}

func NewTraceCaller(serialize def.Serialize, pool *def.MethodsPools) *TraceCaller {
	return &TraceCaller{callerDefault: callerDefault{
		serialize:  serialize,
		mIntercept: NewUserProxyInvokeImpl(methodInvokes),
		pool:       pool,
	}}
}

func (t *TraceCaller) Call(f *def.Entry, req *def.Request) interface{} {
	m := t.pool.FuncInfo(f.Fn)
	if len(m.Middleware) > 0 {
		for i := 0; i < len(m.Middleware); i++ {
			handle := m.Middleware[i]
			if v := handle(req.Request); v != nil {
				return v
			}
		}
	}
	return t.callerDefault.Call(f, req)
}
