package cache

import (
	"container/list"
)

type LRUCache struct {
	baseCache
	maxBytes  int64 //允许使用的最大内存
	nbytes    int64 //当前已使用的内存
	ll        *list.List
	cache     map[string]*list.Element
	OnEvicted func(key string, value Value) //某条记录被移除时的回调函数
}

type entry struct {
	key   string
	value Value
}

// 返回值所占用的内存大小
type Value interface {
	Len() int
}

func NewLRUCache(maxBytes int64, onEvicted func(string, Value)) *LRUCache {
	return &LRUCache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

func (c *LRUCache) Add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Add(key, value)
}

// 添加
func (c *LRUCache) add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		// 键存在 更新对应节点的值 移到队尾
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		// 键不存在 在队尾添加新节点 在字典中添加映射关系
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nbytes += int64(len(key)) + int64(value.Len())
	}
	// 更新当前已使用内存 超过最大值需要移除最少访问节点
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}
}

// 移除最少访问节点
func (c *LRUCache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

func (c *LRUCache) Get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if v, ok := c.get(key); ok {
		return v.(ByteView), ok
	}
	return
}

// 查找
func (c *LRUCache) get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

func (c *LRUCache) Len() int {
	return c.ll.Len()
}
