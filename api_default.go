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
		SetApi(&apiDefault{
			store: newStore(),
		})
	}
	return _api
}

type apiDefault struct {
	fnCaches []*Entry
	store    *store
}

func (d *apiDefault) getStore() *store {
	return d.store
}

func (d *apiDefault) GET(f interface{}, url string) {
	e := &Entry{
		url:    url,
		method: "GET",
		fn:     f,
	}
	d.fnCaches = append(d.fnCaches, e)
	d.store.Add(url, e)
}
func (d *apiDefault) POST(f interface{}, url string) {
	e := &Entry{
		url:    url,
		method: "POST",
		fn:     f,
	}
	d.fnCaches = append(d.fnCaches, e)
	d.store.Add(url, e)
}
func (d *apiDefault) PUT(f interface{}, url string) {
	e := &Entry{
		url:    url,
		method: "PUT",
		fn:     f,
	}
	d.fnCaches = append(d.fnCaches, e)
	d.store.Add(url, e)
}
func (d *apiDefault) DELETE(f interface{}, url string) {
	e := &Entry{
		url:    url,
		method: "DELETE",
		fn:     f,
	}
	d.fnCaches = append(d.fnCaches, e)
	d.store.Add(url, e)
}

func (d *apiDefault) getFnCaches() []*Entry {
	return d.fnCaches
}
