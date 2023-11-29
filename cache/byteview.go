package cache

type ByteView struct {
	B []byte //存储真实的缓存值 只读
}

func (v ByteView) Len() int {
	return len(v.B)
}

func (v ByteView) ByteSlice() []byte {
	return CloneBytes(v.B) // b 只读 返回一个对b的拷贝，防止缓存值呗外部程序修改
}

func (v ByteView) String() string {
	return string(v.B)
}

func CloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
