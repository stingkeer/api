package call

import (
	"container/list"
	"reflect"
	"sync"
	"sync/atomic"

	"go.aew.app/api.v1/def"
)

type MethodInvoke func(fn MethodCaller, m *def.MethodInfo, args []reflect.Value) []reflect.Value

type MethodCaller interface {
	Invoke(m *def.MethodInfo, args []reflect.Value) []reflect.Value
}

var (
	methodInvokes = list.New()
	// methodInvokesVersion is bumped whenever the proxy list changes, so
	// UserProxyInvokeImpl can detect late registrations and rebuild its
	// cached chain instead of rebuilding it on every request.
	methodInvokesVersion atomic.Int64
)

func SetMethodProxy(invoke MethodInvoke) {
	methodInvokes.PushFront(invoke)
	methodInvokesVersion.Add(1)
}

type UserProxyInvokeImpl struct {
	mu      sync.RWMutex
	list    *list.List
	head    MethodCaller // cached composed chain
	version int64        // list version when head was built
}

func NewUserProxyInvokeImpl(list *list.List) *UserProxyInvokeImpl {
	u := &UserProxyInvokeImpl{list: list}
	var fn MethodInvoke = func(fn MethodCaller, m *def.MethodInfo, args []reflect.Value) []reflect.Value {
		return reflect.ValueOf(m.Method.Fn).Call(args)
	}
	list.PushBack(fn)
	methodInvokesVersion.Add(1)
	return u
}

// chain returns the composed invocation chain, rebuilding it only when the
// proxy list changed since the last build. Previously every request walked
// the list and allocated a fresh helper per proxy — pure per-request garbage.
func (d *UserProxyInvokeImpl) chain() MethodCaller {
	if ver := methodInvokesVersion.Load(); ver == d.version && d.head != nil {
		return d.head
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	if ver := methodInvokesVersion.Load(); ver == d.version && d.head != nil {
		return d.head
	}
	var chain MethodCaller
	for e := d.list.Back(); e != nil; e = e.Prev() {
		chain = &methodCallerHelper{d: e.Value.(MethodInvoke), super: chain}
	}
	d.head = chain
	d.version = methodInvokesVersion.Load()
	return chain
}

// Invoke runs the request through the composed proxy chain. The front-most
// registered proxy is the outermost wrapper; the real method call is the
// innermost link (pushed by NewUserProxyInvokeImpl).
func (d *UserProxyInvokeImpl) Invoke(m *def.MethodInfo, args []reflect.Value) []reflect.Value {
	return d.chain().Invoke(m, args)
}

type methodCallerHelper struct {
	d     MethodInvoke
	super MethodCaller
}

func (mp *methodCallerHelper) Invoke(m *def.MethodInfo, args []reflect.Value) []reflect.Value {
	return mp.d(mp.super, m, args)
}
