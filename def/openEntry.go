package def

import (
	"gitee.com/fast_api/api/dwarf"
	"gitee.com/fast_api/api/mg"
	"gitee.com/fast_api/api/utils"
	"net/http"
	"reflect"
	"sync"
)

type MethodsPools struct {
	utils.Map[string, MethodInfo]
}

func init() {
	mg.Provide(func() *MethodsPools {
		return &MethodsPools{}
	})
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
}

type Content struct {
	ContentType string
	Bytes       []byte
}

type Entry struct {
	Url        string
	Group      string
	HttpMethod string
	Fn         interface{}
	Ids        sync.Map
}

type ParamWarp struct {
	http.Request
	PTyp   reflect.Type
	PValue string
	Path   string
	PName  string
}
