package api

import "gitee.com/fast_api/api/call/rettypes"

func Status(code int) *rettypes.Resp {
	return NewResp("").SetCode(code)
}

func Header(header map[string]string) *rettypes.Resp {
	return NewResp("").SetHeader(header)
}
