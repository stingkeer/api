package api

import (
	"github.com/sirupsen/logrus"
	"net/url"
)

type MatchImpl struct {
	pools map[string]Entry
}

/**
if match return func
*/
func (m *MatchImpl) match(url *url.URL, method string) (interface{}, url.Values) {
	logrus.Debugf("url query = %s", url.Query())
	return m.getFuncWithURL(url.Path, method), url.Query()
}

func (m *MatchImpl) getFuncWithURL(url string, method string) interface{} {
	for _url, entry := range m.getMaps() {
		if "/"+_url == url {
			return entry.f
		}
	}
	logrus.Tracef("not match url %s", url)
	return nil
}

func (m *MatchImpl) getMaps() map[string]Entry {
	return m.pools
}
