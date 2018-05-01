package vhost

import (
	"sync"
	"time"
)

type Cache struct {
	Lock sync.RWMutex
	Data map[string]cacheEntry
}

type cacheEntry struct {
	data *VirtualHost
	ttl  int
}

func NewCache() *Cache {
	return &Cache{
		Data: make(map[string]cacheEntry),
	}
}

func (c *Cache) Set(hostname string, host *VirtualHost, ttl int) {
	c.Lock.Lock()
	defer c.Lock.Unlock()

	c.Data[hostname] = cacheEntry{
		data: host,
		ttl:  ttl,
	}
}

func (c *Cache) Get(hostname string) *VirtualHost {
	c.Lock.RLock()
	defer c.Lock.RUnlock()

	if entry, found := c.Data[hostname]; found && entry.ttl > 0 {
		return entry.data
	}
	return nil
}

func (c *Cache) Run() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		<-ticker.C
		c.Lock.Lock()
		for key, value := range c.Data {
			value.ttl--
			if value.ttl <= 0 {
				delete(c.Data, key)
			}
		}
		c.Lock.Unlock()
	}
}
