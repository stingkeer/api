package api

import (
	"fmt"
	"net/http"
)

func Start(addr string) {
	apiServer := ApiService{&Match{GetApi().getMaps()}}
	http.ListenAndServe(":8080", &apiServer)
}

type ApiService struct {
	match *Match
}

func (a *ApiService) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	fun := a.match.match(req.URL, req.Method)
	fmt.Println(fun)
}

type Default struct {
	pools map[string]Entry
}

func (d *Default) getMaps() map[string]Entry {
	return d.pools
}

type Entry struct {
	url    string
	group  string
	method string
	f      interface{}
}

func (a *Default) GET(f interface{}, url string) {
	a.pools[url] = Entry{
		url:    url,
		method: "GET",
		f:      f,
	}
}
func (a *Default) POST(f interface{}, url string) {
	a.pools[url] = Entry{
		url:    url,
		method: "POST",
		f:      f,
	}
}
func (a *Default) PUT(f interface{}, url string) {
	a.pools[url] = Entry{
		url:    url,
		method: "PUT",
		f:      f,
	}
}
func (a *Default) DELETE(f interface{}, url string) {
	a.pools[url] = Entry{
		url:    url,
		method: "DELETE",
		f:      f,
	}
}
