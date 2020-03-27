package api

var _api Api

func SetApi(api Api) {
	_api = api
}

func GetApi() Api {
	if _api == nil {
		SetApi(&APiDefault{map[string]Entry{}})
	}
	return _api
}

type APiDefault struct {
	pools map[string]Entry
}

func (d *APiDefault) getMaps() map[string]Entry {
	return d.pools
}

func (a *APiDefault) GET(f interface{}, url string) {
	a.pools[url] = Entry{
		url:    url,
		method: "GET",
		f:      f,
	}
}
func (a *APiDefault) POST(f interface{}, url string) {
	a.pools[url] = Entry{
		url:    url,
		method: "POST",
		f:      f,
	}
}
func (a *APiDefault) PUT(f interface{}, url string) {
	a.pools[url] = Entry{
		url:    url,
		method: "PUT",
		f:      f,
	}
}
func (a *APiDefault) DELETE(f interface{}, url string) {
	a.pools[url] = Entry{
		url:    url,
		method: "DELETE",
		f:      f,
	}
}
