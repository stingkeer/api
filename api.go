package api

import (
	stdhttp "net/http"

	"gitee.com/fast_api/api/call"
	"gitee.com/fast_api/api/call/rettypes"
	"gitee.com/fast_api/api/def"
	"gitee.com/fast_api/api/http"
	"gitee.com/fast_api/api/kit/core"
)

var (
	HEAD    = core.HttpM(stdhttp.MethodHead, def.DefaultContext)
	GET     = core.HttpM(stdhttp.MethodGet, def.DefaultContext)
	POST    = core.HttpM(stdhttp.MethodPost, def.DefaultContext)
	PUT     = core.HttpM(stdhttp.MethodPut, def.DefaultContext)
	PATCH   = core.HttpM(stdhttp.MethodPatch, def.DefaultContext)
	DELETE  = core.HttpM(stdhttp.MethodDelete, def.DefaultContext)
	OPTIONS = core.HttpM(stdhttp.MethodOptions, def.DefaultContext)

	// RegisterErrorHandler error handler
	RegisterErrorHandler = http.RegisterErrorHandler

	// AddHttpHandle http handler
	AddHttpHandle = http.AddHttpHandle

	// RegisterTypeMapper type handler
	RegisterTypeMapper = call.RegisterTypeMapper

	// RegisterReturnHandler register handler
	RegisterReturnHandler = http.RegisterReturnHandler

	NewStream = rettypes.NewStream

	Html = rettypes.NewHtml

	// HtmlView
	//
	//  go:embed view
	//  var view embed.FS
	//
	//  func htmpList(format def.StringReq) any {
	//        //.....messages
	//        return api.HtmlView(view, "view/list.html", messages)
	//  }
	//
	HtmlView = rettypes.HtmlView

	// Static static web
	Static = http.DefaultStatic.HandleStatic

	NewRedirect = rettypes.NewRedirect

	NewResp = rettypes.NewResp
)
