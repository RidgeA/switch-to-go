package cache

import (
	"sync"
	"time"
)

type (
	cacheEntry struct {
		data   interface{}
		expire time.Time
	}

	MemoryCache struct {
		sync.RWMutex
		duration time.Duration
		storage  map[string]cacheEntry
	}
)

func NewMemoryCache(duration time.Duration) *MemoryCache {
	//todo gc
	mc := &MemoryCache{}
	mc.storage = make(map[string]cacheEntry)
	mc.duration = duration
	return mc
}

func (mc *MemoryCache) Get(key string) (interface{}, bool) {
	mc.RLock()
	defer mc.RUnlock()

	v, exist := mc.storage[key]
	if !exist {
		return nil, false
	}
	if v.expire.Before(time.Now()) {
		delete(mc.storage, key)
		return nil, false
	}
	return v.data, true
}

func (mc *MemoryCache) Del(key string) {
	mc.Lock()
	defer mc.Unlock()

	delete(mc.storage, key)
}

func (mc *MemoryCache) Set(key string, data interface{}) {
	mc.Lock()
	defer mc.Unlock()

	ce := cacheEntry{
		data,
		time.Now().Add(mc.duration),
	}
	mc.storage[key] = ce
}
