package radishcache

import pb "radishcache/radishcachepb"

type PeerPicker interface {
	// 用于根据传入的key选择相应节点PeerGeeter
	PickPeer(key string) (peer PeerGetter, ok bool)
}

type PeerGetter interface {
	// 用于从对应group查找缓存值
	Get(in *pb.Request, out *pb.Response) error
}