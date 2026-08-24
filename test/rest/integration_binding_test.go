package rest

import (
	"net/http"
	"testing"

	"go.aew.app/api.v1"
	"go.aew.app/api.v1/def"
	"go.aew.app/api.v1/test/r"
)

// TestIntegrationQueryTypes verifies binding of every scalar query type,
// both valid values and zero-value defaults when omitted.
func TestIntegrationQueryTypes(t *testing.T) {
	r.Test(t, func() def.Option {
		return api.GET(func(
			b bool, i int, i8 int8, i16 int16, i32 int32, i64 int64,
			u uint, u8 uint8, f32 float32, f64 float64, s string,
		) any {
			return map[string]any{
				"b": b, "i": i, "i8": i8, "i16": i16, "i32": i32, "i64": i64,
				"u": u, "u8": u8, "f32": f32, "f64": f64, "s": s,
			}
		}, "/bind/types")
	}).Request().
		AddParam("b", "true").
		AddParam("i", "-42").
		AddParam("i8", "127").
		AddParam("i16", "32767").
		AddParam("i32", "2147483647").
		AddParam("i64", "9223372036854775807").
		AddParam("u", "42").
		AddParam("u8", "255").
		AddParam("f32", "1.5").
		AddParam("f64", "-2.25").
		AddParam("s", "héllo").
		Do(func(resp *r.Response) {
			resp.AssertStatus(http.StatusOK)
			// map[string]any serializes with sorted keys — compare per field.
			var got map[string]any
			if err := jsonUnmarshalString(resp.BodyString(), &got); err != nil {
				t.Fatalf("unmarshal: %v (body %q)", err, resp.BodyString())
			}
			want := map[string]any{
				"b": true, "i": float64(-42), "i8": float64(127), "i16": float64(32767),
				"i32": float64(2147483647), "i64": float64(9223372036854775807),
				"u": float64(42), "u8": float64(255), "f32": 1.5, "f64": -2.25, "s": "héllo",
			}
			for k, v := range want {
				if got[k] != v {
					t.Errorf("field %s: expect %v but %v", k, v, got[k])
				}
			}
		})
}

// TestIntegrationQueryDefaults: omitted params bind to zero values for every
// scalar kind (uint/float used to panic with "reflect: Call using zero Value").
func TestIntegrationQueryDefaults(t *testing.T) {
	r.Test(t, func() def.Option {
		return api.GET(func(u uint, f float64, b bool, s string, i int) any {
			return map[string]any{"u": u, "f": f, "b": b, "s": s, "i": i}
		}, "/bind/defaults")
	}).Request().Do(func(resp *r.Response) {
		resp.AssertStatus(http.StatusOK)
		resp.AssertBody(`{"b":false,"f":0,"i":0,"s":"","u":0}`)
	})
}

// TestIntegrationInvalidParam: unparseable numbers produce a JSON 500 error
// body, not a hung/empty response.
func TestIntegrationInvalidParam(t *testing.T) {
	r.Test(t, func() def.Option {
		return api.GET(func(n int) any { return n }, "/bind/bad-int")
	}).Request().AddParam("n", "not-a-number").Do(func(resp *r.Response) {
		resp.AssertStatus(http.StatusInternalServerError)
		var e def.Error
		if err := jsonUnmarshalString(resp.BodyString(), &e); err != nil {
			t.Fatal(err)
		}
		if e.ErrorMessage == "" {
			t.Error("expect non-empty error message")
		}
	})
}

// TestIntegrationPathParam: <name> segments bind by position — and the
// bound value is delivered to the function argument whose NAME matches the
// route variable name.
func TestIntegrationPathParam(t *testing.T) {
	r.Test(t, func() def.Option {
		return api.GET(func(id string, sub string) any {
			return id + "/" + sub
		}, "/bind/path/<id>/items/<sub>")
	})

	resp, err := http.Get(r.BaseURL() + "/bind/path/42/items/keys")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	b, _ := ioReadAll(resp)
	// plain string returns are served as text/plain (no JSON quoting)
	if string(b) != "42/keys" {
		t.Errorf("path params: expect %q but %q", "42/keys", string(b))
	}
}

// TestIntegrationPathRegexConstraint: <name:[0-9]+> only matches digits;
// non-matching paths 404 rather than binding garbage.
func TestIntegrationPathRegexConstraint(t *testing.T) {
	r.Test(t, func() def.Option {
		return api.GET(func(id string) any { return "id=" + id }, "/bind/re/<id:[0-9]+>")
	})

	resp, err := http.Get(r.BaseURL() + "/bind/re/123")
	if err != nil {
		t.Fatal(err)
	}
	b, _ := ioReadAll(resp)
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK || string(b) != "id=123" {
		t.Errorf("regex match: code=%d body=%q", resp.StatusCode, string(b))
	}

	resp2, err := http.Get(r.BaseURL() + "/bind/re/abc")
	if err != nil {
		t.Fatal(err)
	}
	defer resp2.Body.Close()
	if resp2.StatusCode != http.StatusNotFound {
		t.Errorf("regex non-match: expect 404 but %d", resp2.StatusCode)
	}
}

// TestIntegrationPathWildcard: <name:.*> captures the rest of the path,
// including slashes. The function argument name must match the route
// variable name — mismatched names bind the zero value silently.
func TestIntegrationPathWildcard(t *testing.T) {
	r.Test(t, func() def.Option {
		return api.GET(func(rest string) any { return rest }, "/bind/wild/<rest:.*>")
	})
	resp, err := http.Get(r.BaseURL() + "/bind/wild/a/b/c.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	b, _ := ioReadAll(resp)
	if string(b) != "a/b/c.txt" {
		t.Errorf("wildcard: expect %q but %q", "a/b/c.txt", string(b))
	}
}

// TestIntegrationQueryAndPathMixed: path params and query params bind
// together in one handler.
func TestIntegrationQueryAndPathMixed(t *testing.T) {
	r.Test(t, func() def.Option {
		return api.GET(func(id string, verbose bool) any {
			return map[string]any{"id": id, "verbose": verbose}
		}, "/bind/mixed/<id>")
	})
	resp, err := http.Get(r.BaseURL() + "/bind/mixed/77?verbose=true")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	b, _ := ioReadAll(resp)
	if string(b) != `{"id":"77","verbose":true}` {
		t.Errorf("mixed: expect %q but %q", `{"id":"77","verbose":true}`, string(b))
	}
}

// TestIntegrationJSONBody: `body` param decodes JSON into a struct.
func TestIntegrationJSONBody(t *testing.T) {
	type login struct {
		Name string `json:"name,omitempty"`
		Pass string `json:"pass,omitempty"`
	}
	r.Test(t, func() def.Option {
		return api.POST(func(body login) any { return body }, "/bind/json")
	}).Request().SetJsonBody(login{Name: "w", Pass: "12345"}).Do(func(resp *r.Response) {
		resp.AssertStatus(http.StatusOK)
		resp.AssertBody(`{"name":"w","pass":"12345"}`)
	})
}

// TestIntegrationRequiredParam: def.StringReq rejects a missing param with
// the framework error JSON.
func TestIntegrationRequiredParam(t *testing.T) {
	r.Test(t, func() def.Option {
		return api.GET(func(name def.StringReq) any { return name.String() }, "/bind/required")
	}).Request().Do(func(resp *r.Response) {
		resp.AssertStatus(http.StatusInternalServerError)
		var e def.Error
		if err := jsonUnmarshalString(resp.BodyString(), &e); err != nil {
			t.Fatal(err)
		}
		if e.Code == 0 {
			t.Error("expect non-zero error code")
		}
	})
}

// TestIntegrationHeaderParam: def.Header exposes request headers.
func TestIntegrationHeaderParam(t *testing.T) {
	r.Test(t, func() def.Option {
		return api.GET(func(h def.Header) any {
			return h.Get("X-Custom")
		}, "/bind/header")
	}).Request().AddHeader("X-Custom", "val-1").Do(func(resp *r.Response) {
		resp.AssertStatus(http.StatusOK)
		// plain string return → text/plain, unquoted
		resp.AssertBody("val-1")
	})
}

// TestIntegrationRequestParam: http.Request / *http.Request inject the real
// request (method, URL, headers).
func TestIntegrationRequestParam(t *testing.T) {
	r.Test(t, func() def.Option {
		return api.GET(func(req *http.Request) any {
			return map[string]any{"method": req.Method, "path": req.URL.Path}
		}, "/bind/request")
	}).Request().AddHeader("probe", "1").Do(func(resp *r.Response) {
		resp.AssertStatus(http.StatusOK)
		resp.AssertBody(`{"method":"GET","path":"/bind/request"}`)
	})
}
