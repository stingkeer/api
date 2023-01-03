package call

import (
	"gitee.com/fast_api/api/call/rettypes"
	"gitee.com/fast_api/api/call/types"
	"gitee.com/fast_api/api/http"
)

func init() {

	base := types.BaseType{}
	RegisterTypeMapper(&base)
	RegisterTypeMapper(&types.BigType{})
	RegisterTypeMapper(&types.FileType{})
	RegisterTypeMapper(&types.HttpType{})
	RegisterTypeMapper(&types.HeadType{})
	RegisterTypeMapper(&types.TypeRequire{BaseType: base})

	http.RegisterReturnHandler(&rettypes.Stream{})
	http.RegisterReturnHandler(&rettypes.Html{})
	http.RegisterReturnHandler(&rettypes.Redirect{})
}
