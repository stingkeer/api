//go:build swagger

package rest

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"go.aew.app/api.v1"
	"go.aew.app/api.v1/def"
	"go.aew.app/api.v1/kit/swagger"
	"go.aew.app/api.v1/test/r"
)

// TestIntegrationSwaggerUIServed: the embedded Swagger UI must actually be
// reachable (regression: /ui/* previously 404'd on every asset because the
// embed FS keeps its ui/ prefix while the handler stripped the URL prefix).
func TestIntegrationSwaggerUIServed(t *testing.T) {
	r.Test(t, func() def.Option {
		return api.GET(func() any { return "boot" }, "/sw/boot")
	})

	resp, err := http.Get(r.BaseURL() + "/ui/index.html")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("GET /ui/index.html: expect 200 but %d", resp.StatusCode)
	}
	b, _ := ioReadAll(resp)
	if !strings.Contains(string(b), "swagger-ui") {
		t.Error("UI html does not reference swagger-ui assets")
	}
	if ct := resp.Header.Get("Content-Type"); !strings.HasPrefix(ct, "text/html") {
		t.Errorf("UI Content-Type: expect text/html but %q", ct)
	}
}

// TestIntegrationSwaggerDocument: the generated OpenAPI document round-trips
// through the server: valid JSON, correct templates for regex routes, merged
// same-path methods, response schemas, operationIds and servers URL.
func TestIntegrationSwaggerDocument(t *testing.T) {
	type loginBody struct {
		Name string `json:"name"`
		Pass string `json:"pass"`
	}
	// same path, two methods → both must survive in the document
	r.Test(t, func() def.Option {
		return api.GET(func(id string) any { return id }, "/sw/user/<id:[0-9]+>")
	})
	r.Test(t, func() def.Option {
		return api.POST(func(body loginBody) any { return body }, "/sw/user")
	})
	r.Test(t, func() def.Option {
		return api.GET(func(name def.StringReq) any { return name.String() }, "/sw/req")
	})
	r.Test(t, func() def.Option {
		return api.GET(func() any { return []map[string]any{{"ok": true}} }, "/sw/complex")
	})

	r.Test(t, func() def.Option {
		return api.GET(func() any { return swagger.GenSwagger(def.DefaultContext) }, "/sw/doc")
	}).Request().Do(func(resp *r.Response) {
		resp.AssertStatus(http.StatusOK)

		var doc map[string]any
		if err := json.Unmarshal(resp.Body(), &doc); err != nil {
			t.Fatalf("invalid JSON document: %v", err)
		}
		if v, _ := doc["openapi"].(string); v != "3.0.3" {
			t.Errorf("openapi version = %v, want 3.0.3", doc["openapi"])
		}

		paths, _ := doc["paths"].(map[string]any)
		if paths == nil {
			t.Fatal("paths missing")
		}

		// regex route → clean {id} template (was {id:[0-9]+})
		userPath, ok := paths["/sw/user/{id}"]
		if !ok {
			t.Errorf("path /sw/user/{id} missing; got keys: %v", pathKeys(paths))
		} else {
			ops := userPath.(map[string]any)
			if _, ok := ops["get"]; !ok {
				t.Error("GET missing on /sw/user/{id}")
			}
			// path param documented with required=true
			getOp := ops["get"].(map[string]any)
			params, _ := getOp["parameters"].([]any)
			found := false
			for _, p := range params {
				po := p.(map[string]any)
				if po["name"] == "id" && po["in"] == "path" && po["required"] == true {
					found = true
				}
			}
			if !found {
				t.Errorf("path param id/required missing: %v", params)
			}
		}

		// same-path POST must be merged under /sw/user (not overwrite the entry)
		userRoot, ok := paths["/sw/user"].(map[string]any)
		if !ok {
			t.Fatalf("path /sw/user missing (merge failure); got keys: %v", pathKeys(paths))
		}
		post, ok := userRoot["post"].(map[string]any)
		if !ok {
			t.Fatal("POST missing on /sw/user (same-path method merge)")
		}
		if _, ok := post["requestBody"]; !ok {
			t.Error("POST requestBody missing")
		}

		// required query param flagged
		reqOp := paths["/sw/req"].(map[string]any)["get"].(map[string]any)
		reqParams, _ := reqOp["parameters"].([]any)
		if len(reqParams) == 0 {
			t.Fatal("def.StringReq param missing from document entirely")
		}
		foundReq := false
		for _, p := range reqParams {
			po := p.(map[string]any)
			if po["name"] == "name" && po["required"] == true {
				foundReq = true
			}
		}
		if !foundReq {
			t.Error("def.StringReq param not marked required")
		}

		// operationIds present and non-empty
		for _, ops := range paths {
			for _, op := range ops.(map[string]any) {
				om := op.(map[string]any)
				if id, _ := om["operationId"].(string); id == "" {
					t.Errorf("operationId missing on %v", om)
				}
			}
		}

		// error schema registered
		comps, _ := doc["components"].(map[string]any)
		if comps == nil {
			t.Fatal("components missing")
		}
		schemas, _ := comps["schemas"].(map[string]any)
		if _, ok := schemas["Error"]; !ok {
			t.Errorf("def.Error schema missing from components; got: %v", schemas)
		}

		// servers reflects the actual listen address
		servers, _ := doc["servers"].([]any)
		if len(servers) == 0 {
			t.Error("servers missing (was sourced from a never-set env var)")
		}
	})
}

func pathKeys(paths map[string]any) []string {
	var ks []string
	for k := range paths {
		ks = append(ks, k)
	}
	return ks
}
