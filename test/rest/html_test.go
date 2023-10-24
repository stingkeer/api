package rest

import (
	"testing"

	"gitee.com/fast_api/api"
)

func TestHtml(t *testing.T) {
	const tpl = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>{{.Title}}</title>
	</head>
	<body>
		{{range .Items}}<div>{{ . }}</div>{{else}}<div><strong>no rows</strong></div>{{end}}
	</body>
</html>`
	api.GET(func() any {
		data := struct {
			Title string
			Items []string
		}{
			Title: "My page",
			Items: []string{
				"My photos",
				"My blog",
			},
		}
		return api.Html(tpl, data)
	}, "/html")
	api.StartService(nil)
}
