package api

type Entry struct {
	url    string
	group  string
	method string
	f      interface{}
}

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

func (d *APiDefault) GET(f interface{}, url string) {
	d.pools[url] = Entry{
		url:    url,
		method: "GET",
		f:      f,
	}
}
func (d *APiDefault) POST(f interface{}, url string) {
	d.pools[url] = Entry{
		url:    url,
		method: "POST",
		f:      f,
	}
}
func (d *APiDefault) PUT(f interface{}, url string) {
	d.pools[url] = Entry{
		url:    url,
		method: "PUT",
		f:      f,
	}
}
func (d *APiDefault) DELETE(f interface{}, url string) {
	d.pools[url] = Entry{
		url:    url,
		method: "DELETE",
		f:      f,
	}
}
