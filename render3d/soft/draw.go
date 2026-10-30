package soft

import (
	"math"

	"github.com/gabrielluizsf/sge/render3d"
)

// vout is a vertex shaded: in clip space, in the world, its normal there,
// and its texture coordinates.
type vout struct {
	clip          [4]float64
	world, normal [3]float64
	uv            [2]float64
}

// Draw draws d's triangles once an instance; a triangle whose indices,
// vertices or matrices are not there is left out.
func (p *pass) Draw(d render3d.Draw) {
	s, ok := p.state()
	if !ok || d.First < 0 || d.Count < 0 {
		return
	}
	idx, verts := p.idx.indices, p.verts.floats
	end := min(d.First+d.Count, len(idx))
	for inst := range max(d.Instances, 1) {
	triangles:
		for t := d.First; t+2 < end; t += 3 {
			var tri [3]vout
			for k := range 3 {
				at := int(idx[t+k]) * s.stride
				if at+s.stride > len(verts) {
					continue triangles
				}
				v, ok := s.vertex(verts[at:at+s.stride], inst)
				if !ok {
					continue triangles
				}
				tri[k] = v
			}
			p.clipped(s, tri)
		}
	}
}

// vertex is the vertex v of instance inst shaded, and whether its
// matrices are there.
func (s *state) vertex(v []float32, inst int) (vout, bool) {
	var m mat
	if s.stride == render3d.SkinnedFloats {
		for k := range 4 {
			j, w := int(v[8+k]), float64(v[12+k])
			if w == 0 {
				continue
			}
			if (j+1)*16 > len(s.objects) {
				return vout{}, false
			}
			for i := range m {
				m[i] += float64(s.objects[j*16+i]) * w
			}
		}
	} else {
		if (inst+1)*16 > len(s.objects) {
			return vout{}, false
		}
		m = matAt(s.objects, inst*16)
	}
	pos := [3]float64{float64(v[0]), float64(v[1]), float64(v[2])}
	w := m.apply(pos, 1)
	n := m.apply([3]float64{float64(v[3]), float64(v[4]), float64(v[5])}, 0)
	out := vout{world: [3]float64{w[0], w[1], w[2]}, normal: unit([3]float64{n[0], n[1], n[2]}), uv: [2]float64{float64(v[6]), float64(v[7])}}
	out.clip = s.vp.apply(out.world, 1)
	return out, true
}

func unit(v [3]float64) [3]float64 {
	l := math.Sqrt(v[0]*v[0] + v[1]*v[1] + v[2]*v[2])
	if l == 0 {
		return v
	}
	return [3]float64{v[0] / l, v[1] / l, v[2] / l}
}
