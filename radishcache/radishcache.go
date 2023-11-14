package radishcache

import (
	pb "radishcache/radishcachepb"
	"radish/singleflight"
	"fmt"
	"log"
	"sync"
)

// 缓存命名空间
type Group struct {
	name      string              // 唯一的名称
	getter    Getter              // 缓存未命中时获取源数据的回调
	mainCache cache               // 并发缓存
	peers     PeerPicker          // 分布式节点
	loader    *singleflight.Group // 请求组
}

// 回调函数，缓存不存在时，调用这个函数，得到源数据
type Getter interface {
	Get(key string) ([]byte, error)
}

type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

// 实例化Group
func NewGroup(name string, cacacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacacheBytes},
		loader:    &singleflight.Group{},
	}
	//存入全局变量groups中
	groups[name] = g
	return g
}

func GetGroup(name string) *Group {
	mu.RLock()
	g := groups[name]
	mu.RUnlock()
	return g
}

func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}

	if v, ok := g.mainCache.get(key); ok {
		log.Println("[RadishCache] hit")
		return v, nil
	}
	// 缓存未命中 触发回调
	return g.load(key)
}

func (g *Group) getLocally(key string) (ByteView, error) {
	// 调用用户回调函数g.getter.Get()获取源数据
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	value := ByteView{b: cloneBytes(bytes)}
	g.populateCache(key, value)
	return value, nil
}

func (g *Group) populateCache(key string, value ByteView) {
	// 将源数据添加到mainCache
	g.mainCache.add(key, value)
}

// 实现了 PeerPicker 接口的 HTTPPool 注入到 Group 中。
func (g *Group) RegisterPeers(peers PeerPicker) {
	if g.peers != nil {
		panic("RegisterPeerPicker called more than once")
	}
	g.peers = peers
}

func (g *Group) load(key string) (value ByteView, err error) {
	viewi, err := g.loader.Do(key, func() (interface{}, error) {
		if g.peers != nil {
			// 选择节点
			if peer, ok := g.peers.PickPeer(key); ok {
				// 非本机节点，调用getFromPeer从远处获取
				if value, err := g.getFromPeer(peer, key); err == nil {
					return value, nil
				}
				log.Println("[GeeCache] Failed to get from peer", err)
			}
		}
		// 本机节点或者失败，回退到getLocally
		return g.getLocally(key)
	})

	if err == nil {
		return viewi.(ByteView), nil
	}
	return
}

// 实现了PeerGetter接口的httpGetter从访问远程节点，获取缓存值
func (g *Group) getFromPeer(peer PeerGetter, key string) (ByteView, error) {
	req := &pb.Request{
		Group: g.name,
		Key: key,
	}
	res := &pb.Response{}
	err := peer.Get(req, res)
	if err != nil {
		return ByteView{}, err
	}
	return ByteView{b: res.Value}, nil
}
