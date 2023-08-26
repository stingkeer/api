package swgger

import (
	"fmt"
	"gitee.com/fast_api/api/call/types"
	"gitee.com/fast_api/api/def"
	"gitee.com/fast_api/api/mg"
	"math/rand"
	"reflect"
	"strings"
)

//impl open-api
//https://swagger.io/resources/open-api/

type Swagger struct {
	Swagger     string                      `json:"swagger,omitempty"`
	Paths       map[string]map[string]Entry `json:"paths"`
	Definitions map[string]Object           `json:"definitions"`
	Schemes     []string                    `json:"schemes"`
	Host        string                      `json:"host"`
}

type Parameter struct {
	Name        string            `json:"name"`
	In          string            `json:"in"`
	Description string            `json:"description"`
	Required    bool              `json:"required"`
	Type        string            `json:"type"`
	Format      string            `json:"format,omitempty"`
	Schema      map[string]string `json:"schema,omitempty"`
}

type Entry struct {
	Tags        []string    `json:"tags"`
	Summary     string      `json:"summary"`
	Description string      `json:"description"`
	OperationId string      `json:"operationId"`
	Consumes    []string    `json:"consumes"`
	Produces    []string    `json:"produces"`
	Parameters  []Parameter `json:"parameters"`
	Responses   struct {
		Field1 struct {
			Description string `json:"description"`
			Schema      struct {
				Ref string `json:"$ref"`
			} `json:"schema"`
		} `json:"200"`
	} `json:"responses"`
	Security []struct {
		PetstoreAuth []string `json:"petstore_auth"`
	} `json:"security"`
}

func GenSwagger() Swagger {
	definitionsMap = make(map[string]Object)
	return Swagger{
		Swagger:     "2.0",
		Paths:       genPaths(),
		Definitions: definitionsMap,
		Schemes:     []string{"http", "https"},
		Host:        "localhost:8080",
	}
}

func genPaths() map[string]map[string]Entry {
	en := make(map[string]map[string]Entry)
	err := mg.Invoke(func(pool *def.MethodsPools) {
		pool.Range(func(s string, info *def.MethodInfo) {
			entry := make(map[string]Entry)
			mEn := Entry{}
			var params []Parameter
			for name, p := range info.Param {
				typ, format := DataType(p.Typ)
				mp, req := ParameterIN(p.Typ)
				parameter := Parameter{
					Name:     name,
					In:       mp,
					Required: req,
				}
				if typ == "Object" {
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
			mEn.Produces = MimeTypes()
			entry[strings.ToLower(info.Method.HttpMethod)] = mEn
			en[info.Method.Url] = entry
		})
	})
	if err != nil {
		panic(err)
	}
	return en
}

var (
	g0 types.TypeRequire
	g  types.TypeRequireG
)

// DataType https://swagger.io/specification/v2/#data-type-format
func DataType(t reflect.Type) (typ, format string) {
	switch t.Kind() {
	case reflect.Struct:
		{
			requireTyps := append(g.Register(), g0.Register()...)
			if index := search(len(requireTyps), func(i int) bool {
				fmt.Println(requireTyps[i], t)
				return requireTyps[i] == t
			}); index > 0 {
				return DataType(requireTyps[index].Field(0).Type)
			} else {
				return "Object", definitions(t)
			}
		}
	case reflect.Int8, reflect.Int16, reflect.Int32:
		return "integer", "int32"
	case reflect.Int, reflect.Int64:
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

func definitions(t reflect.Type) string {
	properties := make(map[string]Cell)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		dt, format := DataType(field.Type)
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

func MimeTypes() []string {
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

// ParameterIN Required. The location of the parameter. Possible values are "query", "header", "path", "formData" or "body"
func ParameterIN(t reflect.Type) (in string, require bool) {
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
