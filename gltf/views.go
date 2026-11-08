package gltf

import "math"

// views are a model's buffers, as read, and its buffer views into them.
type views struct {
	buffers [][]byte
	doc     *document
}

// view is the bytes of buffer view i, and its stride (0 for packed).
func (v views) view(i int) ([]byte, int64, error) {
	if i < 0 || i >= len(v.doc.BufferViews) {
		return nil, 0, malformed("buffer view %d of %d", i, len(v.doc.BufferViews))
	}
	bv := v.doc.BufferViews[i]
	if bv.Buffer < 0 || bv.Buffer >= len(v.buffers) {
		return nil, 0, malformed("buffer %d of %d", bv.Buffer, len(v.buffers))
	}
	buf := v.buffers[bv.Buffer]
	end, ok := addInt(bv.ByteOffset, bv.ByteLength)
	if bv.ByteOffset < 0 || bv.ByteLength < 0 || !ok || end > int64(len(buf)) {
		return nil, 0, malformed("buffer view %d past its buffer", i)
	}
	if bv.ByteStride != 0 && (bv.ByteStride < 4 || bv.ByteStride > 252 || bv.ByteStride%4 != 0) {
		return nil, 0, malformed("buffer view %d stride %d", i, bv.ByteStride)
	}
	return buf[bv.ByteOffset:end], bv.ByteStride, nil
}

// addInt and mulInt are a+b and a*b of sizes, and false when either is
// negative or the result would overflow.
func addInt(a, b int64) (int64, bool) {
	if a < 0 || b < 0 || a > math.MaxInt64-b {
		return 0, false
	}
	return a + b, true
}

func mulInt(a, b int64) (int64, bool) {
	if a < 0 || b < 0 || a != 0 && b > math.MaxInt64/a {
		return 0, false
	}
	return a * b, true
}
