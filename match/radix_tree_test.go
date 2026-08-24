package match

import (
	"fmt"
	"strings"
	"testing"
)

// runGet exercises store.Get with fresh pvalues like the real router does.
func runGet(t *testing.T, s *store, path string) (interface{}, map[string]string) {
	t.Helper()
	pv := make([]string, 10)
	data, pnames := s.Get(path, pv)
	if data == nil {
		return nil, nil
	}
	got := map[string]string{}
	for i, name := range pnames {
		if i < len(pv) && pv[i] != "" {
			got[name] = pv[i]
		}
	}
	return data, got
}

func TestStoreStaticAndParam(t *testing.T) {
	s := newStore()
	s.Add("/users", "list")
	s.Add("/users/<id>", "user")
	s.Add("/users/<id>/posts", "posts")
	s.Add("/users/<id>/posts/<pid:[0-9]+>", "post")

	cases := []struct {
		path string
		want interface{}
		vars map[string]string
	}{
		{"/users", "list", nil},
		{"/users/42", "user", map[string]string{"id": "42"}},
		{"/users/42/posts", "posts", map[string]string{"id": "42"}},
		{"/users/42/posts/7", "post", map[string]string{"id": "42", "pid": "7"}},
		// regex "[0-9]+" rejects "abc"; leftover path means no fallback to /users/<id>/posts
		{"/users/42/posts/abc", nil, nil},
		{"/nope", nil, nil},
		{"/users/42/nope", nil, nil},
	}
	for _, c := range cases {
		data, vars := runGet(t, s, c.path)
		if data != c.want {
			t.Errorf("Get(%q) = %v, want %v", c.path, data, c.want)
			continue
		}
		if c.vars == nil && vars != nil && len(vars) > 0 && c.want != nil {
			// vars may legitimately contain leftovers for non-matching; only check on match
		}
		if c.vars != nil {
			for k, v := range c.vars {
				if vars[k] != v {
					t.Errorf("Get(%q) var %q = %q, want %q (all: %v)", c.path, k, vars[k], v, vars)
				}
			}
		}
	}
}

func TestStoreFirstRegistrationWins(t *testing.T) {
	s := newStore()
	s.Add("/dup", "first")
	s.Add("/dup", "second")
	pv := make([]string, 10)
	data, _ := s.Get("/dup", pv)
	if data != "first" {
		t.Errorf("duplicate registration: got %v, want first", data)
	}
}

func TestStoreParamPriorityOrder(t *testing.T) {
	// static route registered after param route covering the same shape
	s := newStore()
	s.Add("/v1/<name>", "param")
	s.Add("/v1/health", "static")
	pv := make([]string, 10)
	data, _ := s.Get("/v1/health", pv)
	if data != "param" {
		t.Errorf("earlier registration must win: got %v, want param", data)
	}
}

func TestStoreRootAndEmpty(t *testing.T) {
	s := newStore()
	s.Add("/", "root")
	pv := make([]string, 10)
	if data, _ := s.Get("/", pv); data != "root" {
		t.Errorf("Get(/) = %v, want root", data)
	}
}

func TestStoreWildcardRegex(t *testing.T) {
	s := newStore()
	s.Add("/files/<path:.*>", "files")
	data, vars := runGet(t, s, "/files/a/b/c.txt")
	if data != "files" {
		t.Fatalf("Get = %v, want files", data)
	}
	if vars["path"] != "a/b/c.txt" {
		t.Errorf("path = %q, want a/b/c.txt", vars["path"])
	}
}

func TestStoreSplitNode(t *testing.T) {
	// inserting /car then /cart forces a node split
	s := newStore()
	s.Add("/car", "car")
	s.Add("/cart", "cart")
	s.Add("/carts", "carts")
	for path, want := range map[string]interface{}{
		"/car":   "car",
		"/cart":  "cart",
		"/carts": "carts",
		"/ca":    nil,
		"/cars":  nil,
	} {
		pv := make([]string, 10)
		if data, _ := s.Get(path, pv); data != want {
			t.Errorf("Get(%q) = %v, want %v", path, data, want)
		}
	}
}

func TestStoreManyRoutes(t *testing.T) {
	s := newStore()
	const n = 500
	for i := 0; i < n; i++ {
		s.Add(fmt.Sprintf("/api/v1/res%d/<id:[0-9]+>/sub%d", i, i%7), i)
	}
	for i := 0; i < n; i += 37 {
		path := fmt.Sprintf("/api/v1/res%d/99/sub%d", i, i%7)
		pv := make([]string, 10)
		data, _ := s.Get(path, pv)
		if data != i {
			t.Errorf("Get(%q) = %v, want %d", path, data, i)
		}
	}
}

func TestStoreStringDump(t *testing.T) {
	s := newStore()
	s.Add("/a/<id>", 1)
	s.Add("/b", 2)
	dump := s.String()
	// nodes store prefix-compressed keys ("a/" + "<id>"), not the full route literal
	if !strings.Contains(dump, "a/") || !strings.Contains(dump, "<id>") || !strings.Contains(dump, "b") {
		t.Errorf("String() dump missing routes:\n%s", dump)
	}
}

// BenchmarkStoreBuild measures registration-time memory: this is where the
// old [256]*node children arrays showed up (2KB per node, ~24B now).
func BenchmarkStoreBuild(b *testing.B) {
	routes := make([]string, 1000)
	for i := range routes {
		routes[i] = fmt.Sprintf("/api/v%d/resource%d/<id:[0-9]+>/sub", i%10, i)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := newStore()
		for j, r := range routes {
			s.Add(r, j)
		}
	}
}

func BenchmarkStoreGetStatic(b *testing.B) {
	s := newStore()
	for i := 0; i < 100; i++ {
		s.Add(fmt.Sprintf("/static/path%d", i), i)
	}
	pv := make([]string, 10)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		s.Get("/static/path50", pv)
	}
}

func BenchmarkStoreGetParam(b *testing.B) {
	s := newStore()
	for i := 0; i < 100; i++ {
		s.Add(fmt.Sprintf("/api/%d/<id:int>/detail", i), i)
	}
	pv := make([]string, 10)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		s.Get("/api/50/77/detail", pv)
	}
}
