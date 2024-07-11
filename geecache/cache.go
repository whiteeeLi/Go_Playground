package geecache

import (
	"geecache/lru"
	"sync"
)

type cache struct {
	mu  sync.Mutex
	lru *lru.Cache
	//lru里面不是有长度么，为什么这里需要额外添加一个储存数据的长度？
	cacheBytes int64
}

func (c *cache) get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		return
	}

	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), ok
	}

	return
}

// 添加数据
func (c *cache) add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		c.lru = lru.New(c.cacheBytes, nil)
	}
	c.lru.Add(key, value)
}
