package swagger

type RequestBodyObject struct {
	Description string         `json:"description,omitempty"`
	Content     map[string]any `json:"content,omitempty"`
	Required    bool           `json:"required,omitempty"`
}

// application/octet-stream
func RefRequestBody(ref string, typ string) *RequestBodyObject {
	if ref == "binary" && typ == "array" {
		return &RequestBodyObject{
			Content: map[string]any{
				"application/octet-stream": map[string]any{
					"schema": map[string]any{
						"format": ref,
					},
				},
			},
		}
	} else {
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

}
