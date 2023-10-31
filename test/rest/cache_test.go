package rest

import (
	"fmt"
	"testing"
	"time"

	"gitee.com/fast_api/api"
	"gitee.com/fast_api/api/cache"
	"gitee.com/fast_api/api/def"
	"gitee.com/fast_api/api/test/r"
)

func TestCache(t *testing.T) {
	r.Test(t, func() def.Option {
		return api.GET(func(s def.String[cache.Key]) (any, cache.Cache) {
			fmt.Println("invoke")
			return "hello", cache.NewCacheImpl(time.Second * 30)
		}, "/cache")
	}).Request().AddParam("s", "aaaa").Do(func(resp *r.Response) {
		fmt.Println(resp.BodyString())
	})

}
