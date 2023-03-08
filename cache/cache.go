package cache

import "time"

type Cache interface {
	ExpireTime() time.Duration
}

type Key interface {
}

type CacheImpl struct {
	time time.Duration
}

func NewCacheImpl(time time.Duration) *CacheImpl {
	return &CacheImpl{time: time}
}

func (c *CacheImpl) ExpireTime() time.Duration {
	return time.Second
}
