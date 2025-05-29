package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	cache map[string]cacheEntry
	mu    *sync.RWMutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entery, exist := c.cache[key]
	if !exist {
		return nil, false
	}
	return entery.val, true
}

func (c *Cache) reapLoop(ticker *time.Ticker, interval time.Duration) {
	for {
		select {
		case <-ticker.C:
			c.mu.Lock()
			for key, entery := range c.cache {
				if time.Since(entery.createdAt) > interval {
					delete(c.cache, key)
				}
			}
			c.mu.Unlock()
		}
	}
}

func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		cache: make(map[string]cacheEntry),
		mu:    &sync.RWMutex{},
	}
	ticker := time.NewTicker(interval)

	go c.reapLoop(ticker, interval)

	return c
}
