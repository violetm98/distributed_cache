package distributedCache

import "container/list"

type Cache struct {
	maxBytes  int64      //max capacity in the cache
	usedBytes int64      //currently taken bytes
	ll        *list.List //double linked list
	cache     map[string]*list.Element
	OnEvicted func(key string, value Value) //triggered when an item is evicted from the cache
}

type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}

// New is to create cache structure
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToBack(ele)
		kv := ele.Value.(*entry)
		c.usedBytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		ele := c.ll.PushBack(&entry{key, value})
		c.cache[key] = ele
		c.usedBytes += int64(len(key)) + int64(value.Len())
	}
	for c.maxBytes != 0 && c.maxBytes < c.usedBytes {
		c.RemoveOldest()
	}
}

// Get key's value if it exists
func (c *Cache) Get(key string) (value Value, ok bool) {
	element, ok := c.cache[key]
	if ok {
		c.ll.MoveToBack(element)
		kv := element.Value.(*entry)
		return kv.value, true
	}
	return
}

func (c *Cache) RemoveOldest() {
	element := c.ll.Front()
	if element != nil {
		c.ll.Remove(element)
		kv := element.Value.(*entry)
		delete(c.cache, kv.key)
		c.usedBytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

func (c *Cache) Len() int {
	return c.ll.Len()
}
