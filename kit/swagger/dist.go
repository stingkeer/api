package swagger

import (
	"embed"
	"fmt"
	"gitee.com/fast_api/api/http"
	stdhttp "net/http"
)

//go:embed ui/*
var _static embed.FS

func init() {
	http.DefaultStatic.HandleStatic("/ui/*", "", stdhttp.FS(_static))
	fmt.Println("swagger ui http://ip:port/ui/index.html")
}
