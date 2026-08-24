// Package r provides a black-box test harness for services built with the
// api package. Each test registers its routes through Test; the harness
// boots the real HTTP server in-process — once per test binary, on an
// ephemeral 127.0.0.1 port — and hands back a small fluent client with
// built-in assertions.
//
// Typical use:
//
//	r.Test(t, func() def.Option {
//		return api.GET(func(name string) any {
//			return map[string]string{"hello": name}
//		}, "/hello")
//	}).Request().AddParam("name", "world").Do(func(resp *r.Response) {
//		resp.AssertBody(`{"hello":"world"}`)
//	})
package r

import (
	"bytes"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	"go.aew.app/api.v1"
	"go.aew.app/api.v1/def"
)

// Client is the per-test handle returned by Test.
type Client interface {
	// Request returns a builder for a request to the route registered by
	// this test (method and path come from the def.Option).
	Request() *Request
	// DoRequestNobody is sugar for Request().Do(respFn).
	DoRequestNobody(respFn func(resp *Response))
}

var (
	startOnce sync.Once
	baseURL   string // "http://127.0.0.1:<port>", set once when the server boots
)

// Test registers the routes declared by f and returns a Client bound to the
// returned option.
//
// The embedded HTTP server is process-scoped and starts lazily on the first
// call: registering a fresh server per test only ever produces "address
// already in use". An ephemeral port is picked so test binaries never fight
// over a fixed port — with stale processes or with each other.
func Test(t *testing.T, f func() def.Option) Client {
	t.Helper()
	// Force real route registration inside test binaries (isTestMode in
	// kit/core otherwise swaps api.GET for a no-op dummy).
	t.Setenv("API_TEST", "1")

	op := f()
	startOnce.Do(func() { startServer(t) })
	return &client{op: op, t: t}
}

// startServer boots the api server on an ephemeral port and blocks until it
// accepts connections (or fails the test with the reason).
func startServer(t *testing.T) {
	t.Helper()

	// Reserve a free port, release it, and let the server bind it. The
	// release/bind window is tiny and the failure path below reports it.
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("probe ephemeral port: %v", err)
	}
	addr := ln.Addr().String()
	ln.Close()

	bindErr := make(chan error, 1)
	go func() { bindErr <- api.StartService(api.WithListen(addr)) }()

	deadline := time.Now().Add(5 * time.Second)
	for {
		select {
		case err := <-bindErr:
			t.Fatalf("test server on %s exited: %v", addr, err)
		default:
		}
		conn, err := net.DialTimeout("tcp", addr, 250*time.Millisecond)
		if err == nil {
			conn.Close()
			baseURL = "http://" + addr
			return
		}
		if time.Now().After(deadline) {
			t.Fatalf("test server is not listening on %s after 5s", addr)
		}
		time.Sleep(25 * time.Millisecond)
	}
}

// BaseURL returns the base URL (scheme + host + ephemeral port) of the
// shared test server, for requests that need to bypass the Request builder.
func BaseURL() string { return baseURL }

// WSURL converts an HTTP path into a websocket URL on the test server,
// e.g. WSURL("/ws-echo") -> "ws://127.0.0.1:54321/ws-echo".
func WSURL(path string) string {
	return "ws" + strings.TrimPrefix(baseURL, "http") + path
}

// httpClient is shared by every request: bounded timeout so a hung handler
// fails the test instead of hanging it, and redirects are surfaced as the
// first response instead of being followed.
var httpClient = &http.Client{
	Timeout: 10 * time.Second,
	CheckRedirect: func(*http.Request, []*http.Request) error {
		return http.ErrUseLastResponse
	},
}

var _ Client = (*client)(nil)

type client struct {
	op def.Option
	t  *testing.T
}

func (c *client) Request() *Request {
	target, err := url.Parse(BaseURL() + c.op.Path())
	if err != nil {
		c.t.Fatalf("parse %s: %v", c.op.Path(), err)
	}
	req, err := http.NewRequest(c.op.Method(), target.String(), nil)
	if err != nil {
		c.t.Fatalf("create %s %s: %v", c.op.Method(), target, err)
	}
	return &Request{t: c.t, base: target, req: req, vs: make(url.Values)}
}

func (c *client) DoRequestNobody(respFn func(resp *Response)) {
	c.Request().Do(respFn)
}

// Request is a fluent builder around *http.Request.
type Request struct {
	t    *testing.T
	req  *http.Request
	base *url.URL // parsed URL without query, the source of truth for rebuilds
	vs   url.Values
}

// rebuild applies accumulated query params onto the outgoing URL.
func (r *Request) rebuild() {
	u := *r.base
	if len(r.vs) > 0 {
		u.RawQuery = r.vs.Encode()
	}
	r.req.URL = &u
}

// AddParam appends a query parameter.
func (r *Request) AddParam(key, value string) *Request {
	r.vs.Add(key, value)
	r.rebuild()
	return r
}

// AddHeader adds a header value.
func (r *Request) AddHeader(key, value string) *Request {
	r.req.Header.Add(key, value)
	return r
}

// SetHeader replaces a header value.
func (r *Request) SetHeader(key, value string) *Request {
	r.req.Header.Set(key, value)
	return r
}

// SetCookie adds a cookie to the request.
func (r *Request) SetCookie(name, value string) *Request {
	r.req.AddCookie(&http.Cookie{Name: name, Value: value})
	return r
}

// SetBody sets a raw request body.
func (r *Request) SetBody(obj []byte) *Request {
	r.req.Body = io.NopCloser(bytes.NewBuffer(obj))
	return r
}

// SetJsonBody marshals obj and sets it as the JSON request body.
func (r *Request) SetJsonBody(obj any) *Request {
	bs, err := json.Marshal(obj)
	if err != nil {
		r.t.Fatalf("marshal json body: %v", err)
	}
	r.req.Body = io.NopCloser(bytes.NewBuffer(bs))
	return r
}

// Do sends the request and hands the (fully drained) response to respFn.
func (r *Request) Do(respFn func(resp *Response)) {
	r.t.Helper()
	resp, err := httpClient.Do(r.req)
	if err != nil {
		r.t.Fatalf("send %s %s: %v", r.req.Method, r.req.URL, err)
	}
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		r.t.Fatalf("read response body: %v", err)
	}
	// Rewire Body so HttpResponse()/Dump() keep working after the drain.
	resp.Body = io.NopCloser(bytes.NewReader(body))
	respFn(&Response{t: r.t, resp: resp, body: body})
}

// DoTimes sends the same request n times.
func (r *Request) DoTimes(n int, respFn func(resp *Response)) {
	for i := 0; i < n; i++ {
		r.Do(respFn)
	}
}

// Response wraps an *http.Response whose body has already been read into
// memory and closed; assertions fail the owning test via *testing.T.
type Response struct {
	t    *testing.T
	resp *http.Response
	body []byte
}

// Code returns the HTTP status code.
func (res *Response) Code() int { return res.resp.StatusCode }

// Header returns the first value of a response header.
func (res *Response) Header(key string) string { return res.resp.Header.Get(key) }

// Headers returns the full response header map.
func (res *Response) Headers() http.Header { return res.resp.Header }

// Cookies returns the response's Set-Cookie cookies.
func (res *Response) Cookies() []*http.Cookie { return res.resp.Cookies() }

// HttpResponse exposes the underlying response. Its Body is rewired to a
// reader over the buffered body bytes.
func (res *Response) HttpResponse() *http.Response { return res.resp }

// Body returns the buffered response body.
func (res *Response) Body() []byte { return res.body }

// BodyString returns the buffered response body as a string.
func (res *Response) BodyString() string { return string(res.body) }

// Dump logs the full wire-format response (headers + body) via t.Log.
func (res *Response) Dump() {
	res.t.Helper()
	v, err := httputil.DumpResponse(res.resp, true)
	if err != nil {
		res.t.Logf("dump response: %v", err)
		return
	}
	res.t.Log(string(v))
}

// AssertStatus fails the test unless the status code equals want.
func (res *Response) AssertStatus(want int) {
	res.t.Helper()
	if res.Code() != want {
		res.t.Errorf("status: expect %d but %d", want, res.Code())
	}
}

// AssertBody fails the test unless the body equals want exactly.
func (res *Response) AssertBody(want string) {
	res.t.Helper()
	if res.BodyString() != want {
		res.t.Errorf("body: expect %q but %q", want, res.BodyString())
	}
}

// AssertHeader fails the test unless header key has exactly the value want.
func (res *Response) AssertHeader(key, want string) {
	res.t.Helper()
	if h := res.Header(key); h != want {
		res.t.Errorf("header %s: expect %q but %q", key, want, h)
	}
}

// AssetBody is a deprecated typo alias of AssertBody.
func (res *Response) AssetBody(want string) { res.AssertBody(want) }

// AssetHeader is a deprecated typo alias of AssertHeader.
func (res *Response) AssetHeader(key, want string) { res.AssertHeader(key, want) }
