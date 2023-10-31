package r

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"testing"

	"gitee.com/fast_api/api"
	"gitee.com/fast_api/api/def"
)

type Client interface {
	DoRequestNobody(respFn func(resp *Response))
	Request() *Request
	With(f func())
	Wait()
}

func Test(t *testing.T, f func() def.Option) Client {
	ops := f()
	go api.StartService()
	return &client{op: ops, t: t}
}

var _ Client = (*client)(nil)

type client struct {
	op def.Option
	t  *testing.T
}

// Wait implements Client.
func (*client) Wait() {
	var g sync.WaitGroup
	g.Add(1)
	g.Wait()
}

// With implements Client.
func (*client) With(f func()) {
	f()
}

type Request struct {
	req *http.Request
	t   *testing.T
	url *url.URL
}

type Response struct {
	resp *http.Response
	t    *testing.T
	buf  bytes.Buffer
	once sync.Once
}

func (res *Response) Dump() {
	v, err := httputil.DumpResponse(res.HttpResponse(), true)
	if err != nil {
		res.t.Error(err)
	}
	res.t.Log(string(v))
}

func (res *Response) HttpResponse() *http.Response {
	return res.resp
}

func (res *Response) Code() int {
	return res.resp.StatusCode
}

func (res *Response) Header(key string) string {
	return res.resp.Header.Get(key)
}

func (res *Response) BodyString() string {
	res.once.Do(func() {
		rs, err := io.ReadAll(res.resp.Body)
		if err != nil {
			res.t.Error(err)
		}
		res.buf.Write(rs)
	})
	return res.buf.String()
}

func (r *Request) Do(respFn func(resp *Response)) {
	resp, err := http.DefaultClient.Do(r.req)
	if err != nil {
		r.t.Error("The send request failed:", err)
		return
	}
	respFn(&Response{resp: resp, t: r.t})
}

func (r *Request) SetBody(obj []byte) *Request {
	r.req.Body = io.NopCloser(bytes.NewBuffer(obj))
	return r
}

func (r *Request) SetJsonBody(obj any) *Request {
	bs, err := json.Marshal(obj)
	if err != nil {
		r.t.Error(err)
	}
	r.req.Body = io.NopCloser(bytes.NewBuffer(bs))
	return r
}

func (r *Request) AddHeader(key, value string) *Request {
	r.req.Header.Add(key, value)
	return r
}

func (r *Request) AddParam(key, value string) *Request {
	vs := make(url.Values)
	vs.Add(key, value)
	url, err := url.Parse(fmt.Sprintf("%s?%s", r.url, vs.Encode()))
	if err != nil {
		r.t.Error(err)
	}
	r.req.URL = url
	return r
}

func (c *client) Request() *Request {
	url, err := url.Parse(fmt.Sprintf("http://localhost:8080%s", c.op.Path()))
	if err != nil {
		c.t.Error(err)
	}
	req, err := http.NewRequest(c.op.Method(), url.String(), nil)
	if err != nil {
		c.t.Error("The creation request failed:", err)
		return nil
	}
	return &Request{t: c.t, url: url, req: req}
}

// Request implements Client.
func (c *client) DoRequestNobody(respFn func(resp *Response)) {
	req, err := http.NewRequest(c.op.Method(), fmt.Sprintf("http://localhost:8080%s", c.op.Path()), nil)
	if err != nil {
		c.t.Error("The creation request failed:", err)
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		c.t.Error("The send request failed:", err)
		return
	}
	respFn(&Response{resp: resp, t: c.t})
}
