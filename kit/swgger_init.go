//go:build swagger

package kit

import (
	"gitee.com/fast_api/api/def"
	"gitee.com/fast_api/api/kit/core"
	"gitee.com/fast_api/api/kit/swagger"
	"net/http"
)

func init() {
	core.HttpM(http.MethodGet, def.DefaultContext)(func() swagger.Swagger {
		return swagger.GenSwagger(def.DefaultContext)
	}, "/api/swagger").Swagger(func(swagger def.SwaggerOps) {
		swagger.SetSummary("swagger json data")
		swagger.SetDescription("auto gen from api.")
	})
}
