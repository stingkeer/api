package swagger

import (
	"embed"
	"fmt"
	stdhttp "net/http"

	"go.aew.app/api.v1/http"
)

//go:embed ui/*
var _static embed.FS

func init() {
	http.DefaultStatic.HandleStatic("/ui/*", "", stdhttp.FS(_static))
	fmt.Printf("\n[swagger ui] http://ip:port/ui/index.html\n\n")
}
