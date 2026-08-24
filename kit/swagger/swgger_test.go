package swagger

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
	"time"

	"go.aew.app/api.v1/def"
)

type Pet struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type LoginReq struct {
	Name string `json:"name"`
	Pass string `json:"pass"`
}

// TestPathTemplate: framework routes become valid OpenAPI templates,
// including regex-constrained and wildcard segments.
func TestPathTemplate(t *testing.T) {
	cases := map[string]string{
		"/user/<id>":            "/user/{id}",
		"/user/<id:[0-9]+>":     "/user/{id}",
		"/files/<rest:.*>":      "/files/{rest}",
		"/a/<x>/b/<y:[a-z]+>/c": "/a/{x}/b/{y}/c",
		"/plain/path":           "/plain/path",
	}
	for in, want := range cases {
		if got := pathTemplate(in); got != want {
			t.Errorf("pathTemplate(%q) = %q, want %q", in, got, want)
		}
	}
}

// TestPathParamNames strips patterns from route variables.
func TestPathParamNames(t *testing.T) {
	got := pathParamNames("/a/<x>/b/<y:[0-9]+>")
	if len(got) != 2 || got[0] != "x" || got[1] != "y" {
		t.Errorf("pathParamNames = %v, want [x y]", got)
	}
}

// TestParameterINRequired: the first registered require type (previously
// dropped by an `index > 0` bug) and generic instantiations
// def.String[cache.Key] must both be required query params.
func TestParameterINRequired(t *testing.T) {
	if in, req := parameterIN("", reflect.TypeOf(def.IntReq{}), "a"); in != "query" || !req {
		t.Errorf("IntReq: in=%q required=%v, want query/true (first-entry index bug)", in, req)
	}
	if in, req := parameterIN("", reflect.TypeOf(def.StringReq{}), "a"); in != "query" || !req {
		t.Errorf("StringReq: in=%q required=%v, want query/true", in, req)
	}
	// plain (non-required) param
	if _, req := parameterIN("", reflect.TypeOf(""), "a"); req {
		t.Error("plain string must not be required")
	}
	// path param
	if in, req := parameterIN("/u/<id>", reflect.TypeOf(""), "id"); in != "path" || !req {
		t.Errorf("path param: in=%q required=%v, want path/true", in, req)
	}
}

// TestRequireGenericInstantiation: def.String[cache.Key] (an arbitrary
// instantiation, not the pre-registered [any] one) resolves to string+required.
func TestRequireGenericInstantiation(t *testing.T) {
	typ := reflect.TypeOf(def.String[Pet]{})
	base, ok := isRequireType(typ)
	if !ok {
		t.Fatalf("def.String[Pet] not recognized as require type")
	}
	if base.Kind() != reflect.String {
		t.Errorf("base kind = %v, want string", base.Kind())
	}
	typ2, format := parameterDataType(typ)
	if typ2 != "string" || format != "" {
		t.Errorf("parameterDataType = %q/%q, want string/\"\"", typ2, format)
	}
}

// TestSchemaOfPrimitivesAndArrays: scalars, pointers, arrays of primitives
// (previously emitted an empty $ref), and maps.
func TestSchemaOfPrimitivesAndArrays(t *testing.T) {
	// []string → array with string items (was {"$ref": ""} before)
	arr := schemaOf(reflect.TypeOf([]string{}))
	if arr["type"] != "array" {
		t.Fatalf("[]string type = %v, want array", arr["type"])
	}
	items, _ := arr["items"].(map[string]any)
	if items == nil || items["type"] != "string" {
		t.Errorf("[]string items = %v, want {type: string}", arr["items"])
	}
	// pointer transparent
	p := schemaOf(reflect.TypeOf((*int)(nil)))
	if p["type"] != "integer" {
		t.Errorf("*int type = %v, want integer", p["type"])
	}
	// map
	m := schemaOf(reflect.TypeOf(map[string]int{}))
	if m["type"] != "object" {
		t.Errorf("map type = %v, want object", m["type"])
	}
	// time.Time
	tt := schemaOf(reflect.TypeOf(time.Now()))
	if tt["type"] != "string" || tt["format"] != "date-time" {
		t.Errorf("time.Time = %v, want string/date-time", tt)
	}
}

// TestDefinitions: struct schema registration, json tag handling (with and
// without options), json:"-" skipping, embedded flattening, and anonymous
// struct naming determinism.
func TestDefinitions(t *testing.T) {
	clear(definitionsMap)
	clear(definitionsBuilding)

	type Inner struct {
		A string `json:"a"`
		B int    `json:"b,omitempty"`
	}
	type Outer struct {
		Inner
		C       []Pet  `json:"pets"`
		Skipped string `json:"-"`
		NoTag   string
	}
	ref := definitions(reflect.TypeOf(Outer{}))
	if !strings.HasPrefix(ref, refPrefix) {
		t.Fatalf("ref = %q, want prefix %q", ref, refPrefix)
	}
	key := strings.TrimPrefix(ref, refPrefix)
	obj, ok := definitionsMap[key]
	if !ok {
		t.Fatalf("schema %q not registered", key)
	}
	// embedded flattened
	if _, ok := obj.Properties["a"]; !ok {
		t.Error("embedded field json:a missing")
	}
	if _, ok := obj.Properties["b"]; !ok {
		t.Error("embedded field json:b missing")
	}
	// slice of structs → array of $ref
	pets, _ := obj.Properties["pets"].(map[string]any)
	if pets == nil || pets["type"] != "array" {
		t.Fatalf("pets = %v, want array", obj.Properties["pets"])
	}
	petItems, _ := pets["items"].(map[string]any)
	if petItems == nil || !strings.HasPrefix(petItems["$ref"].(string), refPrefix) {
		t.Errorf("pets items = %v, want $ref to Pet", pets["items"])
	}
	// json:"-" skipped
	if _, ok := obj.Properties["Skipped"]; ok {
		t.Error("json:\"-\" field must be skipped")
	}
	if _, ok := obj.Properties["NoTag"]; ok {
		t.Error("untagged exported field must be skipped")
	}

	// anonymous struct gets a deterministic (non-random) name, stable across
	// regenerations
	type anonOuter struct {
		X int `json:"x"`
	}
	ref1 := definitions(reflect.TypeOf(anonOuter{}))
	clear(definitionsMap)
	ref2 := definitions(reflect.TypeOf(anonOuter{}))
	if ref1 != ref2 {
		t.Errorf("anonymous schema unstable: %q vs %q (rand naming bug)", ref1, ref2)
	}
}

// TestRecursiveSchema: self-referencing types terminate instead of
// stack-overflowing.
func TestRecursiveSchema(t *testing.T) {
	clear(definitionsMap)
	clear(definitionsBuilding)
	type Node struct {
		Next *Node `json:"next,omitempty"`
	}
	ref := definitions(reflect.TypeOf(Node{}))
	if !strings.HasPrefix(ref, refPrefix) {
		t.Fatalf("recursive schema ref = %q", ref)
	}
}

// TestOperationId: runtime function names become swagger-safe identifiers.
func TestOperationId(t *testing.T) {
	got := operationID("go.aew.app/api.v1/test.(*S).Handle.func1")
	if strings.ContainsAny(got, "./()") {
		t.Errorf("operationID = %q, still contains separators", got)
	}
}

// TestStructQueryParams: a struct parameter expands to one query param per
// json-tagged field.
func TestStructQueryParams(t *testing.T) {
	type Q struct {
		Page int          `json:"page"`
		Size def.Int32Req `json:"size"`
		Tag  string       `json:"-"`
	}
	ps := structQueryParams(reflect.TypeOf(Q{}), 3)
	if len(ps) != 2 {
		t.Fatalf("params = %d, want 2", len(ps))
	}
	if ps[0].Name != "page" || ps[0].In != "query" {
		t.Errorf("page param = %+v", ps[0])
	}
	if !ps[1].Required {
		t.Error("def.Int32Req field must be required")
	}
}

// TestBodyOfJSON serializes sanely.
func TestBodyOfJSON(t *testing.T) {
	rb := bodyOf(reflect.TypeOf(LoginReq{}))
	bs, err := json.Marshal(rb)
	if err != nil {
		t.Fatal(err)
	}
	s := string(bs)
	if !strings.Contains(s, `"application/json"`) || !strings.Contains(s, refPrefix) {
		t.Errorf("body json = %s", s)
	}
	if !rb.Required {
		t.Error("body must be required")
	}
}
