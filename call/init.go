package call

import (
	"go.aew.app/api.v1/call/rettypes"
	"go.aew.app/api.v1/call/types"
	"go.aew.app/api.v1/http"
)

func init() {

	base := types.BaseType{}
	RegisterTypeMapper(&base)
	RegisterTypeMapper(&types.BigType{})
	RegisterTypeMapper(&types.FileType{})
	RegisterTypeMapper(&types.HttpType{})
	RegisterTypeMapper(&types.HeadType{})
	RegisterTypeMapper(&types.WSType{})
	RegisterTypeMapper(&types.TypeRequire{BaseType: base})

	RegisterGenericTypeMapper(&types.TypeRequireG{})

	http.RegisterReturnHandler(&rettypes.Stream{})
	http.RegisterReturnHandler(&rettypes.Html{})
	http.RegisterReturnHandler(&rettypes.Redirect{})
	http.RegisterReturnHandler(&rettypes.Resp{})

}
