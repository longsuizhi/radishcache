package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// 采取依赖注入的方式，允许用于替换成自定义的Hash函数，默认crc32.ChecksumIEEE算法
type Hash func(data []byte) uint32

type Map struct {
	hash     Hash           // hash函数
	replicas int            // 虚拟节点倍数
	keys     []int          // 哈希环
	hashMap  map[int]string // 虚拟节点与真实节点的映射表 key是虚拟节点的hash值 值时真实节点的名称
}

func New(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}
	// 没有自定义hash算法则默认crc32.ChecksumIEEE算法
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// 允许传入0或多个真实节点的名称
func (m *Map) Add(keys ...string) {
	// 对每个真实节点创建m.replicas个虚拟节点
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			// 虚拟节点名称：strconv.Itoa(i) + key
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			// 将虚拟节点添加到hash环上
			m.keys = append(m.keys, hash)
			// 添加虚拟节点和真实节点的映射关系
			m.hashMap[hash] = key
		}
	}
	// 将环上hash值排序
	sort.Ints(m.keys)
}

func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return " "
	}
	// 计算key的哈希值
	hash := int(m.hash([]byte(key)))
	// 顺时针找到第一个匹配的虚拟节点的下标index
	index := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})
	// 从m.key中获取对应哈希值
	return m.hashMap[m.keys[index%len(m.keys)]]
}