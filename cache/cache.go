package cache

import (
	"fmt"
	"sync"
	"time"

	gocache "github.com/patrickmn/go-cache"
)

var (
	ErrCacheIsBeenBuilding = fmt.Errorf("o cache já está sendo efetuado")
)

type Cache struct {
	cache        *gocache.Cache
	buildingLock sync.Mutex
	building     map[string]bool
}

func NewCache() *Cache {
	return &Cache{
		cache:    gocache.New(1*time.Hour, 30*time.Minute),
		building: make(map[string]bool),
	}
}

func (c *Cache) Get(key string, populateFunc func() interface{}) interface{} {
	profilesCache, found := c.cache.Get(key)

	if !found && c.isBuildingCacheFor(key) {
		return ErrCacheIsBeenBuilding
	}

	if !found {
		c.markAsBuilding(key)
		data := populateFunc()
		c.cache.SetDefault(key, data)
		profilesCache = data

		_, isError := profilesCache.(error)
		if isError {
			c.cache.Delete(key)
		}

		c.markAsDoneBuilding(key)
	}

	return profilesCache
}

func (c *Cache) isBuildingCacheFor(key string) bool {
	c.buildingLock.Lock()
	defer c.buildingLock.Unlock()
	_, found := c.building[key]
	return found
}

func (c *Cache) markAsBuilding(key string) {
	c.buildingLock.Lock()
	c.building[key] = true
	c.buildingLock.Unlock()
}

func (c *Cache) markAsDoneBuilding(key string) {
	c.buildingLock.Lock()
	delete(c.building, key)
	c.buildingLock.Unlock()
}
