package rest

import (
	"encoding/json"
	"io"
	"net/http"
)

// jsonUnmarshalString unmarshals a JSON string payload into out.
func jsonUnmarshalString(s string, out any) error {
	return json.Unmarshal([]byte(s), out)
}

// ioReadAll reads resp.Body fully (ignoring errors, for short assertions).
func ioReadAll(resp *http.Response) ([]byte, error) {
	return io.ReadAll(resp.Body)
}
