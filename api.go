package api

import (
	"gitee.com/fast_api/api/call"
	"gitee.com/fast_api/api/call/rettypes"
	"gitee.com/fast_api/api/def"
	"gitee.com/fast_api/api/http"
	"gitee.com/fast_api/api/kit/core"
	stdhttp "net/http"
)

var (
	GET  = core.HttpM(stdhttp.MethodGet, def.DefaultContext)
	POST = core.HttpM(stdhttp.MethodPost, def.DefaultContext)
	PUT  = core.HttpM(stdhttp.MethodPut, def.DefaultContext)

	// RegisterErrorHandler error handler
	RegisterErrorHandler = http.RegisterErrorHandler

	// AddHttpHandle http handler
	AddHttpHandle = http.AddHttpHandle

	// RegisterTypeMapper type handler
	RegisterTypeMapper = call.RegisterTypeMapper

	// RegisterReturnHandler register handler
	RegisterReturnHandler = http.RegisterReturnHandler

	NewStream = rettypes.NewStream

	Html     = rettypes.NewHtml
	HtmlView = rettypes.HtmlView

	// Static static web
	Static = http.DefaultStatic.HandleStatic

	NewRedirect = rettypes.NewRedirect
)
