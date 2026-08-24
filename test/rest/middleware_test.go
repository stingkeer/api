package rest

import (
	"net/http"
	"testing"

	"go.aew.app/api.v1"
	"go.aew.app/api.v1/def"
	"go.aew.app/api.v1/test/r"
)

// authPass / authBlock are middlewares: returning nil continues the chain,
// any non-nil value short-circuits and becomes the response.
func mwAuth(token string) def.MiddleWare {
	return func(req *http.Request) any {
		if req.Header.Get("Authorization") != "Bearer "+token {
			return api.NewResp(map[string]string{"error": "unauthorized"}).SetCode(http.StatusUnauthorized)
		}
		return nil
	}
}

func mwAudit(counter *int) def.MiddleWare {
	return func(req *http.Request) any {
		*counter++
		return nil
	}
}

// TestIntegrationMiddlewarePass: all middlewares return nil → handler runs.
func TestIntegrationMiddlewarePass(t *testing.T) {
	audits := 0
	r.Test(t, func() def.Option {
		return api.GET(func() any { return "protected" }, "/mw/pass").
			SetMiddleware(mwAudit(&audits), mwAuth("secret"))
	}).Request().AddHeader("Authorization", "Bearer secret").Do(func(resp *r.Response) {
		resp.AssertStatus(http.StatusOK)
		resp.AssertBody("protected")
	})
	if audits != 1 {
		t.Errorf("audit middleware ran %d times, want 1", audits)
	}
}

// TestIntegrationMiddlewareShortCircuit: first middleware rejects → 401,
// handler never runs.
func TestIntegrationMiddlewareShortCircuit(t *testing.T) {
	handlerRan := false
	r.Test(t, func() def.Option {
		return api.GET(func() any {
			handlerRan = true
			return "protected"
		}, "/mw/block").SetMiddleware(mwAuth("secret"))
	}).Request().AddHeader("Authorization", "Bearer wrong").Do(func(resp *r.Response) {
		resp.AssertStatus(http.StatusUnauthorized)
		resp.AssertBody(`{"error":"unauthorized"}`)
	})
	if handlerRan {
		t.Error("handler must not run when middleware short-circuits")
	}
}

// TestIntegrationMiddlewareChainOrder: middlewares run in registration
// order; the first non-nil result wins.
func TestIntegrationMiddlewareChainOrder(t *testing.T) {
	var order []string
	rec := func(name string, ret any) def.MiddleWare {
		return func(req *http.Request) any {
			order = append(order, name)
			return ret
		}
	}
	r.Test(t, func() def.Option {
		return api.GET(func() any { return "never" }, "/mw/order").
			SetMiddleware(rec("one", nil), rec("two", "stopped-by-two"), rec("three", nil))
	}).Request().Do(func(resp *r.Response) {
		resp.AssertStatus(http.StatusOK)
		resp.AssertBody("stopped-by-two")
	})
	if len(order) != 2 || order[0] != "one" || order[1] != "two" {
		t.Errorf("middleware order: expect [one two] but %v", order)
	}
}

// TestIntegrationAddRoutesGroupMiddleware: AddRoutes(...).Middleware(...)
// applies to every route in the group.
func TestIntegrationAddRoutesGroupMiddleware(t *testing.T) {
	audits := 0
	var opA def.Option
	// r.Test boots the shared server (BaseURL) — routes register eagerly via
	// api.GET, so the group is completed right after.
	r.Test(t, func() def.Option {
		opA = api.GET(func() any { return "a" }, "/mw/group/a")
		return opA
	})
	api.AddRoutes(
		opA,
		api.GET(func() any { return "b" }, "/mw/group/b"),
	).Middleware(mwAudit(&audits))

	for _, path := range []string{"/mw/group/a", "/mw/group/b"} {
		resp, err := http.Get(r.BaseURL() + path)
		if err != nil {
			t.Fatal(err)
		}
		b, _ := ioReadAll(resp)
		resp.Body.Close()
		if string(b) != path[len(path)-1:] {
			t.Errorf("%s: expect %q but %q", path, path[len(path)-1:], string(b))
		}
	}
	if audits != 2 {
		t.Errorf("audit middleware ran %d times, want 2", audits)
	}
}
