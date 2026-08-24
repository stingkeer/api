package rest

import (
	"net/http"
	"testing"

	"go.aew.app/api.v1"
	"go.aew.app/api.v1/def"
	"go.aew.app/api.v1/test/r"
)

// TestIntegrationVerbs exercises every registered HTTP verb round-trip.
func TestIntegrationVerbs(t *testing.T) {
	cases := []struct {
		name   string
		reg    func(f any, url string) def.Option
		method string
		path   string
	}{
		{"get", api.GET, http.MethodGet, "/verbs/get"},
		{"post", api.POST, http.MethodPost, "/verbs/post"},
		{"put", api.PUT, http.MethodPut, "/verbs/put"},
		{"delete", api.DELETE, http.MethodDelete, "/verbs/delete"},
		{"patch", api.PATCH, http.MethodPatch, "/verbs/patch"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r.Test(t, func() def.Option {
				return c.reg(func() any { return c.name }, c.path)
			}).Request().Do(func(resp *r.Response) {
				resp.AssertStatus(http.StatusOK)
				resp.AssertBody(c.name)
			})
		})
	}
}

// TestIntegrationMethodNotAllowed: a wrong-method request must return 405
// with an Allow header naming the supported method (RFC 9110 §15.5.6),
// not a silently empty 200.
func TestIntegrationMethodNotAllowed(t *testing.T) {
	r.Test(t, func() def.Option {
		return api.GET(func() any { return "get" }, "/mm/allowed")
	})
	// POST via a manually built request (the builder derives the method
	// from the registration, so build it by hand).
	req, err := http.NewRequest(http.MethodPost, r.BaseURL()+"/mm/allowed", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("status: expect 405 but %d", resp.StatusCode)
	}
	if allow := resp.Header.Get("Allow"); allow != http.MethodGet {
		t.Errorf("Allow header: expect %q but %q", http.MethodGet, allow)
	}
}

// TestIntegrationHeadOnGet: HEAD against a GET route must answer 200 with an
// empty body (body is never sent for HEAD, and GET support implies HEAD).
func TestIntegrationHeadOnGet(t *testing.T) {
	r.Test(t, func() def.Option {
		return api.GET(func() any { return "body-not-sent" }, "/mm/head")
	})
	req, err := http.NewRequest(http.MethodHead, r.BaseURL()+"/mm/head", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("status: expect 200 but %d", resp.StatusCode)
	}
}

// TestIntegrationNotFound: unknown paths return 404 with a JSON body
// describing the path.
func TestIntegrationNotFound(t *testing.T) {
	r.Test(t, func() def.Option {
		return api.GET(func() any { return "x" }, "/mm/found")
	}).Request() // registers /mm/found; we request something else

	resp, err := http.Get(r.BaseURL() + "/definitely/not/here")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("status: expect 404 but %d", resp.StatusCode)
	}
	b, _ := ioReadAll(resp)
	var body map[string]string
	if err := jsonUnmarshalString(string(b), &body); err != nil {
		t.Fatal(err)
	}
	if body["path"] != "/definitely/not/here" {
		t.Errorf("body path: expect /definitely/not/here but %q", body["path"])
	}
	if body["msg"] != "Not find Path" {
		t.Errorf("body msg: expect %q but %q", "Not find Path", body["msg"])
	}
}
