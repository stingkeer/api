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
	i, j := 0, 0
	r.Test(t, func() def.Option {
		return api.GET(func(s def.String[cache.Key]) (any, cache.Cache) {
			fmt.Println("invoke")
			i++
			return "hello", cache.NewCacheImpl(time.Second * 30)
		}, "/cache")
	}).Request().AddParam("s", "aaaa").DoTimes(3, func(resp *r.Response) {
		j++
	})
	if i != 1 || j != 3 {
		t.Errorf("TestCache Error")
	}

}
