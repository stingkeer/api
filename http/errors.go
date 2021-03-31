package http

import (
	"gitee.com/fast_api/api/public"
	"reflect"
)

type ErrorHandler func(err interface{}) interface{}

var errorsMap map[reflect.Type]ErrorHandler

func handleError(err interface{}) interface{} {
	if v, b := errorsMap[reflect.TypeOf(err)]; b {
		return v(err)
	}
	if v, b := err.(string); b {
		return public.NewError(v)
	}
	if v, b := err.(error); b {
		return public.NewError(v.Error())
	}
	return ""
}

func RegisterErrorHandler(p reflect.Type, handler ErrorHandler) {
	errorsMap[p] = handler
}
