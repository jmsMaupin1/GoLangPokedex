package pokecache

import (
	"time"
	"sync"
)

type CacheEntry struct {
	createdAt time.Time
	val []byte
}

type Cache struct {
	cacheEntries map[string]CacheEntry
	mu *sync.RWMutex
	interval time.Duration
}

func setInterval(fn func(), tickRate time.Duration) chan bool {	
	ticker := time.NewTicker(tickRate)
	done := make(chan bool)

	go func() {
		select {
		case <-ticker.C:
			fn()
		case <-done:
			ticker.Stop()
			return
		}
	}()

	return done
}

func (c *Cache) reap() {
	currentTime := time.Now()

	for key, cache := range c.cacheEntries {
		timeDelta := currentTime.Sub(cache.createdAt)
		if timeDelta >= c.interval {
			delete(c.cacheEntries, key)
		}
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, ok := c.cacheEntries[key]

	return entry.val, ok
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cacheEntries[key] = CacheEntry{
		createdAt: time.Now(),
		val: val,
	}
}

func NewCache(interval time.Duration) Cache {
	cache := Cache{
		cacheEntries: map[string]CacheEntry{},
		mu: &sync.RWMutex{},
		interval: interval,
	}

	setInterval(cache.reap, 5 * time.Second)

	return cache
}


