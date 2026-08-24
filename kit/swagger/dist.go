package swagger

import (
	"embed"
	"io/fs"
	stdhttp "net/http"

	"go.aew.app/api.v1/http"
)

//go:embed ui/*
var _static embed.FS

func init() {
	// embed keeps the "ui/" prefix in file names (ui/index.html), so expose
	// a sub-FS rooted at ui/ and strip the /ui/ URL prefix via rewrite:
	// GET /ui/index.html → sub-FS "index.html". Without both halves every
	// UI request 404s.
	sub, err := fs.Sub(_static, "ui")
	if err != nil {
		panic("swagger ui: embedded fs missing ui/: " + err.Error())
	}
	http.DefaultStatic.HandleStatic("/ui/*", "", stdhttp.FS(sub),
		http.StaticRewrite("/ui/", ""),
	)
}
