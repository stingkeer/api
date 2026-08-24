// Package api is a declarative Go web framework: URL parameters are inferred
// from handler function signatures (parameter names are resolved from DWARF
// debug information), so a route is one line:
//
//	api.GET(func(name string) any { return "hello " + name }, "/hello")
//
// Build tags:
//
//	swagger — serve Swagger UI at /ui/ and the OpenAPI document at /api/swagger
//	http3   — serve over HTTP/3 (QUIC) instead of plain TCP
package api

import (
	stdhttp "net/http"

	"go.aew.app/api.v1/call"
	"go.aew.app/api.v1/call/rettypes"
	"go.aew.app/api.v1/def"
	"go.aew.app/api.v1/http"
	"go.aew.app/api.v1/kit/core"
)

// Route registration.
//
// Each verb is a func(f any, url string) def.Option:
//
//   - f must be a func. Its PARAMETER NAMES bind request values (resolved
//     from DWARF debug info — keep DWARF in production builds):
//
//     func(name string)            // GET /hello?name=x   → query param "name"
//     func(id string)  + "/u/<id>" // path segment <id> binds to parameter "id"
//
//   - Special parameter names and types:
//
//     body    any                  // the request body (JSON-decoded for structs)
//     def.StringReq / def.IntReq…  // required params (missing → HTTP 500 error)
//     http.Request / *http.Request // the raw request
//     def.Header                   // request/response header access
//     multipart.Reader             // multipart (file upload) body
//     *ws.WSCtx                    // upgrades the connection to WebSocket
//
//   - Return values become the response:
//
//     any / struct / slice         // serialized as JSON (strings as text/plain)
//     def.Error                    // JSON error body
//     api.NewResp(...)             // custom status code / headers / content type
//     api.NewStream(...)           // file download (Range + rate limit support)
//     api.Html(...)                // rendered HTML template
//     api.NewRedirect(...)         // HTTP 302
//
//   - URL patterns: static segments ("/users"), named parameters
//     ("/users/<id>"), regex constraints ("/users/<id:[0-9]+>"), and
//     wildcards ("/files/<rest:.*>" — matches across "/").
//
//   - The returned def.Option carries registration metadata; chain .Swagger
//     or .SetMiddleware on it, or group several via api.AddRoutes.
var (
	// HEAD registers a handler answering HEAD requests on url.
	// Handlers return values like GET; bodies are discarded per RFC 9110.
	HEAD = core.HttpM(stdhttp.MethodHead, def.DefaultContext)

	// GET registers a handler answering GET requests on url.
	//
	//	api.GET(func(name string) any {
	//		return map[string]string{"hello": name}
	//	}, "/hello")
	//
	// GET routes also answer HEAD requests automatically.
	GET = core.HttpM(stdhttp.MethodGet, def.DefaultContext)

	// POST registers a handler answering POST requests on url.
	// A parameter named "body" receives the request body:
	//
	//	api.POST(func(body LoginRequest) any { return body }, "/login")
	POST = core.HttpM(stdhttp.MethodPost, def.DefaultContext)

	// PUT registers a handler answering PUT requests on url.
	PUT = core.HttpM(stdhttp.MethodPut, def.DefaultContext)

	// PATCH registers a handler answering PATCH requests on url.
	PATCH = core.HttpM(stdhttp.MethodPatch, def.DefaultContext)

	// DELETE registers a handler answering DELETE requests on url.
	DELETE = core.HttpM(stdhttp.MethodDelete, def.DefaultContext)

	// OPTIONS registers a handler answering OPTIONS requests on url.
	OPTIONS = core.HttpM(stdhttp.MethodOptions, def.DefaultContext)

	// RegisterErrorHandler customizes how a recovered panic value of a given
	// type becomes the HTTP error response. When a handler panics with a
	// value of the registered reflect.Type, the handler converts it into a
	// response payload (serialized as JSON, status 500):
	//
	//	type ForbiddenError struct{ Resource string }
	//	api.RegisterErrorHandler(reflect.TypeOf(ForbiddenError{}),
	//		func(err any) any {
	//			return map[string]string{
	//				"error":    "forbidden",
	//				"resource": err.(ForbiddenError).Resource,
	//			}
	//		})
	//
	// Without a registered handler, panics with a string or error become
	// def.NewError(value) — {"error":…, "code":…} — and other values become
	// an empty body.
	RegisterErrorHandler = http.RegisterErrorHandler

	// AddHttpHandle registers a global HTTP interceptor executed around
	// every request, in Order() sequence. Handlers with a lower order run
	// earlier; returning true from Http stops the chain (the response is
	// considered written).
	//
	// Order ranges: 0 = runs before routing; 1-99 and ≥1000 = user space
	// (≥1000 runs after the route handler, in the response phase);
	// 100-999 are reserved for the framework itself.
	//
	//	api.AddHttpHandle(myintercept) // implements intercept.HttpIntercept
	AddHttpHandle = http.AddHttpHandle

	// RegisterTypeMapper teaches parameter binding about a custom type.
	// Implement def.Adapter — Mapper converts an incoming request value into
	// the parameter's reflect.Value, Register lists the handled types:
	//
	//	api.RegisterTypeMapper(&myUUIDAdapter{})
	//
	// Use RegisterGenericTypeMapper (in package call) for generic types.
	RegisterTypeMapper = call.RegisterTypeMapper

	// RegisterReturnHandler teaches the framework how to write a custom
	// return type. Implement def.RetAdapter (ContentType, Return, Register,
	// and optionally HttpStatus / AppendHeader / io.Closer) and register it:
	//
	//	api.RegisterReturnHandler(&myCsvReturn{})
	RegisterReturnHandler = http.RegisterReturnHandler

	// SetMethodProxy wraps every handler invocation with a custom proxy.
	// Proxies form a chain around the real call; fn.Invoke(m, args) invokes
	// the next link (the innermost link is the actual handler). Every proxy
	// may transform the []reflect.Value arguments in and the return values
	// out — this is how the method cache (package cache) is implemented.
	//
	// Registration order matters: proxies registered later wrap earlier ones
	// (they run first). Register during package init, before serving starts.
	//
	//	api.SetMethodProxy(func(fn call.MethodCaller, m *def.MethodInfo, args []reflect.Value) []reflect.Value {
	//		start := time.Now()
	//		ret := fn.Invoke(m, args) // call the wrapped chain
	//		log.Printf("%s took %s", m.MethodName, time.Since(start))
	//		return ret
	//	})
	SetMethodProxy = call.SetMethodProxy

	// NewStream wraps an io.Reader as a file-download response, with HTTP
	// Range support, content sniffing and optional rate limiting:
	//
	//	api.GET(func() any {
	//		f, _ := os.Open("video.mp4")
	//		return api.NewStream(f).
	//			SetName("video.mp4").          // Content-Disposition attachment
	//			SetRateLimit(1024 * 1024)      // bytes per second
	//	}, "/download")
	NewStream = rettypes.NewStream

	// Html renders a Go text/template and returns it as a text/html response:
	//
	//	api.Html("<h>{{.Title}}</h>", struct{ Title string }{"Hi"})
	Html = rettypes.NewHtml

	// HtmlView renders a template embedded in an fs.FS (embed.FS works):
	//
	//	//go:embed view
	//	var view embed.FS
	//
	//	func htmlList(name def.StringReq) any {
	//		return api.HtmlView(view, "view/list.html", messages)
	//	}
	HtmlView = rettypes.HtmlView

	// Static mounts a file system at a wildcard URL path:
	//
	//	//go:embed public
	//	var public embed.FS
	//
	//	func init() {
	//		api.Static("/admin/*", "prefix/dir", http.FS(public),
	//			mhttp.StaticRewrite("/admin/", ""),   // strip the /admin/ prefix
	//			mhttp.StaticDefaultFile("index.html"), // fallback for missing files
	//		)
	//	}
	//
	// Files not found fall through to regular routing (404 if nothing
	// matches). Both StaticRewrite and StaticDefaultFile are optional.
	Static = http.DefaultStatic.HandleStatic

	// NewRedirect answers the request with HTTP 302 and a Location header:
	//
	//	api.GET(func() any { return api.NewRedirect("https://example.com") }, "/go")
	NewRedirect = rettypes.NewRedirect

	// NewResp builds a fully customizable response:
	//
	//	return api.NewResp(payload).
	//		SetCode(http.StatusCreated).                  // any status code
	//		SetHeader(map[string]string{"X-Trace": id})   // extra headers
	//		SetReader(r).                                  // raw body (opt-out of serialization)
	//		SetContentType("application/xml")              // override content type
	NewResp = rettypes.NewResp
)
