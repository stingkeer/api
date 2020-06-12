package api

import (
	"github.com/sirupsen/logrus"
	"net/url"
)

type MatchImpl struct {
	store *store
}

/**
if match return func
*/
func (m *MatchImpl) match(url *url.URL) *Entry {
	data, _ := m.store.Get(url.Path, queryToValues(url.Query()))
	if data == nil {
		return nil
	}
	logrus.Debugf("url path = %s is matched", url.Path)
	return data.(*Entry)
}

func queryToValues(v url.Values) []string {
	var strs []string
	for _, strings := range v {
		strs = append(strs, strings[0])
	}
	return strs
}
