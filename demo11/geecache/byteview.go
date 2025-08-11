package geecache

type BytesView struct {
	b []byte
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}

func (v BytesView) Len() int {
	return len(v.b)
}

func (v BytesView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

func (v BytesView) String() string {
	return string(v.b)
}
