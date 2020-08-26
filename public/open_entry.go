package public

import "gitee.com/aifuturewell/methods"

//Fn [name]->
type MetaMethods map[string]MethodInfo

var MethodsPools MetaMethods

type Param struct {
	Order int    `json:"order"`
	Name  string `json:"name"`
}

type MethodInfo struct {
	Pkg        string                      `json:"pkg"`
	Receive    string                      `json:"receive"`
	Method     interface{}                 `json:"-"`
	MethodName string                      `json:"method_name"`
	Param      map[string]methods.ArgsMeta `json:"param"`
}

type Header struct {
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
