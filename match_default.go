package api

import (
	"fmt"
	"net/url"
)

type MatchImpl struct {
	pools map[string]Entry
}

/**
if match return func
*/
func (m *MatchImpl) match(url *url.URL, method string) interface{} {
	fmt.Println(url)
	return m.getFuncWithURL(url.Path, method)
}

func (m *MatchImpl) getFuncWithURL(url string, method string) interface{} {
	for _url, entry := range m.getMaps() {
		if _url == url {
			return entry.f
		}
	}
	return nil
}

func (m *MatchImpl) getMaps() map[string]Entry {
	return m.pools
}
