package rest

import (
	"fmt"
	"testing"
	"time"

	"gitee.com/fast_api/api"
	"gitee.com/fast_api/api/cache"
	"gitee.com/fast_api/api/def"
)

func TestCache(t *testing.T) {
	api.GET(func(s def.String[cache.Key]) (any, cache.Cache) {
		fmt.Println("invoke")
		return "hello", cache.NewCacheImpl(time.Second * 30)
	}, "/cache")
	api.StartService(nil)
}
