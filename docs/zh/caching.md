# 缓存系统

框架提供方法级缓存，透明缓存 handler 的返回值。

## 基本用法

将 `cache.ImplCache` 作为**第二个返回值**启用缓存：

```go
import "go.aew.app/api.v1/cache"

api.GET(func() (interface{}, *cache.ImplCache) {
    data := expensiveDatabaseQuery()
    return data, cache.NewCacheImpl(5 * time.Minute)
}, "/cached")
```

调用此接口时：
1. 框架根据方法名 + 参数值计算缓存键
2. 缓存命中 → 返回缓存数据，**不调用** handler
3. 缓存未命中 → 调用 handler，将序列化结果缓存指定 TTL

## 缓存接口

```go
type Cache interface {
    ExpireTime_() time.Duration
}
```

`ExpireTime_()` 方法返回缓存 TTL。任何实现此接口的类型都可以作为第二个返回值。

## 自定义持久化缓存

默认使用内存中的 `sync.Map` 缓存。替换为自定义实现（如 Redis）：

```go
type RedisCache struct{}

func (r *RedisCache) Set(key []byte, value []byte, ttl time.Duration) {
    // 存储到 Redis
}

func (r *RedisCache) Get(key []byte) []byte {
    // 从 Redis 读取，过期或不存在返回 nil
}

// 在 init 时注册
func init() {
    cache.SetPersistenceCacheImpl(&RedisCache{})
}
```

接口：

```go
type PersistenceCache interface {
    Set(key []byte, value []byte, ttl time.Duration)
    Get(key []byte) []byte
}
```

## 自定义缓存键/值编码

覆盖缓存键和值的编码方式：

```go
type CustomEncoder struct{}

func (c *CustomEncoder) EncodeKey(m *def.MethodInfo, args []reflect.Value) []byte {
    // 返回缓存键字节
}

func (c *CustomEncoder) EncodeValue(s def.Serialize, v reflect.Value) []byte {
    // 返回缓存值字节
}

func init() {
    cache.SetProcessCacheImpl(&CustomEncoder{})
}
```

接口：

```go
type ProcessCache interface {
    EncodeKey(m *def.MethodInfo, args []reflect.Value) []byte
    EncodeValue(s def.Serialize, v reflect.Value) []byte
}
```

## 工作原理

缓存通过**方法代理**（`call.SetMethodProxy`）实现，包装 handler 调用：

```
请求到达
    │
    ├── 计算缓存键（方法名 + 参数值）
    │
    ├── 检查持久化缓存
    │   ├── 命中 → 直接返回缓存字节
    │   └── 未命中 → 调用 handler
    │       ├── 检查第二个返回值是否实现 Cache 接口
    │       ├── 如果存在 → 序列化结果，存入缓存
    │       └── 返回结果
    │
    └── 响应
```

缓存代理是方法调用链的最外层，在其他代理之前执行。
