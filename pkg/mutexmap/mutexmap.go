package runtime

import (
	"sync"
)

type MutexMap[K comparable, V any] struct {
	mp map[K]V
	mu sync.RWMutex
}

func NewMutexMapFilled[K comparable, V any](mp map[K]V) *MutexMap[K, V] {
	return &MutexMap[K, V]{
		mp: mp,
	}
}

func (m *MutexMap[K, V]) GetOK(k K) (V, bool) {
	m.mu.RLock()
	v, ok := m.mp[k]
	m.mu.RUnlock()
	return v, ok
}

func (m *MutexMap[K, V]) GetSomeValues(count int) []V {
	m.mu.RLock()

	result := make([]V, 0, count)
	for _, v := range m.mp {
		result = append(result, v)
	}

	m.mu.RUnlock()
	return result
}

func (m *MutexMap[K, V]) Set(k K, v V) {
	m.mu.Lock()
	m.mp[k] = v
	m.mu.Unlock()
}

func (m *MutexMap[K, V]) DeleteIfExist(k K) (isDeleted bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.mp[k]; !ok {
		return false
	}
	delete(m.mp, k)
	return true
}
