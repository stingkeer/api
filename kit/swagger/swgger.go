package swagger

import (
	"fmt"
	"math/big"
	"math/rand"
	"net/http"
	"os"
	"reflect"
	"strings"

	"gitee.com/fast_api/api/call/types"
	"gitee.com/fast_api/api/def"
	"gitee.com/fast_api/api/kit/core"
)

//TODO Response type no implement

//Implement the open-api protocol
//https://swagger.io/resources/open-api/
//https://swagger.io/specification/v3/
//https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.0.md

type Swagger struct {
	Servers    []Server                              `json:"servers,omitempty"`
	Openapi    string                                `json:"openapi,omitempty"`
	Paths      map[string]map[string]OperationObject `json:"paths,omitempty"`
	Schemes    []string                              `json:"schemes,omitempty"`
	Host       string                                `json:"host,omitempty"`
	Components map[string]any                        `json:"components,omitempty"`
	Info       SwaggerInfo                           `json:"info,omitempty"`
}

type Server struct {
	URL         string `json:"url"`
	Description string `json:"description"`
}

type SwaggerInfo struct {
	Title          string  `json:"title"`
	Description    string  `json:"description"`
	TermsOfService string  `json:"termsOfService"`
	Contact        Contact `json:"contact"`
	License        License `json:"license"`
	Version        string  `json:"version"`
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
		Openapi: "3.0.0",
		Paths:   genPaths(ctx),
		Info: SwaggerInfo{
			Title:       "api Gen",
			Description: "This is a sample server for api",
		},
		Schemes: []string{"http", "https"},
		Host:    host,
	}
	s.Components = make(map[string]any)
	s.Components["schemas"] = definitionsMap
	s.Components["securitySchemes"] = securityDefinitionsMap
	return s
}

func genPaths(ctx *def.Context) map[string]map[string]OperationObject {
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
						"application/json": map[string]any{},
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

		var reqBodys []*RequestBodyObject
		var params []ParameterObject = []ParameterObject{}
		for name, p := range info.Param {
			typ, format := parameterDataType(p.Typ)
			in, req := parameterIN(p.Typ)
			parameter := ParameterObject{
				Name:     name,
				In:       in,
				Required: req,
			}

			if in == "pass" {
				continue
			}

			if typ == "Object" || typ == "array" {
				if info.Method.HttpMethod == http.MethodPost {
					reqBodys = append(reqBodys, JsonRefRequestBody(format, typ))
				} else {
					parameter.Schema = map[string]any{
						"type": typ,
						"$ref": format,
					}
				}
			} else {
				parameter.Schema = map[string]any{
					"type": typ,
				}
				if format != "" {
					parameter.Schema["format"] = format
				}
				//Description of setting parameters
				if description, b := info.KV.Load(fmt.Sprintf("swagger.parameter.%s", name)); b {
					parameter.Description = description.(string)
				}
			}
			params = append(params, parameter)
		}
		//TODO mul req bodys
		if len(reqBodys) > 0 {
			mEn.RequestBody = reqBodys[0]
		} else {
			mEn.Parameters = params
		}
		entry[strings.ToLower(info.Method.HttpMethod)] = mEn
		en[info.Method.Url] = entry
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
				return "Object", definitions(t)
			}
		}
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Uint8, reflect.Uint16, reflect.Uint32:
		return "integer", "int32"
	case reflect.Int, reflect.Int64, reflect.Uint64:
		return "integer", "int64"
	case reflect.String:
		return "string", ""
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
	Typ        string          `json:"type,omitempty"`
	Properties map[string]Cell `json:"properties"`
}

type Cell struct {
	Typ    string `json:"type"`
	Format string `json:"format,omitempty"`
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
func definitions(t reflect.Type) string {
	properties := make(map[string]Cell)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		dt, format := parameterDataType(field.Type)
		fName := getFieldName("json", field)
		if fName == "" {
			fName = field.Name
		}
		properties[fName] = Cell{Typ: dt, Format: format}
	}
	key := t.Name()
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

var (
	header types.HeadType
)

// ParameterIN Required. The location of the parameter.
// The location of the parameter. Possible values are "query", "header", "path" or "cookie"
// TODO only support query
func parameterIN(t reflect.Type) (in string, require bool) {
	requireTyps := append(g.Register(), g0.Register()...)
	if index := search(len(requireTyps), func(i int) bool {
		return requireTyps[i] == t
	}); index > 0 {
		return "query", true
	}
	headerRegs := header.Register()
	if index := search(len(headerRegs), func(i int) bool {
		return headerRegs[i] == t
	}); index > 0 {
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
