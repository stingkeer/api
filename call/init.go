package call

import (
	"gitee.com/fast_api/api/call/rettypes"
	"gitee.com/fast_api/api/call/types"
	"gitee.com/fast_api/api/http"
)

func init() {

	RegisterTypeMapper(&types.BaseType{})
	RegisterTypeMapper(&types.BigType{})
	RegisterTypeMapper(&types.FileType{})
	RegisterTypeMapper(&types.HttpType{})
	RegisterTypeMapper(&types.HeadType{})

	http.RegisterReturnHandler(&rettypes.Stream{})
}
