package api

type Entry struct {
	url    string
	group  string
	method string
	fn     interface{}
}

var _api *apiDefault

func SetApi(api *apiDefault) {
	_api = api
}

func GetApi() *apiDefault {
	if _api == nil {
		SetApi(&apiDefault{map[string]Entry{}})
	}
	return _api
}

type apiDefault struct {
	pools map[string]Entry
}

func (d *apiDefault) getMaps() map[string]Entry {
	return d.pools
}

func (d *apiDefault) GET(f interface{}, url string) {
	d.pools[url] = Entry{
		url:    url,
		method: "GET",
		fn:     f,
	}
}
func (d *apiDefault) POST(f interface{}, url string) {
	d.pools[url] = Entry{
		url:    url,
		method: "POST",
		fn:     f,
	}
}
func (d *apiDefault) PUT(f interface{}, url string) {
	d.pools[url] = Entry{
		url:    url,
		method: "PUT",
		fn:     f,
	}
}
func (d *apiDefault) DELETE(f interface{}, url string) {
	d.pools[url] = Entry{
		url:    url,
		method: "DELETE",
		fn:     f,
	}
}
