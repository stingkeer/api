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
	pv := make([]string, 10)
	data, e := m.store.Get(url.Path, pv)
	if data == nil {
		return nil
	}
	ent := data.(*Entry)
	for i := 0; i < len(e); i++ {
		ent.ids[e[i]] = pv[i]
	}
	logrus.Debugf("url path = %s is matched", url.Path)
	return data.(*Entry)
}
