// Package swagger generates an OpenAPI 3.0 document from registered routes.
//
// https://swagger.io/specification/v3/
// https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.0.md
package swagger

import (
	"math/big"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"time"

	"go.aew.app/api.v1/call/rettypes"
	"go.aew.app/api.v1/call/types"
	"go.aew.app/api.v1/def"
	"go.aew.app/api.v1/kit/core"
)

type Swagger struct {
	Openapi    string                                `json:"openapi,omitempty"`
	Info       SwaggerInfo                           `json:"info,omitempty"`
	Servers    []Server                              `json:"servers,omitempty"`
	Paths      map[string]map[string]OperationObject `json:"paths,omitempty"`
	Host       string                                `json:"host,omitempty"`
	Components map[string]any                        `json:"components,omitempty"`
}

type Server struct {
	URL         string `json:"url"`
	Description string `json:"description"`
}

type SwaggerInfo struct {
	Title          string  `json:"title"`
	Description    string  `json:"description"`
	TermsOfService string  `json:"termsOfService,omitempty"`
	Contact        Contact `json:"contact,omitempty"`
	License        License `json:"license,omitempty"`
	Version        string  `json:"version"`
}

type Contact struct {
	Name  string `json:"name,omitempty"`
	URL   string `json:"url,omitempty"`
	Email string `json:"email,omitempty"`
}

type License struct {
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}

type ParameterObject struct {
	Index       int            `json:"-"`
	Name        string         `json:"name,omitempty"`
	In          string         `json:"in,omitempty"`
	Description string         `json:"description,omitempty"`
	Required    bool           `json:"required"`
	Schema      map[string]any `json:"schema,omitempty"`
	Deprecated  bool           `json:"deprecated,omitempty"`
	Style       string         `json:"style,omitempty"`
}

type OperationObject struct {
	Deprecated  bool                  `json:"deprecated,omitempty"`
	Tags        []string              `json:"tags,omitempty"`
	Summary     string                `json:"summary,omitempty"`
	Description string                `json:"description,omitempty"`
	OperationId string                `json:"operationId,omitempty"`
	RequestBody *RequestBodyObject    `json:"requestBody,omitempty"`
	Parameters  []ParameterObject     `json:"parameters,omitempty"`
	Responses   map[uint]any          `json:"responses,omitempty"`
	Security    []map[string][]string `json:"security,omitempty"`
}

// customInfo lets applications brand the generated document via SetInfo
// (typically from an init function before the server starts).
var customInfo *SwaggerInfo

// SetInfo overrides the default document metadata (title, version, ...).
func SetInfo(info SwaggerInfo) { customInfo = &info }

func docInfo() SwaggerInfo {
	if customInfo != nil {
		return *customInfo
	}
	return SwaggerInfo{
		Title:       "Golang API Generate",
		Description: "This is a sample server for api",
		Version:     "1.0",
	}
}

// GenSwagger builds the OpenAPI document for every route registered in ctx.
// The server URL is taken from ctx.Listen, which StartService fills in.
func GenSwagger(ctx *def.Context) Swagger {
	s := Swagger{
		Openapi: "3.0.3",
		Paths:   genPaths(ctx),
		Info:    docInfo(),
	}
	if listen := strings.TrimSpace(ctx.Listen); listen != "" {
		host := strings.Replace(listen, "0.0.0.0", "localhost", 1)
		s.Servers = []Server{{URL: "http://" + host, Description: "api server"}}
		s.Host = host
	}
	s.Components = map[string]any{
		"schemas":         definitionsMap,
		"securitySchemes": securityDefinitionsMap,
	}
	return s
}

// ---------------------------------------------------------------------------
// Paths
// ---------------------------------------------------------------------------

func genPaths(ctx *def.Context) map[string]map[string]OperationObject {
	//clear all maps
	clear(definitionsMap)
	clear(securityDefinitionsMap)

	en := make(map[string]map[string]OperationObject)
	ctx.Pool.Range(func(_ string, info *def.MethodInfo) {
		mEn := OperationObject{
			OperationId: operationID(info.MethodName),
			Responses:   genResponses(info),
		}

		if securit, b := info.KV.Load("swagger.securit"); b {
			if sec, sb := securit.(map[string]*core.SecurityObject); sb {
				name := loadSecurityDefinition(sec)
				mEn.Security = []map[string][]string{
					{name: {}},
				}
			}
		}
		if commit, b := info.KV.Load("swagger.description"); b {
			mEn.Description = commit.(string)
		}
		if summary, b := info.KV.Load("swagger.summary"); b {
			mEn.Summary = summary.(string)
		}
		if tag, b := info.KV.Load("swagger.tag"); b {
			mEn.Tags = []string{tag.(string)}
		}

		var requestBody *RequestBodyObject
		var params []ParameterObject
		for _, p := range info.ParamList {
			in, req := parameterIN(info.Method.Url, p.Typ, p.Name)
			switch in {
			case "pass":
				continue
			case "file":
				requestBody = multipartBody()
				continue
			}
			if p.Name == "body" {
				requestBody = bodyOf(p.Typ)
				continue
			}
			// Struct-shaped scalars (def.StringReq, big.Int, time.Time, ...) are
			// single query params — parameterDataType reports their scalar type;
			// only real documentable structs (object → $ref) expand per-field.
			if typ, format := parameterDataType(p.Typ); typ == "object" && strings.HasPrefix(format, refPrefix) {
				params = append(params, structQueryParams(derefStruct(p.Typ), p.Order)...)
				continue
			}
			po := ParameterObject{
				Index:    p.Order,
				Name:     p.Name,
				In:       in,
				Required: req,
				Schema:   schemaOf(p.Typ),
			}
			if description, b := info.KV.Load("swagger.parameter." + p.Name); b {
				if d, ok := description.(string); ok {
					po.Description = d
				}
			}
			params = append(params, po)
		}

		sort.SliceStable(params, func(i, j int) bool {
			return params[i].Index < params[j].Index
		})
		mEn.Parameters = params
		mEn.RequestBody = requestBody

		url := pathTemplate(info.Method.Url)
		method := strings.ToLower(info.Method.HttpMethod)
		// merge with other methods already registered on the same path —
		// overwriting the whole entry dropped every method but the last
		entry, ok := en[url]
		if !ok {
			entry = make(map[string]OperationObject)
			en[url] = entry
		}
		entry[method] = mEn
	})
	return en
}

// operationID derives a stable, unique operation id from the runtime method
// name (e.g. "go.aew.app/api.v1/pkg.(*S).Handler.func1").
func operationID(methodName string) string {
	return identRe.ReplaceAllString(methodName, "_")
}

// genResponses builds the default response map. The 200 schema is derived
// from the handler's return type; 500 documents the def.Error JSON shape.
func genResponses(info *def.MethodInfo) map[uint]any {
	resp200 := map[string]any{"description": "successful operation"}
	if fnType := reflect.TypeOf(info.Method.Fn); fnType != nil && fnType.NumOut() > 0 {
		if ct, schema := returnSchema(fnType.Out(0)); ct != "" {
			resp200["content"] = map[string]any{
				ct: map[string]any{"schema": schema},
			}
		}
	}
	responses := map[uint]any{200: resp200}
	if errRef := definitions(reflect.TypeOf(def.Error{})); errRef != "" {
		responses[500] = map[string]any{
			"description": "error",
			"content": map[string]any{
				"application/json": map[string]any{
					"schema": map[string]any{"$ref": errRef},
				},
			},
		}
	}
	return responses
}

var (
	streamType = reflect.TypeOf((*rettypes.Stream)(nil)).Elem()
	htmlType   = reflect.TypeOf((*rettypes.Html)(nil)).Elem()
)

// returnSchema maps a handler return type to its response content type and
// schema. Stream → binary download, Html → text/html, interface → generic
// JSON (no schema), everything else → JSON schema of the concrete type.
func returnSchema(t reflect.Type) (string, map[string]any) {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t == streamType {
		return "application/octet-stream", map[string]any{"type": "string", "format": "binary"}
	}
	if t == htmlType {
		return "text/html", map[string]any{"type": "string"}
	}
	if t.Kind() == reflect.Interface {
		return "", nil
	}
	return "application/json", schemaOf(t)
}

// ---------------------------------------------------------------------------
// Parameter location
// ---------------------------------------------------------------------------

var (
	bigIntType = reflect.TypeOf((*big.Int)(nil)).Elem()
	timeType   = reflect.TypeOf((*time.Time)(nil)).Elem()
)

// requireTypes holds the required-parameter wrapper types (def.IntReq ...);
// requireBases holds their generic base names so that ANY instantiation
// (def.String[cache.Key], def.Int[MyTag], ...) is recognized, not just the
// [any] ones that happen to be pre-registered.
var (
	requireTypes = func() []reflect.Type {
		var g0 types.TypeRequire
		var g types.TypeRequireG
		return append(g.Register(), g0.Register()...)
	}()
	// defPkgPath pins the wrapper package exactly; the base names are the
	// generic type names (String, Int, ...).
	defPkgPath       = reflect.TypeOf(def.StringReq{}).PkgPath()
	requireBaseNames = map[string]bool{
		"Int":    true,
		"Int8":   true,
		"Int16":  true,
		"Int32":  true,
		"Int64":  true,
		"String": true,
	}
)

// isRequireType reports whether t is a required-parameter wrapper and returns
// the underlying value type.
func isRequireType(t reflect.Type) (reflect.Type, bool) {
	for _, rt := range requireTypes {
		if rt == t {
			return rt.Field(0).Type, true
		}
	}
	// generic instantiation: t.String() is "def.String[full/Type.Args]"
	if s := t.String(); strings.Contains(s, "[") {
		base := strings.TrimPrefix(s[:strings.IndexByte(s, '[')], "def.")
		if t.PkgPath() == defPkgPath && requireBaseNames[base] {
			if t.Kind() == reflect.Struct && t.NumField() > 0 {
				return t.Field(0).Type, true
			}
		}
	}
	return nil, false
}

// pathParamRe matches route variables: <id>, <id:[0-9]+>, <rest:.*>
var pathParamRe = regexp.MustCompile(`<([^<>:]+)(?::[^<>]*)?>`)
var identRe = regexp.MustCompile(`[^A-Za-z0-9_]+`)

// pathTemplate converts a framework route to an OpenAPI path template:
// /user/<id:[0-9]+> → /user/{id}
func pathTemplate(url string) string {
	return pathParamRe.ReplaceAllString(url, "{$1}")
}

// pathParamNames returns the variable names of a route, patterns stripped:
// /user/<id:[0-9]+>/x → ["id"]
func pathParamNames(url string) []string {
	matches := pathParamRe.FindAllStringSubmatch(url, -1)
	names := make([]string, 0, len(matches))
	for _, m := range matches {
		names = append(names, m[1])
	}
	return names
}

// parameterIN decides where a parameter lives: path, query, file (upload) or
// pass (framework-injected, not documented).
func parameterIN(url string, t reflect.Type, pName string) (in string, required bool) {
	if url != "" && pName != "" {
		for _, n := range pathParamNames(url) {
			if n == pName {
				return "path", true
			}
		}
	}
	if _, ok := isRequireType(t); ok {
		return "query", true
	}
	deref := t
	for deref.Kind() == reflect.Ptr {
		deref = deref.Elem()
	}
	if isPassType(t) || isPassType(deref) {
		return "pass", false
	}
	if t == multipartReaderType {
		return "file", true
	}
	return "query", false
}

// ---------------------------------------------------------------------------
// Schemas
// ---------------------------------------------------------------------------

const refPrefix = "#/components/schemas/"

// schemaOf builds a complete JSON schema object for t: primitives carry
// type/format, arrays carry items, structs become $refs into components.
func schemaOf(t reflect.Type) map[string]any {
	if t == nil {
		return map[string]any{"type": "object"}
	}
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() == reflect.Map {
		return map[string]any{
			"type":                 "object",
			"additionalProperties": schemaOf(t.Elem()),
		}
	}
	typ, format := parameterDataType(t)
	if typ == "object" && strings.HasPrefix(format, refPrefix) {
		return map[string]any{"$ref": format}
	}
	if typ == "array" {
		return map[string]any{
			"type":  "array",
			"items": schemaOf(t.Elem()),
		}
	}
	s := map[string]any{"type": typ}
	if format != "" {
		s["format"] = format
	}
	return s
}

// DataType maps Go types onto OpenAPI type/format pairs.
// https://swagger.io/docs/specification/data-models/data-types/
func parameterDataType(t reflect.Type) (typ, format string) {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t == bigIntType {
		return "integer", "int64"
	}
	if t == timeType {
		return "string", "date-time"
	}
	if v, ok := isRequireType(t); ok {
		return parameterDataType(v)
	}
	switch t.Kind() {
	case reflect.Struct:
		return "object", definitions(t)
	case reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Uint8, reflect.Uint16, reflect.Uint32:
		return "integer", "int32"
	case reflect.Int, reflect.Int64, reflect.Uint, reflect.Uint64:
		return "integer", "int64"
	case reflect.Float32:
		return "number", "float"
	case reflect.Float64:
		return "number", "double"
	case reflect.String:
		return "string", ""
	case reflect.Bool:
		return "boolean", ""
	case reflect.Interface:
		return "object", ""
	case reflect.Map:
		return "object", ""
	case reflect.Array, reflect.Slice:
		if t.Elem().Kind() == reflect.Uint8 {
			return "string", "byte"
		}
		return "array", ""
	}
	return "object", ""
}

// ---------------------------------------------------------------------------
// Components (schemas)
// ---------------------------------------------------------------------------

type Object struct {
	Typ                  string         `json:"type,omitempty"`
	Properties           map[string]any `json:"properties,omitempty"`
	AdditionalProperties bool           `json:"additionalProperties,omitempty"`
}

// Cell / CellArray are kept for API compatibility with earlier versions.
type Cell struct {
	Typ    string `json:"type"`
	Format string `json:"format,omitempty"`
}

type CellArray struct {
	Typ   string         `json:"type"`
	Items map[string]any `json:"items"`
}

var (
	definitionsMap         = make(map[string]Object)
	securityDefinitionsMap = make(map[string]any)
	definitionsBuilding    = map[reflect.Type]bool{}
)

func loadSecurityDefinition(v map[string]*core.SecurityObject) string {
	for name, obj := range v {
		if _, ex := securityDefinitionsMap[name]; !ex {
			securityDefinitionsMap[name] = obj
		}
		return name
	}
	return ""
}

// schemaKey names a schema: the type name for named structs, a sanitized
// deterministic signature for anonymous ones (previously rand.Int31() —
// every regeneration produced new schema names).
func schemaKey(t reflect.Type) string {
	if n := t.Name(); n != "" {
		return n
	}
	return identRe.ReplaceAllString(t.String(), "_")
}

// derefStruct returns the struct type behind t (through pointers), or nil.
func derefStruct(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() == reflect.Struct {
		return t
	}
	return nil
}

// definitions registers a struct in components/schemas and returns its $ref.
// Recursive types terminate: a type being built right now resolves to its
// final ref instead of recursing again.
func definitions(t reflect.Type) (ref string) {
	t = derefStruct(t)
	if t == nil {
		return ""
	}
	if isPassType(t) || t == multipartReaderType {
		return ""
	}
	key := schemaKey(t)
	ref = refPrefix + key
	if _, done := definitionsMap[key]; done {
		return ref
	}
	if definitionsBuilding[t] {
		return ref
	}
	definitionsBuilding[t] = true
	defer delete(definitionsBuilding, t)

	props := map[string]any{}
	collectProps(t, props, map[reflect.Type]bool{t: true})
	definitionsMap[key] = Object{Typ: "object", Properties: props}
	return ref
}

// jsonFieldName resolves a field's JSON name from its tag (splitting
// options like "name,omitempty" — a bare "name" previously yielded "").
// skip=true for json:"-".
func jsonFieldName(f reflect.StructField) (name string, skip bool) {
	tag, ok := f.Tag.Lookup("json")
	if !ok || tag == "" {
		return "", false
	}
	parts := strings.Split(tag, ",")
	if parts[0] == "-" {
		return "", true
	}
	return parts[0], false
}

// collectProps gathers the JSON properties of t into props. Embedded
// (anonymous) structs are flattened, mirroring encoding/json; visiting
// guards against embedding cycles.
func collectProps(t reflect.Type, props map[string]any, visiting map[reflect.Type]bool) {
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.PkgPath != "" && !f.Anonymous {
			continue // unexported field
		}
		ft := derefStruct(f.Type)
		if f.Anonymous && ft != nil {
			if !visiting[ft] {
				visiting[ft] = true
				collectProps(ft, props, visiting)
			}
			continue
		}
		name, skip := jsonFieldName(f)
		if skip || name == "" {
			continue
		}
		props[name] = schemaOf(f.Type)
	}
}

// structQueryParams expands a struct parameter into one query parameter per
// json-tagged field (embedded structs flattened).
func structQueryParams(t reflect.Type, index int) []ParameterObject {
	var out []ParameterObject
	var walk func(rt reflect.Type, visiting map[reflect.Type]bool)
	walk = func(rt reflect.Type, visiting map[reflect.Type]bool) {
		for i := 0; i < rt.NumField(); i++ {
			f := rt.Field(i)
			if f.PkgPath != "" && !f.Anonymous {
				continue
			}
			ft := derefStruct(f.Type)
			if f.Anonymous && ft != nil {
				if !visiting[ft] {
					visiting[ft] = true
					walk(ft, visiting)
				}
				continue
			}
			name, skip := jsonFieldName(f)
			if skip || name == "" {
				continue
			}
			_, req := isRequireType(f.Type)
			out = append(out, ParameterObject{
				Index:    index,
				Name:     name,
				In:       "query",
				Required: req,
				Schema:   schemaOf(f.Type),
			})
		}
	}
	walk(t, map[reflect.Type]bool{t: true})
	return out
}

func ResponseMimeTypes() []string {
	return []string{"application/json"}
}
