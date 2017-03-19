package main

import (
	"fmt"
	"sync"

)


// Cache is thread-safe in-memory state manager.
type Cache struct {
	state map[string]interface{}
	mutex *sync.RWMutex
}

// Set will insert a new key-value record to state map.
func (c *Cache) Set(key string, val interface{}) {
	c.mutex.Lock()
	c.state[key] = val
	c.mutex.Unlock()
}

// Get will retrieve value under corresponding key from state map.
func (c *Cache) Get(key string) (interface{}, error) {
	c.mutex.RLock()
	val := c.state[key]
	c.mutex.RUnlock()
	if val != nil {
		return val, nil
	}
	return nil, fmt.Errorf(KeyNotFoundError)
}

// NewCache constructs empty cache object.
func NewCache() *Cache {
	return &Cache{
		state: make(map[string]interface{}),
		mutex: &sync.RWMutex{},
	}
}

var singleCache *Cache
var once sync.Once

// SingleCache constructs a thread-safe singleton cache object.
func SingleCache() *Cache {
	once.Do(func() { singleCache = NewCache() })
	return singleCache
}
