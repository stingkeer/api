package swagger

import (
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"reflect"
	"regexp"
	"slices"
	"sort"
	"strings"

	"gitee.com/fast_api/api/call/types"
	"gitee.com/fast_api/api/def"
	"gitee.com/fast_api/api/dwarf"
	"gitee.com/fast_api/api/kit/core"
)

//TODO Response type no implement

//Implement the open-api protocol
//https://swagger.io/resources/open-api/
//https://swagger.io/specification/v3/
//https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.0.md

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
	Title          string `json:"title"`
	Description    string `json:"description"`
	TermsOfService string `json:"termsOfService"`
	// Contact        Contact `json:"contact"`
	License License `json:"license"`
	Version string  `json:"version"`
}
type Contact struct {
	Name  string `json:"name"`
	URL   string `json:"url"`
	Email string `json:"email"`
}
type License struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type ParameterObject struct {
	Index           int            `json:"-"`
	Name            string         `json:"name,omitempty"`
	In              string         `json:"in,omitempty"`
	Description     string         `json:"description,omitempty"`
	Required        bool           `json:"required"`
	Schema          map[string]any `json:"schema,omitempty"`
	Deprecated      bool           `json:"deprecated,omitempty"`
	RequestBody     map[string]any `json:"requestBody,omitempty"`
	AllowEmptyValue bool           `json:"allow_empty_value,omitempty"`
	Style           string         `json:"style,omitempty"`
}

type OperationObject struct {
	Deprecated  bool                  `json:"deprecated,omitempty"`
	Tags        []string              `json:"tags,omitempty"`
	Summary     string                `json:"summary,omitempty"`
	Description string                `json:"description,omitempty"`
	OperationId string                `json:"operation_id,omitempty"`
	RequestBody *RequestBodyObject    `json:"requestBody,omitempty"`
	Parameters  []ParameterObject     `json:"parameters,omitempty"`
	Responses   map[uint]any          `json:"responses,omitempty"`
	Security    []map[string][]string `json:"security,omitempty"`
}

func GenSwagger(ctx *def.Context) Swagger {
	host := os.Getenv("api.listen")
	s := Swagger{
		Servers: []Server{
			{
				URL:         host,
				Description: "swagger api ui",
			},
		},
		Openapi: "3.0.3",
		Paths:   genPaths(ctx),
		Info: SwaggerInfo{
			Title:       "Golang API Generate",
			Description: "This is a sample server for api",
			// Contact: Contact{
			// 	Name:  "api",
			// 	URL:   "api",
			// 	Email: "golang@gmail.com",
			// },
		},
		Host: host,
	}
	s.Components = make(map[string]any)
	s.Components["schemas"] = definitionsMap
	s.Components["securitySchemes"] = securityDefinitionsMap
	return s
}

func genStructParam(t reflect.Type) []ParameterObject {
	if t.Kind() != reflect.Struct {
		return nil
	}
	var res []ParameterObject
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		typ, format := parameterDataType(field.Type)
		in, req := parameterIN("", field.Type, "")
		newVar := ParameterObject{
			In:       in,
			Required: req,
		}
		newVar.Schema = map[string]any{
			"type":   typ,
			"format": format,
		}
		if v, ok := field.Tag.Lookup("json"); ok {
			newVar.Name = v
		} else {
			continue
		}

		res = append(res, newVar)
	}
	return res
}

func parseQuery(typ, format string, p dwarf.ArgsMeta) (params []ParameterObject, parameter *ParameterObject) {
	switch typ {
	case "object":
		{

			ps := genStructParam(p.Typ)
			if len(ps) > 0 {
				params = append(params, ps...)
			}

		}
	case "array":
		{
			var p ParameterObject
			p.Schema = map[string]any{
				"type": "array",
				"items": map[string]any{
					"$ref": format,
				},
			}
			parameter = &p
		}
	default:
		{
			var p ParameterObject
			p.Schema = map[string]any{
				"type": typ,
			}
			if format != "" {
				p.Schema["format"] = format
			}
			parameter = &p
		}
	}
	return
}

func genPaths(ctx *def.Context) map[string]map[string]OperationObject {
	//clear all maps
	clear(definitionsMap)
	clear(securityDefinitionsMap)

	en := make(map[string]map[string]OperationObject)
	ctx.Pool.Range(func(s string, info *def.MethodInfo) {
		entry := make(map[string]OperationObject)
		//initialization
		mEn := OperationObject{
			// Consumes: consumes(),
			Responses: map[uint]any{
				200: map[string]any{
					"description": "successful operation",
					"content": map[string]any{
						"application/json": map[string]any{
							"schema": map[string]string{
								"type": "object",
							},
						},
					},
				},
			},
		}
		//map[string]*SecurityObject
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

		var reqBodys []*RequestBodyObject
		var params []ParameterObject = []ParameterObject{}
		url := strings.ReplaceAll(info.Method.Url, "<", "{")
		url = strings.ReplaceAll(url, ">", "}")
		for name, p := range info.Param {
			typ, format := parameterDataType(p.Typ)
			in, req := parameterIN(info.Method.Url, p.Typ, name)

			if in == "pass" {
				continue
			}

			if name == "body" {
				reqBodys = append(reqBodys, JsonRefRequestBody(format, typ))
			} else {
				arrs, parameter := parseQuery(typ, format, p)
				if parameter != nil {
					parameter.Index = p.Order
					parameter.Name = name
					parameter.In = in
					parameter.Required = req
					//Description of setting parameters
					if description, b := info.KV.Load(fmt.Sprintf("swagger.parameter.%s", name)); b {
						parameter.Description = description.(string)
					}
					params = append(params, *parameter)
				}

				if len(arrs) > 0 {
					params = append(params, arrs...)
				}
			}

		}

		//sort param
		sort.Slice(params, func(i, j int) bool {
			return params[i].Index < params[j].Index
		})

		//TODO mul req bodys
		if len(reqBodys) > 0 {
			mEn.RequestBody = reqBodys[0]
		}
		mEn.Parameters = params
		entry[strings.ToLower(info.Method.HttpMethod)] = mEn
		en[url] = entry
	})
	return en
}

var (
	g0 types.TypeRequire
	g  types.TypeRequireG
)

// DataType https://swagger.io/specification/v2/#data-type-format
func parameterDataType(t reflect.Type) (typ, format string) {
	if reflect.TypeOf((*big.Int)(nil)) == t {
		return "integer", "int64"
	}
	switch t.Kind() {
	case reflect.Struct:
		{
			requireTyps := append(g.Register(), g0.Register()...)
			if index := search(len(requireTyps), func(i int) bool {
				return requireTyps[i] == t
			}); index > 0 {
				return parameterDataType(requireTyps[index].Field(0).Type)
			} else {
				return "object", definitions(t)
			}
		}
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Uint8, reflect.Uint16, reflect.Uint32:
		return "integer", "int32"
	case reflect.Int, reflect.Int64, reflect.Uint64:
		return "integer", "int64"
	case reflect.String:
		return "string", ""
	case reflect.Interface:
		return "object", ""
	case reflect.Bool:
		return "boolean", ""
	case reflect.Array, reflect.Slice:
		return "array", definitions(t.Elem())
	}
	return "", ""
}

// var mDefinitions = make(map[string]map[string])
var _ = `"definitions": {
    "ApiResponse": {
      "type": "object",
      "properties": {
        "code": {
          "type": "integer",
          "format": "int32"
        },
        "type": {
          "type": "string"
        },
        "message": {
          "type": "string"
        }
      }
    },`

type Object struct {
	Typ                  string         `json:"type,omitempty"`
	Properties           map[string]any `json:"properties,omitempty"`
	AdditionalProperties bool           `json:"additionalProperties"`
}

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

// Check type
func definitions(t reflect.Type) (ref string) {
	if t.Kind() != reflect.Struct {
		return ""
	}

	key := t.Name()

	if isDefTypes(t) {
		return ""
	}

	properties := make(map[string]any)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		dt, format := parameterDataType(field.Type)
		fName := getFieldName("json", field)
		if fName == "" {
			fName = field.Name
		}
		if dt == "array" {
			properties[fName] = CellArray{Typ: dt, Items: map[string]any{
				"$ref": format,
			}}
		} else {
			properties[fName] = Cell{Typ: dt, Format: format}
		}
	}

	if t.Name() == "" {
		key = fmt.Sprintf("Struct%d", rand.Int31())
	}

	definitionsMap[key] = Object{
		Typ:        "object",
		Properties: properties,
	}
	return fmt.Sprintf("#/components/schemas/%s", key)
}

func consumes() []string {
	return []string{"application/json"}
}

//Mime Types
//text/plain; charset=utf-8
//application/json
//application/vnd.github+json
//application/vnd.github.v3+json
//application/vnd.github.v3.raw+json
//application/vnd.github.v3.text+json
//application/vnd.github.v3.html+json
//application/vnd.github.v3.full+json
//application/vnd.github.v3.diff
//application/vnd.github.v3.patch

func ResponseMimeTypes() []string {
	return []string{"application/json"}
}

func search(t int, f func(int) bool) int {
	for i := 0; i < t; i++ {
		if f(i) {
			return i
		}
	}
	return -1
}

func getPaths(url string) []string {
	re := regexp.MustCompile(`<([^>]+)>`)
	matches := re.FindAllStringSubmatch(url, -1)

	var paths []string
	for _, match := range matches {
		if len(match) > 1 {
			paths = append(paths, match[1])
		}
	}
	return paths
}

// ParameterIN Required. The location of the parameter.
// The location of the parameter. Possible values are "query", "header", "path" or "cookie"
// TODO only support query
func parameterIN(url string, t reflect.Type, pName string) (in string, require bool) {
	if url != "" && pName != "" {
		paths := getPaths(url)
		if slices.Index(paths, pName) >= 0 {
			return "path", true
		}
	}
	requireTyps := append(g.Register(), g0.Register()...)
	if index := search(len(requireTyps), func(i int) bool {
		return requireTyps[i] == t
	}); index > 0 {
		return "query", true
	}

	if isDefTypes(t) {
		return "pass", true
	}
	return "query", false
}

func getFieldName(key string, s reflect.StructField) (fieldname string) {
	get := s.Tag.Get(key)
	if strings.Contains(get, ",") {
		return strings.Split(get, ",")[0]
	}
	return ""
}
