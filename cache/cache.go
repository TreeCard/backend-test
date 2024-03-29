package cache

import (
	"container/list"
	"time"
)

// GetterFn loads data for the key and optionally can set the time to live.
type GetterFn func(string) (string, time.Time, bool)

// Cache implements a string key value read-through cache.
//
// The returned cache is a fixed sized with a LRU eviction strategy.
// If a getter function is provided the cache should load the fetched item.
// A TTL for each value can be set, not in the cache it will optionally fetch the items if a getter function is provided.
type Cache struct {
	values map[string]*list.Element
	getter GetterFn
	lru    list.List
	size   int
}

// NewCache creates a limited size cache.
func NewCache(getter GetterFn, size int) *Cache {
	return &Cache{
		size:   10,
		getter: getter,
	}
}

// cacheValue hold the key, value mapping with the optional ttl.
type cacheValue struct {
	key string
	val string
	ttl time.Time
}

func (cv *cacheValue) valid() bool {
	if cv.ttl.IsZero() {
		return true
	}
	return cv.ttl.Before(time.Now())
}

func (c *Cache) remove(e *list.Element) {
	if e == nil {
		return
	}
	cv := e.Value.(cacheValue)
	delete(c.values, cv.key)
}

// GetAt fetches the value. It optionally loads and caches the value if none is
// found or the TTL has passed the given time.
func (c *Cache) GetAt(key string, at time.Time) (string, bool) {
	elem, ok := c.values[key]
	if ok {
		cv := elem.Value.(cacheValue)
		if cv.valid() {
			c.lru.MoveToFront(elem)
			return cv.val, ok
		}

		// Remove non-valid key
		c.remove(elem)
	}

	value, ttl, ok := c.getter(key)
	if !ok {
		return "", false
	}

	l := c.lru.Len()
	if l > c.size {
		c.remove(c.lru.Front())
	}

	elem = c.lru.PushFront(cacheValue{
		key: key,
		val: value,
		ttl: ttl,
	})

	if c.values == nil {
		c.values = make(map[string]*list.Element)
	}
	c.values[key] = elem

	return value, true
}

// Get is a helper for calling GetAt at the current time.
func (c *Cache) Get(key string) (string, bool) {
	return c.GetAt(key, time.Now())
}
