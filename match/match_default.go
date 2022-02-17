package match

import (
	"gitee.com/fast_api/api/def"
	"gitee.com/fast_api/api/log"
	"net/url"
)

type MatchImpl struct {
	store *store
}

// Match
//if match return func
func (m *MatchImpl) Match(url *url.URL) *def.Entry {
	pv := make([]string, 10)
	data, e := m.store.Get(url.Path, pv)
	if data == nil {
		return nil
	}
	ent := data.(*def.Entry)
	for i := 0; i < len(e); i++ {
		ent.Ids[e[i]] = pv[i]
	}
	log.Debugf("Url path = %s is matched", url.Path)
	return data.(*def.Entry)
}

func (m *MatchImpl) Add(key string, data interface{}) {
	m.store.Add(key, data)
}

func NewMatchImpl() *MatchImpl {
	return &MatchImpl{store: newStore()}
}
