package swagger

import (
	"embed"
	"fmt"
	stdhttp "net/http"

	"gitee.com/fast_api/api/http"
)

//go:embed ui/*
var _static embed.FS

func init() {
	http.DefaultStatic.HandleStatic("/ui/*", "", stdhttp.FS(_static))
	fmt.Printf("\n[swagger ui] http://ip:port/ui/index.html\n\n")
}
