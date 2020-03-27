package api

type Api interface {
	GET(interface{}, string)
	POST(interface{}, string)
	PUT(interface{}, string)
	DELETE(interface{}, string)
	getMaps() map[string]Entry
}

var _api Api

func SetApi(api Api) {
	_api = api
}

func GetApi() Api {
	if _api == nil {
		SetApi(&Default{map[string]Entry{}})
	}
	return _api
}
