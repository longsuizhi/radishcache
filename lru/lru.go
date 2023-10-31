package lru

import "container/list"

type Cache struct {
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

func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// 添加
func (c *Cache) Add(key string, value Value) {
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
func (c *Cache) RemoveOldest() {
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

// 查找
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

func (c *Cache) Len() int {
	return c.ll.Len()
}
