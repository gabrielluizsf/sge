package building

import (
	"math"
	"testing"

	"github.com/gabrielluizsf/sge/gltf"
)

// Every building is made of closed solids facing out: each edge one
// triangle runs along is run back along by another, and each piece the
// edges join encloses a volume of its own, its faces outward.
func TestBuildingsAreClosedFacingOut(t *testing.T) {
	for _, k := range All() {
		m := k.Model(Plain())
		edges := map[[2]point]int{}
		for _, tr := range triangles(m) {
			for e := range 3 {
				edges[[2]point{tr[e], tr[(e+1)%3]}]++
			}
		}
		for e, n := range edges {
			if back := edges[[2]point{e[1], e[0]}]; back != n {
				t.Fatalf("%s: the edge %v is run %d times, back %d", k.Name, e, n, back)
			}
		}
		for i, v := range volumes(m) {
			if v <= 0 {
				t.Errorf("%s: piece %d encloses %v m³, facing in", k.Name, i, v)
			}
		}
	}
}

// point is a point of a model, to a tenth of a millimetre.
type point [3]int64

func pointOf(v [3]float32) point {
	q := func(x float32) int64 { return int64(math.Round(float64(x) * 1e4)) }
	return point{q(v[0]), q(v[1]), q(v[2])}
}

// triangles are every triangle of m, its corners as points.
func triangles(m *gltf.Model) [][3]point {
	var out [][3]point
	for _, p := range m.Meshes[0].Primitives {
		for k := 0; k+2 < len(p.Indices); k += 3 {
			out = append(out, [3]point{pointOf(p.Positions[p.Indices[k]]), pointOf(p.Positions[p.Indices[k+1]]), pointOf(p.Positions[p.Indices[k+2]])})
		}
	}
	return out
}

// volumes are the volumes the pieces of m enclose, a piece being the
// triangles joined by their corners, each signed: positive facing out.
func volumes(m *gltf.Model) []float64 {
	tris := triangles(m)
	parent := map[point]point{}
	var find func(p point) point
	find = func(p point) point {
		if q, ok := parent[p]; ok && q != p {
			r := find(q)
			parent[p] = r
			return r
		}
		parent[p] = p
		return p
	}
	for _, tr := range tris {
		parent[find(tr[1])] = find(tr[0])
		parent[find(tr[2])] = find(tr[0])
	}
	// Each piece's volume is summed about a corner of its own, as built,
	// so a rope's is not lost in the rounding of a house's.
	by, about := map[point]float64{}, map[point]vec{}
	for _, p := range m.Meshes[0].Primitives {
		for k := 0; k+2 < len(p.Indices); k += 3 {
			a := f64(p.Positions[p.Indices[k]])
			root := find(pointOf(p.Positions[p.Indices[k]]))
			o, ok := about[root]
			if !ok {
				o, about[root] = a, a
			}
			b, c := f64(p.Positions[p.Indices[k+1]]), f64(p.Positions[p.Indices[k+2]])
			by[root] += dot(sub(a, o), cross(sub(b, o), sub(c, o))) / 6
		}
	}
	var out []float64
	for _, v := range by {
		out = append(out, v)
	}
	return out
}
