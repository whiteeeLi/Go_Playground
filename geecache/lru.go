package geecache

import "container/list"

type Cache struct {
	//多考虑了Cache的两个参数maxBytes 是允许使用的最大内存，nbytes 是当前已使用的内存
	maxBytes int64
	nbytes   int64
	// 这里队列使用的是指针，合理，直接使用队列太大了
	ll *list.List
	// 使用队列中定义的Element
	cache map[string]*list.Element
	//  是某条记录被移除时的回调函数，可以为 nil。
	OnEvicted func(key string, value Value)
}

type entry struct {
	key   string
	value Value
}

// Value use Len to count how many bytes it takes
type Value interface {
	Len() int
}

// New 构造函数
func New(maxBytes int64, onEvicted func(key string, value Value)) *Cache {
	result := &Cache{
		maxBytes:  maxBytes,
		nbytes:    0,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
	return result
}

// Get 用于查询
func (c *Cache) Get(key string) (value Value, ok bool) {
	//如果查询到则执行if语句中的内容
	if ele, ok := c.cache[key]; ok {
		//将对应元素放到队尾，这里以Front为队尾
		c.ll.MoveToFront(ele)
		//这里其实要和插入的函数相互呼应，插入的时候element中的Value是什么
		//这里就因该提取出什么
		//这里有一个类型断言
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	//没查询到直接返回空值
	return nil, false
}

// Add 新增
func (c *Cache) Add(key string, value Value) {
	//en := entry{key, value}
	//检查字典，如果已经存在key则进行修改
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		//更新Cache储存的长度数值
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		//没找到的话就是新增
		//插入队尾
		ele := c.ll.PushFront(&entry{key, value})
		//插入map
		c.cache[key] = ele
		c.nbytes += int64(len(key)) + int64(value.Len())
	}
	//当cache容量满了，需要删除队尾的元素保证cache有空的地方
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}
}

// RemoveOldest removes the oldest item
func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		//学到了一个新的用法，使用delete删除map中的元素
		delete(c.cache, kv.key)
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// Len the number of cache entries
func (c *Cache) Len() int {
	return c.ll.Len()
}
