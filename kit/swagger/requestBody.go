package swagger

import "reflect"

type RequestBodyObject struct {
	Description string         `json:"description,omitempty"`
	Content     map[string]any `json:"content,omitempty"`
	Required    bool           `json:"required,omitempty"`
}

// bodyOf builds the application/json request body schema for a `body`
// parameter of struct type t. Named structs become $ref schemas; everything
// else (string, slices, maps) inlines its schema.
func bodyOf(t reflect.Type) *RequestBodyObject {
	return &RequestBodyObject{
		Required: true,
		Content: map[string]any{
			"application/json": map[string]any{
				"schema": schemaOf(t),
			},
		},
	}
}

// multipartBody documents a multipart/form-data file upload (a handler
// parameter of type multipart.Reader).
func multipartBody() *RequestBodyObject {
	return &RequestBodyObject{
		Description: "file upload",
		Required:    true,
		Content: map[string]any{
			"multipart/form-data": map[string]any{
				"schema": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"file": map[string]any{
							"type":   "string",
							"format": "binary",
						},
					},
				},
			},
		},
	}
}

// RefRequestBody is kept for API compatibility with earlier versions.
func RefRequestBody(ref string, typ string) *RequestBodyObject {
	if typ == "array" {
		return multipartBody()
	}
	return &RequestBodyObject{
		Required: true,
		Content: map[string]any{
			"application/json": map[string]any{
				"schema": map[string]any{
					"$ref": ref,
				},
			},
		},
	}
}
