package rest

import (
	"fmt"
	"testing"
	"time"

	"go.aew.app/api"
	"go.aew.app/api/cache"
	"go.aew.app/api/def"
	"go.aew.app/api/test/r"
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
