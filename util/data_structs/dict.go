package util

import "sync"

type Map[K comparable, V any] map[K]V

func (dict Map[K, V]) Keys() []K {
	var keys = make([]K, 0, len(dict))
	for k := range dict {
		keys = append(keys, k)
	}
	return keys
}

func (dict Map[K, V]) Values() []V {
	var vals = make([]V, 0, len(dict))
	for _, v := range dict {
		vals = append(vals, v)
	}
	return vals
}

func (dict Map[K, V]) Flatten() []KeyValPair[K, V] {
	var pairs = make([]KeyValPair[K, V], 0, len(dict))
	for k, v := range dict {
		pairs = append(pairs, KeyValPair[K, V]{k, v})
	}
	return pairs
}

func MakeMap[K comparable, V any](size ...int) Map[K, V] {
	var n = 0
	if len(size) > 1 {
		n = size[0]
	}
	return make(map[K]V, n)
}

func FlattenMap[K comparable, V any](dict map[K]V) []KeyValPair[K, V] {
	var pairs = make([]KeyValPair[K, V], 0)
	for k, v := range dict {
		pairs = append(pairs, KeyValPair[K, V]{k, v})
	}
	return pairs
}

func KeysToMap[K comparable, V any](keys []K, val V) Map[K, V] {
	var dict = make(map[K]V, len(keys))
	for _, k := range keys {
		dict[k] = val
	}
	return dict
}

type SyncMap[K comparable, V any] struct {
	sync.RWMutex
	data map[K]V
}

func MakeSyncMap[K comparable, V any](size ...int) *SyncMap[K, V] {
	var n = 0
	if len(size) > 1 {
		n = size[0]
	}
	return &SyncMap[K, V]{
		data: MakeMap[K, V](n),
	}
}

func (dict *SyncMap[K, V]) Set(key K, value V) {
	dict.Lock()
	defer dict.Unlock()
	dict.data[key] = value
}

func (dict *SyncMap[K, V]) Update(data Map[K, V]) {
	dict.Lock()
	defer dict.Unlock()
	for k, v := range data {
		dict.data[k] = v
	}
}

func (dict *SyncMap[K, V]) Delete(key K) {
	dict.Lock()
	defer dict.Unlock()
	delete(dict.data, key)
}

func (dict *SyncMap[K, V]) Get(key K) (V, bool) {
	dict.RLock()
	defer dict.RUnlock()
	var val, ok = dict.data[key]
	return val, ok
}

func (dict *SyncMap[K, V]) Len() int {
	dict.RLock()
	defer dict.RUnlock()
	return len(dict.data)
}

func (dict *SyncMap[K, V]) ForEach(callback func(k K, v V)) {
	dict.RLock()
	defer dict.RUnlock()

	for k, v := range dict.data {
		callback(k, v)
	}
}
