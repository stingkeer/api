package call

import (
	"container/list"
	"reflect"

	"go.aew.app/api.v1/def"
)

type MethodInvoke func(fn MethodCaller, m *def.MethodInfo, args []reflect.Value) []reflect.Value

type MethodCaller interface {
	Invoke(m *def.MethodInfo, args []reflect.Value) []reflect.Value
}

var methodInvokes = list.New()

func SetMethodProxy(invoke MethodInvoke) {
	methodInvokes.PushFront(invoke)
}

type UserProxyInvokeImpl struct {
	list *list.List
}

func NewUserProxyInvokeImpl(list *list.List) *UserProxyInvokeImpl {
	u := &UserProxyInvokeImpl{list: list}
	var fn MethodInvoke = func(fn MethodCaller, m *def.MethodInfo, args []reflect.Value) []reflect.Value {
		return reflect.ValueOf(m.Method.Fn).Call(args)
	}
	list.PushBack(fn)
	return u
}

func (d *UserProxyInvokeImpl) Invoke(m *def.MethodInfo, args []reflect.Value) []reflect.Value {
	var temp *methodCallerHelper
	for e := d.list.Back(); e != nil; e = e.Prev() {
		temp = newMethodCallerHelper(e.Value.(MethodInvoke), temp)
	}
	var begin MethodInvoke = func(fn MethodCaller, m *def.MethodInfo, args []reflect.Value) []reflect.Value {
		return fn.Invoke(m, args)
	}
	values := begin(temp, m, args)
	return values
}

type methodCallerHelper struct {
	d     MethodInvoke
	super MethodCaller
}

func newMethodCallerHelper(d MethodInvoke, super MethodCaller) *methodCallerHelper {
	return &methodCallerHelper{d: d, super: super}
}

func (mp *methodCallerHelper) Invoke(m *def.MethodInfo, args []reflect.Value) []reflect.Value {
	return mp.d(mp.super, m, args)
}
