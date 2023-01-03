package api

import "gitee.com/fast_api/api/def"

var (
	initFnCache InitFnCache
)

// InitFnCache used for cache init
type InitFnCache struct {
	fnCaches []*def.Entry
	initDone bool
}

func (c *InitFnCache) Add(def *def.Entry) {
	c.fnCaches = append(c.fnCaches, def)
}

func (c *InitFnCache) Range(v func(index int, en *def.Entry)) {
	for i, cache := range c.fnCaches {
		v(i, cache)
	}
}

func (c *InitFnCache) Len() int {
	return len(c.fnCaches)
}

func (c *InitFnCache) Init() bool {
	return c.initDone
}

func (c *InitFnCache) SetDone() {
	c.initDone = true
}

func (c *InitFnCache) Clear() {
	c.fnCaches = nil
}
