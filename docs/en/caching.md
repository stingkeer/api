# Caching

The framework provides method-level caching that transparently caches handler return values.

## Basic Usage

Return `cache.ImplCache` as the **second return value** to enable caching:

```go
import "go.aew.app/api.v1/cache"

api.GET(func() (interface{}, *cache.ImplCache) {
    data := expensiveDatabaseQuery()
    return data, cache.NewCacheImpl(5 * time.Minute)
}, "/cached")
```

When this endpoint is called:
1. The framework computes a cache key from the method name + parameter values
2. On cache hit → returns cached data, handler is **not called**
3. On cache miss → calls handler, caches the serialized result with the given TTL

## Cache Interface

```go
type Cache interface {
    ExpireTime_() time.Duration
}
```

The `ExpireTime_()` method returns the cache TTL. Any type implementing this interface can be used as the second return value.

## Custom Persistence Cache

By default, caching uses an in-memory `sync.Map`. Replace it with a custom implementation (e.g., Redis):

```go
type RedisCache struct{}

func (r *RedisCache) Set(key []byte, value []byte, ttl time.Duration) {
    // store in Redis
}

func (r *RedisCache) Get(key []byte) []byte {
    // retrieve from Redis, return nil if expired/missing
}

// Register at init time
func init() {
    cache.SetPersistenceCacheImpl(&RedisCache{})
}
```

Interface:

```go
type PersistenceCache interface {
    Set(key []byte, value []byte, ttl time.Duration)
    Get(key []byte) []byte
}
```

## Custom Cache Key/Value Encoding

Override how cache keys and values are encoded:

```go
type CustomEncoder struct{}

func (c *CustomEncoder) EncodeKey(m *def.MethodInfo, args []reflect.Value) []byte {
    // return cache key bytes
}

func (c *CustomEncoder) EncodeValue(s def.Serialize, v reflect.Value) []byte {
    // return cache value bytes
}

func init() {
    cache.SetProcessCacheImpl(&CustomEncoder{})
}
```

Interface:

```go
type ProcessCache interface {
    EncodeKey(m *def.MethodInfo, args []reflect.Value) []byte
    EncodeValue(s def.Serialize, v reflect.Value) []byte
}
```

## How It Works

Caching is implemented as a **method proxy** (`call.SetMethodProxy`) that wraps the handler invocation:

```
Request arrives
    │
    ├── Compute cache key (method name + parameter values)
    │
    ├── Check persistence cache
    │   ├── Hit → return cached bytes directly
    │   └── Miss → call handler
    │       ├── Check 2nd return value for Cache interface
    │       ├── If present → serialize result, store in cache
    │       └── Return result
    │
    └── Response
```

The cache proxy is the outermost layer in the method invocation chain, so it executes before any other proxies.
