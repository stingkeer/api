package api

import (
	"encoding/base64"
	"github.com/sirupsen/logrus"
	"net/url"
)

/**
convert func result to []byte
*/
type Convert interface {
	convert(interface{}) []byte
	getContentType() string
}

/**
 为了用户更加方便的路由设置，beego 参考了 sinatra 的路由实现，支持多种方式的路由：
	beego.Router(“/api/?:id”, &controllers.RController{})
	默认匹配 //例如对于URL”/api/123”可以匹配成功，此时变量”:id”值为”123”
	beego.Router(“/api/:id”, &controllers.RController{})
	默认匹配 //例如对于URL”/api/123”可以匹配成功，此时变量”:id”值为”123”，但URL”/api/“匹配失败
	beego.Router(“/api/:id([0-9]+)“, &controllers.RController{})
	自定义正则匹配 //例如对于URL”/api/123”可以匹配成功，此时变量”:id”值为”123”
	beego.Router(“/user/:username([\\w]+)“, &controllers.RController{})
	正则字符串匹配 //例如对于URL”/user/astaxie”可以匹配成功，此时变量”:username”值为”astaxie”
	beego.Router(“/download/*.*”, &controllers.RController{})
	*匹配方式 //例如对于URL”/download/file/api.xml”可以匹配成功，此时变量”:path”值为”file/api”， “:ext”值为”xml”
	beego.Router(“/download/ceshi/*“, &controllers.RController{})
	*全匹配方式 //例如对于URL”/download/ceshi/file/api.json”可以匹配成功，此时变量”:splat”值为”file/api.json”
	beego.Router(“/:id:int”, &controllers.RController{})
	int 类型设置方式，匹配 :id为int 类型，框架帮你实现了正则 ([0-9]+)
	beego.Router(“/:hi:string”, &controllers.RController{})
	string 类型设置方式，匹配 :hi 为 string 类型。框架帮你实现了正则 ([\w]+)
	beego.Router(“/cms_:id([0-9]+).html”, &controllers.CmsController{})
	带有前缀的自定义正则 //匹配 :id 为正则类型。匹配 cms_123.html 这样的 url :id = 123
*/

type Match interface {
	/**
	 * return (result,param)
	 */
	match(url *url.URL, method string) (interface{}, url.Values)
	getMaps() map[string]Entry
}

type Api interface {
	GET(interface{}, string)
	POST(interface{}, string)
	PUT(interface{}, string)
	DELETE(interface{}, string)
	getMaps() map[string]Entry
}

type Caller interface {
	//function --> return
	call(f interface{}, params url.Values) interface{}
}

var M string

var methods map[interface{}]methodInfo

type param struct {
	order int
	name  string
	typ   string
}

type methodInfo struct {
	pkg    string
	method interface{}
	//map[order]Param
	param map[int]param
}

func initDef() {
	bytes, _ := base64.RawStdEncoding.DecodeString(M)
	logrus.Debugf("has decode string %s", string(bytes))
}
