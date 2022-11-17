package def

import (
	"gitee.com/fast_api/api/dwarf"
	"gitee.com/fast_api/api/mg"
	"net/http"
	"reflect"
	"sync"
)

//Fn [name]->

type MethodsPools struct {
	kv sync.Map
}

func init() {
	mg.Provide(func() *MethodsPools {
		return &MethodsPools{}
	})
}

func (m *MethodsPools) Get(name string) *MethodInfo {
	if v, b := m.kv.Load(name); b {
		return v.(*MethodInfo)
	} else {
		return nil
	}
}

func (m *MethodsPools) Set(name string, methodInfo *MethodInfo) {
	m.kv.Store(name, methodInfo)
}

type Param struct {
	Order int    `json:"order"`
	Name  string `json:"name"`
}

type MethodInfo struct {
	Pkg        string                    `json:"pkg"`
	Receive    string                    `json:"receive"`
	Method     interface{}               `json:"-"`
	MethodName string                    `json:"method_name"`
	Param      map[string]dwarf.ArgsMeta `json:"param"`
}

type Content struct {
	ContentType string
	Bytes       []byte
}

type Entry struct {
	Url    string
	Group  string
	Method string
	Fn     interface{}
	Ids    map[string]string
}

type ParamWarp struct {
	http.Request
	PTyp   reflect.Type
	PValue string
	Path   string
	PName  string
}
