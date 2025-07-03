//go:build runtest

package core

import (
	"go.aew.app/api.v1/def"
)

func HttpM(method string, ctx *def.Context) def.HttpMethod {
	return func(f interface{}, url string) def.Option {
		return &option{url: url, method: method}
	}
}
