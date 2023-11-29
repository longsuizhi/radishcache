package cache

import (
	"sync"

)

const (
	TYPE_LRU    = "lru"
	TYPE_LFU    = "lfu"
	TYPE_ARC    = "arc"
)

type baseCache struct {
	mu         sync.Mutex // 互斥锁
	cacheBytes int64
	tp         string
}

type Cache interface {
	Add(key string, value ByteView)
	Get(key string) (value ByteView, ok bool)
	Len() int
}

func New(cacacheBytes int64, strategy string) Cache {
	switch strategy {
	case TYPE_LRU:
		return NewLRUCache(cacacheBytes, nil)
	default:
		panic("radishcache: Unknown type " + strategy)
	}
}