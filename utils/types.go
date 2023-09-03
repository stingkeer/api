package utils

import (
	"sync"
)

type Map[K string, V any] struct {
	kv sync.Map
}

func (m *Map[K, V]) Get(name K) (v1 V) {
	if v, b := m.kv.Load(name); b {
		v1 = v.(V)
	}
	return
}

func (m *Map[K, V]) Set(name K, v V) {
	m.kv.Store(name, v)
}

func (m *Map[K, V]) Range(f func(K, V)) {
	m.kv.Range(func(key, value any) bool {
		f(key.(K), value.(V))
		return true
	})
}
