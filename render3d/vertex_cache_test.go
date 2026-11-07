package render3d

import (
	"slices"
	"testing"
)

// missRate is how many times a cache of size runs a vertex for each triangle.
func missRate(indices []uint32, size int) float64 {
	cache := make([]uint32, 0, size+1)
	miss := 0
	for _, v := range indices {
		if slices.Contains(cache, v) {
			continue
		}
		miss++
		cache = append(cache, v)
		if len(cache) > size {
			cache = cache[1:]
		}
	}
	return float64(miss) / (float64(len(indices)) / 3)
}

// grid is the triangles of a grid of n by n squares, in an order that keeps
// nothing in the cache: every other row, then the rest.
func grid(n int) ([]uint32, int) {
	var out []uint32
	at := func(c, r int) uint32 { return uint32(r*(n+1) + c) }
	for pass := range 2 {
		for r := pass; r < n; r += 2 {
			for c := range n {
				out = append(out, at(c, r), at(c+1, r), at(c, r+1), at(c+1, r), at(c+1, r+1), at(c, r+1))
			}
		}
	}
	return out, (n + 1) * (n + 1)
}

// sorted is each triangle of indices, its first index the least so a
// triangle reads the same however it was started, keeping its turn, sorted.
func sorted(indices []uint32) [][3]uint32 {
	out := make([][3]uint32, 0, len(indices)/3)
	for i := 0; i+2 < len(indices); i += 3 {
		a, b, c := indices[i], indices[i+1], indices[i+2]
		for a > b || a > c {
			a, b, c = b, c, a
		}
		out = append(out, [3]uint32{a, b, c})
	}
	slices.SortFunc(out, func(x, y [3]uint32) int {
		for k := range 3 {
			if x[k] != y[k] {
				return int(x[k]) - int(y[k])
			}
		}
		return 0
	})
	return out
}

// Ordering keeps every triangle whole and turned the same way, none lost or
// added, and asks far less of the vertex shader.
func TestTheCacheOrderKeepsEveryTriangle(t *testing.T) {
	in, verts := grid(60)
	out := inCacheOrder(in, verts)
	if !slices.Equal(sorted(in), sorted(out)) {
		t.Fatal("the triangles changed")
	}
	before, after := missRate(in, cacheSize), missRate(out, cacheSize)
	if after > 0.75 || after >= before {
		t.Errorf("%.2f vertices a triangle, from %.2f", after, before)
	}
}

// Nothing to order is left alone, and the same input gives the same order.
func TestTheCacheOrderIsTheSameEachTime(t *testing.T) {
	if got := inCacheOrder(nil, 0); len(got) != 0 {
		t.Errorf("%v", got)
	}
	in, verts := grid(20)
	if !slices.Equal(inCacheOrder(in, verts), inCacheOrder(in, verts)) {
		t.Error("two orders")
	}
	// A vertex used by no triangle, and a degenerate one, are no trouble.
	odd := []uint32{0, 1, 2, 2, 2, 3}
	if len(inCacheOrder(odd, 6)) != 6 {
		t.Error("lost a triangle")
	}
}
