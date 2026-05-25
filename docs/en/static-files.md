# Static Files

Serve embedded static files using Go's `embed.FS`.

## Basic Usage

```go
//go:embed public
var public embed.FS

func init() {
    api.Static("/admin/*", "public", http.FS(public))
}
```

Parameters:
1. **URL path pattern** — supports `*` wildcard
2. **Directory path** — path within the filesystem
3. **File system** — `http.FS(embed.FS)` or any `http.FileSystem`

## Path Rewriting

Rewrite URL paths before file lookup:

```go
api.Static("/admin/*", "public", http.FS(public),
    api.StaticRewrite("/admin/", ""),
)
```

This strips `/admin/` from the URL before looking up the file:
- Request `/admin/style.css` → serves `public/style.css`
- Request `/admin/js/app.js` → serves `public/js/app.js`

## Default File (SPA Support)

Set a default file for unmatched paths — useful for Single Page Applications:

```go
api.Static("/admin/*", "public", http.FS(public),
    api.StaticRewrite("/admin/", ""),
    api.StaticDefaultFile("index.html"),
)
```

When a file is not found, the framework serves `index.html` instead, allowing client-side routing to handle the URL.

## Multiple Static Entries

Register multiple static file handlers for different paths:

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

## How It Works

Static file serving is implemented as an HTTP interceptor with `Order = 99`, running before the API handler. When a static file is found:
1. The response is written directly using `http.ServeContent`
2. `SkipResponse()` is called to prevent the response pipeline from running

If no file matches, the request continues to the API handler.

## Root Path

When the URL path is `/`, the framework automatically serves `index.html`:

```go
api.Static("/*", "public", http.FS(public))
```
