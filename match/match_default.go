package match

import (
	"net/url"
	"sync"

	"go.aew.app/api.v1/def"
	"go.aew.app/api.v1/log"
)

type MatchImpl struct {
	store *store
}

// pvaluesPool recycles the param-value scratch slice across requests. The
// previous make([]string, 10) per request was pure hot-path garbage.
var pvaluesPool = sync.Pool{
	New: func() any {
		s := make([]string, 16)
		return &s
	},
}

// Match
// if match return func
func (m *MatchImpl) Match(url *url.URL) *def.Entry {
	pvP := pvaluesPool.Get().(*[]string)
	pv := *pvP
	data, e := m.store.Get(url.Path, pv)
	if data == nil {
		clear(pv)
		pvaluesPool.Put(pvP)
		return nil
	}
	ent := data.(*def.Entry)
	for i := 0; i < len(e); i++ {
		ent.Ids.Set(e[i], pv[i])
	}
	// The captured strings are immutable; only clear the scratch slice
	// before recycling it back to the pool.
	clear(pv)
	pvaluesPool.Put(pvP)
	log.Debugf("Url path = %s is matched", url.Path)
	return ent
}

func (m *MatchImpl) Add(key string, data interface{}) {
	m.store.Add(key, data)
}

func NewMatchImpl() *MatchImpl {
	return &MatchImpl{store: newStore()}
}
