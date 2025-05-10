package api

import "go.aew.app/api.v1/call/rettypes"

func Status(code int) *rettypes.Resp {
	return NewResp("").SetCode(code)
}

func Header(header map[string]string) *rettypes.Resp {
	return NewResp("").SetHeader(header)
}
