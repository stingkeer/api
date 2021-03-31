package match

import (
	"gitee.com/fast_api/api/public"
	"github.com/sirupsen/logrus"
	"net/url"
)

type MatchImpl struct {
	store *store
}

/**
if match return func
*/
func (m *MatchImpl) Match(url *url.URL) *public.Entry {
	pv := make([]string, 10)
	data, e := m.store.Get(url.Path, pv)
	if data == nil {
		return nil
	}
	ent := data.(*public.Entry)
	for i := 0; i < len(e); i++ {
		ent.Ids[e[i]] = pv[i]
	}
	logrus.Debugf("Url path = %s is matched", url.Path)
	return data.(*public.Entry)
}

func (m *MatchImpl) Add(key string, data interface{}) {
	m.store.Add(key, data)
}

func NewMatchImpl() *MatchImpl {
	return &MatchImpl{store: newStore()}
}
