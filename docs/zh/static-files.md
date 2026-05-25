# 静态文件服务

使用 Go 的 `embed.FS` 嵌入并提供静态文件服务。

## 基本用法

```go
//go:embed public
var public embed.FS

func init() {
    api.Static("/admin/*", "public", http.FS(public))
}
```

参数说明：
1. **URL 路径模式** — 支持 `*` 通配符
2. **目录路径** — 文件系统内的路径
3. **文件系统** — `http.FS(embed.FS)` 或任意 `http.FileSystem`

## 路径重写

在查找文件前重写 URL 路径：

```go
api.Static("/admin/*", "public", http.FS(public),
    api.StaticRewrite("/admin/", ""),
)
```

从 URL 中去除 `/admin/` 后查找文件：
- 请求 `/admin/style.css` → 提供 `public/style.css`
- 请求 `/admin/js/app.js` → 提供 `public/js/app.js`

## 默认文件（SPA 支持）

设置未匹配路径的默认文件 — 适用于单页应用：

```go
api.Static("/admin/*", "public", http.FS(public),
    api.StaticRewrite("/admin/", ""),
    api.StaticDefaultFile("index.html"),
)
```

当文件未找到时，框架提供 `index.html`，让客户端路由处理 URL。

## 多个静态文件入口

为不同路径注册多个静态文件处理器：

```go
//go:embed public
var public embed.FS

//go:embed assets
var assets embed.FS

func init() {
    api.Static("/admin/*", "public", http.FS(public),
        api.StaticRewrite("/admin/", ""),
        api.StaticDefaultFile("index.html"),
    )
    api.Static("/static/*", "assets", http.FS(assets))
}
```

## 工作原理

静态文件服务实现为 `Order = 99` 的 HTTP 拦截器，在 API handler 之前运行。文件找到时：
1. 使用 `http.ServeContent` 直接写入响应
2. 调用 `SkipResponse()` 阻止响应管道继续执行

如果没有文件匹配，请求继续传递给 API handler。

## 根路径

URL 路径为 `/` 时，框架自动提供 `index.html`：

```go
api.Static("/*", "public", http.FS(public))
```
