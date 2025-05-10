package def

import (
	"reflect"
	"runtime"
	"sync"

	"go.aew.app/api.v1/dwarf"
	"go.aew.app/api.v1/log"
	"go.aew.app/api.v1/utils"
)

type MethodsPools struct {
	utils.Map[string, *MethodInfo]
}

func (m *MethodsPools) FuncInfo(fn any) *MethodInfo {
	name := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
	mInfo := m.Get(name)
	if mInfo != nil {
		return mInfo
	}
	log.Errorf("not find name [%s]", name)
	return nil
}

type Param struct {
	Order int    `json:"order"`
	Name  string `json:"name"`
}

type MethodInfo struct {
	Pkg        string                    `json:"pkg"`
	Receive    string                    `json:"receive"`
	Method     *Entry                    `json:"-"`
	MethodName string                    `json:"method_name"`
	Param      map[string]dwarf.ArgsMeta `json:"param"`
	Middleware []MiddleWare
	KV         sync.Map
}

type Content struct {
	ContentType string
	Bytes       []byte
}

// Empty
type Empty string

type Entry struct {
	Url        string
	Group      string
	HttpMethod string
	Fn         interface{}
	Ids        utils.Map[string, string]
}

type ParamWarp struct {
	Request
	PTyp   reflect.Type
	PValue string
	Path   string
	PName  string
}
