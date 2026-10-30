package gltf

import (
	"encoding/binary"
	"math"
)

// accessor is an accessor's elements, checked to lie in their view:
// count of them, comps components each, stride bytes apart.
type accessor struct {
	data         []byte
	count, comps int
	stride, size int
	component    int
	normalized   bool
}

// The component types, and their sizes in bytes.
var componentSize = map[int]int{5120: 1, 5121: 1, 5122: 2, 5123: 2, 5125: 4, 5126: 4}

// typeComps are the components of each accessor type.
var typeComps = map[string]int{"SCALAR": 1, "VEC2": 2, "VEC3": 3, "VEC4": 4, "MAT4": 16}

// accessor is accessor i, of one of the types allowed, checked: its type
// and components known, not sparse, and every element inside its buffer
// view, the sizes worked out without overflow.
func (v views) accessor(i int, allowed ...string) (accessor, error) {
	if i < 0 || i >= len(v.doc.Accessors) {
		return accessor{}, malformed("accessor %d of %d", i, len(v.doc.Accessors))
	}
	a := v.doc.Accessors[i]
	comps, ok := typeComps[a.Type]
	csize, okc := componentSize[a.ComponentType]
	if !ok || !okc || !contains(allowed, a.Type) {
		return accessor{}, malformed("accessor %d is %s of component %d", i, a.Type, a.ComponentType)
	}
	if a.Sparse != nil {
		return accessor{}, malformed("accessor %d is sparse, which is not read", i)
	}
	if a.BufferView == nil || a.Count < 0 || a.ByteOffset < 0 {
		return accessor{}, malformed("accessor %d has no data", i)
	}
	data, stride, err := v.view(*a.BufferView)
	if err != nil {
		return accessor{}, err
	}
	size := int64(comps * csize)
	if stride == 0 {
		stride = size
	}
	if stride < size {
		return accessor{}, malformed("accessor %d's elements overlap", i)
	}
	if a.Count > 0 {
		last, ok1 := mulInt(a.Count-1, stride)
		end, ok2 := addInt(last, size)
		end, ok3 := addInt(end, a.ByteOffset)
		if !ok1 || !ok2 || !ok3 || end > int64(len(data)) {
			return accessor{}, malformed("accessor %d past its view", i)
		}
	}
	if a.ByteOffset > int64(len(data)) {
		return accessor{}, malformed("accessor %d past its view", i)
	}
	return accessor{data: data[a.ByteOffset:], count: int(a.Count), comps: comps, stride: int(stride), size: int(size),
		component: a.ComponentType, normalized: a.Normalized}, nil
}

func contains(list []string, s string) bool {
	for _, x := range list {
		if x == s {
			return true
		}
	}
	return false
}

// float is component c of element e, normalized to 0..1 or -1..1 when the
// accessor says so.
func (a accessor) float(e, c int) float64 {
	at := e*a.stride + c*componentSize[a.component]
	d := a.data[at:]
	switch a.component {
	case 5126:
		return float64(math.Float32frombits(binary.LittleEndian.Uint32(d)))
	case 5120:
		return a.scale(float64(int8(d[0])), 127, true)
	case 5121:
		return a.scale(float64(d[0]), 255, false)
	case 5122:
		return a.scale(float64(int16(binary.LittleEndian.Uint16(d))), 32767, true)
	case 5123:
		return a.scale(float64(binary.LittleEndian.Uint16(d)), 65535, false)
	}
	return a.scale(float64(binary.LittleEndian.Uint32(d)), 4294967295, false)
}

func (a accessor) scale(x, max float64, signed bool) float64 {
	if !a.normalized {
		return x
	}
	if signed {
		return math.Max(x/max, -1)
	}
	return x / max
}

// uint is component c of element e as an unsigned integer; for an
// integer accessor only.
func (a accessor) uint(e, c int) uint32 {
	at := e*a.stride + c*componentSize[a.component]
	d := a.data[at:]
	switch a.component {
	case 5121, 5120:
		return uint32(d[0])
	case 5123, 5122:
		return uint32(binary.LittleEndian.Uint16(d))
	}
	return binary.LittleEndian.Uint32(d)
}

// integer reports whether the accessor holds unsigned integers, as
// indices and joints must.
func (a accessor) integer() bool {
	return (a.component == 5121 || a.component == 5123 || a.component == 5125) && !a.normalized
}
