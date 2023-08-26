package utils

import "sync"

type Map[K string, V any] struct {
	kv sync.Map
}

func (m *Map[K, V]) Get(name K) *V {
	if v, b := m.kv.Load(name); b {
		return v.(*V)
	} else {
		return nil
	}
}

func (m *Map[K, V]) Set(name K, methodInfo *V) {
	m.kv.Store(name, methodInfo)
}

func (m *Map[K, V]) Range(f func(K, *V)) {
	m.kv.Range(func(key, value any) bool {
		f(key.(K), value.(*V))
		return true
	})
}
