package cache

type ByteView struct {
	b []byte //存储真实的缓存值 只读
}

func (v ByteView) Len() int {
	return len(v.b)
}

func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b) // b 只读 返回一个对b的拷贝，防止缓存值呗外部程序修改
}

func (v ByteView) String() string {
	return string(v.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
