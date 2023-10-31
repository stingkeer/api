package rest

import (
	"testing"

	"gitee.com/fast_api/api"
	"gitee.com/fast_api/api/def"
	"gitee.com/fast_api/api/test/r"
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
	r.Test(t, func() def.Option {
		return api.GET(func() any {
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
	}).Request().Do(func(resp *r.Response) {
		if resp.Header(def.Content_Type) != def.CONTENT_HTML {
			t.Error("not html Content_Type")
		}
		t.Log(resp.BodyString())
	})
}
