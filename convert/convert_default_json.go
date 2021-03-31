package convert

import (
	"encoding/json"
	"gitee.com/fast_api/api/public"
	"net/http"
	"reflect"
)

type JsonConvertImpl struct{}

func (c *JsonConvertImpl) Decode(bytes []byte, vpr interface{}) error {
	return json.Unmarshal(bytes, vpr)
}

func (c *JsonConvertImpl) Encode(f interface{}) *public.Content {
	var ctxt public.Content
	ctxt.Code = http.StatusOK
	if e, b := f.(error); b {
		bytes, _ := json.Marshal(public.NewError(e.Error()))
		ctxt.Bytes = bytes
		ctxt.ContentType = public.Json
		ctxt.Code = http.StatusInternalServerError
		return &ctxt
	}

	kind := reflect.Indirect(reflect.ValueOf(f)).Kind()
	switch kind {
	case reflect.String:
		ctxt.Bytes = []byte(f.(string))
		ctxt.ContentType = public.Json
		return &ctxt
	default:
		bytes, _ := json.Marshal(f)
		ctxt.Bytes = bytes
		ctxt.ContentType = public.Json
		return &ctxt
	}
}
