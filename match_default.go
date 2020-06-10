package api

import (
	"github.com/sirupsen/logrus"
	"net/url"
	"strings"
)

type MatchImpl struct {
	pools map[string]Entry
}

/**
if match return func
*/
func (m *MatchImpl) match(url *url.URL) interface{} {
	logrus.Debugf("url query = %s", url.Query())
	return m.getFuncWithURL(url.Path)
}

func (m *MatchImpl) getFuncWithURL(url string) interface{} {
	for _url, entry := range m.getMaps() {
		str := strings.ReplaceAll("/"+_url, "//", "/")
		if str == url {
			return entry.fn
		}
	}
	logrus.Tracef("not match url %s", url)
	return nil
}

func (m *MatchImpl) getMaps() map[string]Entry {
	return m.pools
}
