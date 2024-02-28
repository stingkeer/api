package swagger

import (
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"reflect"
	"strings"

	"gitee.com/fast_api/api/call/types"
	"gitee.com/fast_api/api/def"
)

//TODO Response type no implement

//Implement the open-api protocol
//https://swagger.io/resources/open-api/
//https://swagger.io/specification/v3/

type Swagger struct {
	Swagger     string                      `json:"swagger,omitempty"`
	Paths       map[string]map[string]Entry `json:"paths"`
	Definitions map[string]Object           `json:"definitions"`
	Schemes     []string                    `json:"schemes"`
	Host        string                      `json:"host"`
}

type Parameter struct {
	Name        string            `json:"name,omitempty"`
	In          string            `json:"in,omitempty"`
	Description string            `json:"description,omitempty"`
	Required    bool              `json:"required,omitempty"`
	Type        string            `json:"type,omitempty"`
	Format      string            `json:"format,omitempty"`
	Schema      map[string]string `json:"schema,omitempty"`
}

type Entry struct {
	Tags        []string     `json:"tags,omitempty"`
	Summary     string       `json:"summary,omitempty"`
	Description string       `json:"description,omitempty"`
	OperationId string       `json:"operation_id,omitempty"`
	Consumes    []string     `json:"consumes,omitempty"`
	Produces    []string     `json:"produces,omitempty"`
	Parameters  []Parameter  `json:"parameters,omitempty"`
	Responses   map[uint]any `json:"responses,omitempty"`
	Security    []struct {
		PetstoreAuth []string `json:"petstore_auth,omitempty"`
	} `json:"security,omitempty"`
}

func GenSwagger(ctx *def.Context) Swagger {
	definitionsMap = make(map[string]Object)
	host := os.Getenv("api.listen")
	return Swagger{
		Swagger:     "2.0",
		Paths:       genPaths(ctx),
		Definitions: definitionsMap,
		Schemes:     []string{"http", "https"},
		Host:        host,
	}
}

func genPaths(ctx *def.Context) map[string]map[string]Entry {
	en := make(map[string]map[string]Entry)
	ctx.Pool.Range(func(s string, info *def.MethodInfo) {
		entry := make(map[string]Entry)
		//initialization
		mEn := Entry{
			Consumes: consumes(),
			Responses: map[uint]any{
				200: map[string]string{
					"description": "successful operation",
				},
			},
		}

		if commit, b := info.KV.Load("swagger.description"); b {
			mEn.Description = commit.(string)
		}
		if summary, b := info.KV.Load("swagger.summary"); b {
			mEn.Summary = summary.(string)
		}
		var params []Parameter = []Parameter{}
		for name, p := range info.Param {
			typ, format := parameterDataType(p.Typ)
			mp, req := parameterIN(p.Typ)
			parameter := Parameter{
				Name:     name,
				In:       mp,
				Required: req,
			}
			//Description of setting parameters
			if description, b := info.KV.Load(fmt.Sprintf("swagger.parameter.%s", name)); b {
				parameter.Description = description.(string)
			}
			if typ == "Object" && format != "" {
				parameter.Schema = map[string]string{
					"$ref": format,
				}
			} else {
				parameter.Type = typ
				parameter.Format = format
			}
			params = append(params, parameter)
		}
		mEn.Parameters = params
		mEn.Produces = ResponseMimeTypes()
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

var definitionsMap map[string]Object

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
	return fmt.Sprintf("#/definitions/%s", key)
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

// ParameterIN Required. The location of the parameter.
// Possible values are "query", "header", "path", "formData" or "body"
func parameterIN(t reflect.Type) (in string, require bool) {
	requireTyps := append(g.Register(), g0.Register()...)
	switch t.Kind() {
	case reflect.Struct:
		if index := search(len(requireTyps), func(i int) bool {
			return requireTyps[i] == t
		}); index > 0 {
			return "query", true
		}
		return "body", false
	case reflect.String:
		return "query", false
	}
	return "", false
}

func getFieldName(key string, s reflect.StructField) (fieldname string) {
	get := s.Tag.Get(key)
	if strings.Contains(get, ",") {
		return strings.Split(get, ",")[0]
	}
	return ""
}
