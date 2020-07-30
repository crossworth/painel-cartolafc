package cache

import (
	"time"

	gocache "github.com/patrickmn/go-cache"
)

type Cache struct {
	cache *gocache.Cache
}

func NewCache() *Cache {
	return &Cache{
		cache: gocache.New(5*time.Hour, 1*time.Hour),
	}
}

func (c *Cache) Get(key string, populateFunc func() interface{}) interface{} {
	profilesCache, found := c.cache.Get(key)

	if !found {
		data := populateFunc()
		c.cache.SetDefault(key, data)
		profilesCache = data
	}

	return profilesCache
}
