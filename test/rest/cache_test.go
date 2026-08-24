package rest

import (
	"fmt"
	"testing"
	"time"

	"go.aew.app/api.v1"
	"go.aew.app/api.v1/cache"
	"go.aew.app/api.v1/def"
	"go.aew.app/api.v1/test/r"
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

// TestCacheMultiParam is a regression test: cache keys were previously built
// by ranging over the Param map, so with 2+ params the key order randomized
// per process and every request missed the cache.
func TestCacheMultiParam(t *testing.T) {
	invoked := 0
	const hits = 5
	r.Test(t, func() def.Option {
		return api.GET(func(a, b def.String[cache.Key]) (any, cache.Cache) {
			invoked++
			return a.V + b.V, cache.NewCacheImpl(time.Minute)
		}, "/cache-multi")
	}).Request().AddParam("a", "foo").AddParam("b", "bar").DoTimes(hits, func(resp *r.Response) {
		// The cached copy is stored/replayed as the serialized string value.
		resp.AssertBody("foobar")
	})
	if invoked != 1 {
		t.Errorf("multi-param cache: handler invoked %d times, want 1 (cache key not deterministic?)", invoked)
	}
}
