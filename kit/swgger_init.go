//go:build swagger

package kit

import (
	"net/http"

	"go.aew.app/api.v1/def"
	"go.aew.app/api.v1/kit/core"
	"go.aew.app/api.v1/kit/swagger"
)

func init() {
	core.HttpM(http.MethodGet, def.DefaultContext)(func() swagger.Swagger {
		return swagger.GenSwagger(def.DefaultContext)
	}, "/api/swagger").Swagger(func(swagger def.SwaggerOps) {
		swagger.SetSummary("swagger json data")
		swagger.SetDescription("auto gen from api.")
	})
}
