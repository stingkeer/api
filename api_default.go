package api

import "gitee.com/fast_api/api/public"

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
	fnCaches []*public.Entry
	store    *store
}

func (d *apiDefault) getStore() *store {
	return d.store
}

func (d *apiDefault) GET(f interface{}, url string) {
	e := &public.Entry{
		Url:    url,
		Method: "GET",
		Fn:     f,
		Ids:    make(map[string]string),
	}
	d.fnCaches = append(d.fnCaches, e)
	d.store.Add(url, e)
}
func (d *apiDefault) POST(f interface{}, url string) {
	e := &public.Entry{
		Url:    url,
		Method: "POST",
		Fn:     f,
		Ids:    make(map[string]string),
	}
	d.fnCaches = append(d.fnCaches, e)
	d.store.Add(url, e)
}
func (d *apiDefault) PUT(f interface{}, url string) {
	e := &public.Entry{
		Url:    url,
		Method: "PUT",
		Fn:     f,
		Ids:    make(map[string]string),
	}
	d.fnCaches = append(d.fnCaches, e)
	d.store.Add(url, e)
}
func (d *apiDefault) DELETE(f interface{}, url string) {
	e := &public.Entry{
		Url:    url,
		Method: "DELETE",
		Fn:     f,
		Ids:    make(map[string]string),
	}
	d.fnCaches = append(d.fnCaches, e)
	d.store.Add(url, e)
}

func (d *apiDefault) getFnCaches() []*public.Entry {
	return d.fnCaches
}
