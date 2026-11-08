package render3d

import "math"

// cacheSize is how many vertices of the recent past a GPU is reckoned to
// remember, not running its vertex shader again for them.
const cacheSize = 32

// inCacheOrder is the triangles of indices, each whole and turned the same
// way, in an order that uses again the vertices the GPU has just run, so it
// runs each as few times as may be: Forsyth's linear-speed ordering. What is
// drawn is the same; only opaque triangles exactly at one depth could tell.
func inCacheOrder(indices []uint32, vertices int) []uint32 {
	tris := len(indices) / 3
	if tris == 0 || vertices == 0 {
		return indices
	}
	o := newOrderer(indices, vertices)
	out := make([]uint32, 0, len(indices))
	next := o.best()
	for len(out) < 3*tris {
		if next < 0 {
			next = o.firstLeft()
		}
		t := indices[3*next : 3*next+3]
		out = append(out, t...)
		next = o.emit(next)
	}
	return out
}

// orderer is Forsyth's ordering in progress: each vertex's triangles yet to
// be drawn, its place in the cache and its score, each triangle's score and
// whether it is drawn, and the cache.
type orderer struct {
	indices []uint32
	uses    [][]int32
	at      []int32
	score   []float64
	tscore  []float64
	done    []bool
	cache   []int32
	scan    int
}

// newOrderer is an orderer of the triangles of indices.
func newOrderer(indices []uint32, vertices int) *orderer {
	tris := len(indices) / 3
	o := &orderer{indices: indices, uses: make([][]int32, vertices), at: make([]int32, vertices), score: make([]float64, vertices),
		tscore: make([]float64, tris), done: make([]bool, tris), cache: make([]int32, 0, cacheSize+3)}
	counts := make([]int32, vertices)
	for _, v := range indices {
		counts[v]++
	}
	all := make([]int32, len(indices))
	from := 0
	for v, c := range counts {
		o.uses[v] = all[from : from : from+int(c)]
		from += int(c)
		o.at[v] = -1
	}
	for i, v := range indices {
		o.uses[v] = append(o.uses[v], int32(i/3))
	}
	for v := range o.score {
		o.score[v] = vertexScore(-1, len(o.uses[v]))
	}
	for t := range o.tscore {
		o.tscore[t] = o.triangleScore(t)
	}
	return o
}

// vertexScore is how good a vertex is to draw next: better near the top of
// the cache, the last three of it a little less so as they are in the
// triangle just drawn, and better with few triangles left, not to strand it.
func vertexScore(place int32, left int) float64 {
	if left == 0 {
		return -1
	}
	s := 0.0
	switch {
	case place < 0:
	case place < 3:
		s = 0.75
	default:
		s = math.Pow(1-float64(place-3)/float64(cacheSize-3), 1.5)
	}
	return s + 2/math.Sqrt(float64(left))
}

// triangleScore is the sum of the scores of t's vertices.
func (o *orderer) triangleScore(t int) float64 {
	return o.score[o.indices[3*t]] + o.score[o.indices[3*t+1]] + o.score[o.indices[3*t+2]]
}

// best is the undrawn triangle of the highest score, the first found of those.
func (o *orderer) best() int {
	at, top := -1, -1.0
	for t, s := range o.tscore {
		if !o.done[t] && s > top {
			at, top = t, s
		}
	}
	return at
}

// firstLeft is the first triangle not yet drawn.
func (o *orderer) firstLeft() int {
	for o.done[o.scan] {
		o.scan++
	}
	return o.scan
}

// emit draws triangle t: its vertices go to the top of the cache and lose
// it from their lists, and the scores of what the cache's change touches are
// worked out anew. It is the best of those to draw next, or -1.
func (o *orderer) emit(t int) int {
	o.done[t] = true
	old := append([]int32(nil), o.cache...)
	fresh := make([]int32, 0, cacheSize+3)
	for k := range 3 {
		v := int32(o.indices[3*t+k])
		o.drop(v, int32(t))
		fresh = append(fresh, v)
	}
	for _, v := range old {
		if v != fresh[0] && v != fresh[1] && v != fresh[2] {
			fresh = append(fresh, v)
		}
	}
	o.cache = fresh
	return o.rescore()
}

// drop takes triangle t off v's list.
func (o *orderer) drop(v, t int32) {
	u := o.uses[v]
	for i, x := range u {
		if x == t {
			u[i] = u[len(u)-1]
			o.uses[v] = u[:len(u)-1]
			return
		}
	}
}

// rescore gives the vertices of the cache their new places and scores, those
// pushed out losing theirs, and the triangles they touch their new scores,
// and is the best of those not drawn, or -1.
func (o *orderer) rescore() int {
	for i, v := range o.cache {
		if i < cacheSize {
			o.at[v] = int32(i)
		} else {
			o.at[v] = -1
		}
		o.score[v] = vertexScore(o.at[v], len(o.uses[v]))
	}
	best, top := -1, -1.0
	for _, v := range o.cache {
		for _, t := range o.uses[v] {
			s := o.triangleScore(int(t))
			o.tscore[t] = s
			if s > top {
				best, top = int(t), s
			}
		}
	}
	if len(o.cache) > cacheSize {
		o.cache = o.cache[:cacheSize]
	}
	return best
}
