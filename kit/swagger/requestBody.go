package swagger

type RequestBodyObject struct {
	Description string         `json:"description,omitempty"`
	Content     map[string]any `json:"content,omitempty"`
	Required    bool           `json:"required,omitempty"`
}

func JsonRefRequestBody(ref string, typ string) *RequestBodyObject {
	return &RequestBodyObject{
		Content: map[string]any{
			"application/json": map[string]any{
				"schema": map[string]any{
					"type": typ,
					"$ref": ref,
				},
			},
		},
	}
}
