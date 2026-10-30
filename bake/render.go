package bake

import (
	"image"
	"image/color"
	"math"

	"github.com/gabrielluizsf/sge/gltf"
)

// Sprite is a baked picture, and the pixel of it the model's origin, the
// point it stands on, falls on.
type Sprite struct {
	Image *image.NRGBA
	Pivot image.Point
}

// view is the way toward the viewer in the game's world: the projection
// shows the world from its +X+Y corner, above.
var view = normalize([3]float64{1, 1, 1})

// Render bakes the model as s shoots it, lit and projected as o says.
func Render(m *gltf.Model, s Shot, o Options) Sprite {
	tris := triangles(m, s, o.Unit)
	sun := normalize(o.Sun)
	f := frame{o: o, n: max(1, o.Samples), sun: sun}
	f.bounds(tris)
	f.shadows = newShadowMap(tris, sun)
	f.raster(tris)
	return f.resolve()
}

// frame is a render in progress: its options, the samples of its pixels,
// and where on the screen it lies.
type frame struct {
	o       Options
	n       int // samples across a pixel
	sun     [3]float64
	minX    float64
	minY    float64
	w, h    int // in pixels
	depth   []float64
	color   [][4]float64 // linear light and coverage, a sample each
	shadows *shadowMap
}

// screen is where the game's projection puts the world point p, in
// pixels from the world's origin.
func (f *frame) screen(p [3]float64) (float64, float64) {
	pr := f.o.Projection
	hx, hy := pr.Width/2/pr.Cell, pr.Height/2/pr.Cell
	return (p[0] - p[1]) * hx, (p[0]+p[1])*hy - p[2]*2*hy
}

// bounds sizes the frame to hold the triangles and their shadows on the
// ground.
func (f *frame) bounds(tris []triangle) {
	minX, minY, maxX, maxY := math.Inf(1), math.Inf(1), math.Inf(-1), math.Inf(-1)
	add := func(p [3]float64) {
		x, y := f.screen(p)
		minX, minY, maxX, maxY = min(minX, x), min(minY, y), max(maxX, x), max(maxY, y)
	}
	add([3]float64{})
	for _, t := range tris {
		for _, v := range t.v {
			add(v.p)
			if f.o.Shadow > 0 && f.sun[2] > 0.05 && v.p[2] > 0 {
				add(sub(v.p, scale(f.sun, v.p[2]/f.sun[2])))
			}
		}
	}
	f.minX, f.minY = math.Floor(minX)-2, math.Floor(minY)-2
	f.w, f.h = int(math.Ceil(maxX)-f.minX)+3, int(math.Ceil(maxY)-f.minY)+3
	samples := f.w * f.n * f.h * f.n
	f.depth = make([]float64, samples)
	f.color = make([][4]float64, samples)
	for i := range f.depth {
		f.depth[i] = math.Inf(-1)
	}
}

func scale(v [3]float64, k float64) [3]float64 { return [3]float64{v[0] * k, v[1] * k, v[2] * k} }

// raster draws each triangle's nearest samples, shaded, and the ground's
// shadow where no triangle is.
func (f *frame) raster(tris []triangle) {
	sw, sh := f.w*f.n, f.h*f.n
	for _, t := range tris {
		var p [3][3]float64
		for k, v := range t.v {
			x, y := f.screen(v.p)
			p[k] = [3]float64{(x - f.minX) * float64(f.n), (y - f.minY) * float64(f.n), dot(v.p, [3]float64{1, 1, 1})}
		}
		facing := dot(cross(sub(t.v[1].p, t.v[0].p), sub(t.v[2].p, t.v[0].p)), view)
		avg := add3(t.v[0].n, t.v[1].n, t.v[2].n)
		if !t.mat.DoubleSided && dot(avg, view) < 0 && math.Abs(facing) > 0 {
			continue // its back is to the viewer
		}
		fill(p, [2]int{sw, sh}, func(x, y int, w [3]float64) {
			i := y*sw + x
			d := w[0]*p[0][2] + w[1]*p[1][2] + w[2]*p[2][2]
			if d <= f.depth[i] {
				return
			}
			c, ok := f.shade(t, w)
			if !ok {
				return
			}
			if z := w[0]*t.v[0].p[2] + w[1]*t.v[1].p[2] + w[2]*t.v[2].p[2]; f.o.Water != nil && z < 0 {
				c = f.o.Water.under(c, -z/max(f.o.Unit, 1e-9))
			}
			f.depth[i], f.color[i] = d, c
		})
	}
	if f.o.Shadow > 0 && f.o.Water == nil {
		f.groundShadow()
	}
}

func add3(a, b, c [3]float64) [3]float64 {
	return [3]float64{a[0] + b[0] + c[0], a[1] + b[1] + c[1], a[2] + b[2] + c[2]}
}

// groundShadow darkens the samples showing bare ground, the plane under
// the model, where the model keeps the sun from it.
func (f *frame) groundShadow() {
	pr := f.o.Projection
	hx, hy := pr.Width/2/pr.Cell, pr.Height/2/pr.Cell
	sw := f.w * f.n
	for i := range f.color {
		if f.color[i][3] > 0 {
			continue
		}
		sx := (float64(i%sw)+0.5)/float64(f.n) + f.minX
		sy := (float64(i/sw)+0.5)/float64(f.n) + f.minY
		a, b := sx/hx, sy/hy // X-Y and X+Y on the ground
		g := [3]float64{(a + b) / 2, (b - a) / 2, 0}
		if dark := 1 - f.shadows.lit(g); dark > 0 {
			f.color[i] = [4]float64{0, 0, 0, -dark * f.o.Shadow} // a shadow: negative coverage
		}
	}
}

// shade is the colour, in linear light, of triangle t at weights w, lit
// by the sun where it reaches, the sky and the ground; ok is false where
// the surface's alpha is cut away.
func (f *frame) shade(t triangle, w [3]float64) ([4]float64, bool) {
	var p, n [3]float64
	var tg [4]float64
	var uv [2]float64
	for k, v := range t.v {
		for i := range 3 {
			p[i] += v.p[i] * w[k]
			n[i] += v.n[i] * w[k]
			tg[i] += v.t[i] * w[k]
		}
		tg[3] += v.t[3] * w[k]
		uv[0] += v.uv[0] * w[k]
		uv[1] += v.uv[1] * w[k]
	}
	n = normalize(n)
	back := dot(n, view) < 0
	if back {
		n = scale(n, -1) // the back of a double-sided surface
	}
	m := t.mat
	if m.normal != nil {
		n = bend(n, tg, texel{m.normal.at(uv), m.NormalScale})
	}
	rough, metal, occlusion := m.Roughness, m.Metallic, 1.0
	if m.mr != nil {
		s := m.mr.at(uv)
		rough, metal = rough*s[1], metal*s[2]
	}
	if m.ao != nil {
		occlusion = 1 + m.OcclusionStrength*(m.ao.at(uv)[0]-1)
	}
	albedo := m.BaseColor // glTF's factor is in linear light already
	if m.tex != nil {
		tc := m.tex.at(uv)
		for i := range albedo {
			albedo[i] *= tc[i]
		}
	}
	if m.AlphaMode != "OPAQUE" && albedo[3] < m.AlphaCutoff {
		return [4]float64{}, false
	}
	lit := f.shadows.lit(p) * max(0, dot(n, f.sun))
	up := 0.5 + 0.5*n[2]
	rough = min(max(rough, 0.05), 1)
	shine := min(max(2/math.Pow(rough, 4)-2, 2), 512)
	hv := normalize([3]float64{f.sun[0] + view[0], f.sun[1] + view[1], f.sun[2] + view[2]})
	spec := math.Pow(max(0, dot(n, hv)), shine) * (1 - rough) * lit
	var c [4]float64
	for i := range 3 {
		ambient := (f.o.Ground[i]*(1-up) + f.o.Sky[i]*up) * occlusion
		diffuse := albedo[i] * (1 - metal) * (ambient + f.o.Light[i]*lit)
		reflect := 0.04*(1-metal) + albedo[i]*metal
		c[i] = diffuse + reflect*spec*f.o.Light[i] + albedo[i]*metal*ambient
	}
	c[3] = 1
	return c, true
}

// resolve averages each pixel's samples into the sprite: colours by how
// much each covers, tone mapped and encoded as sRGB; the ground's shadow
// as black, as dark as it is.
func (f *frame) resolve() Sprite {
	img := image.NewNRGBA(image.Rect(0, 0, f.w, f.h))
	sw := f.w * f.n
	k := float64(f.n * f.n)
	for y := range f.h {
		for x := range f.w {
			var c [3]float64
			var cover, shade float64
			for sy := range f.n {
				for sx := range f.n {
					s := f.color[(y*f.n+sy)*sw+x*f.n+sx]
					if s[3] < 0 {
						shade -= s[3]
						continue
					}
					for i := range 3 {
						c[i] += s[i] * s[3]
					}
					cover += s[3]
				}
			}
			if cover == 0 {
				img.SetNRGBA(x, y, color.NRGBA{A: uint8(math.Round(shade / k * 255))})
				continue
			}
			alpha := cover/k + shade/k*(1-cover/k)
			px := color.NRGBA{A: uint8(math.Round(min(alpha, 1) * 255))}
			// The shadow mixed in at the model's edge darkens it a little.
			darken := cover / (cover + shade)
			e := f.o.Exposure
			if e == 0 {
				e = 1
			}
			px.R = uint8(math.Round(srgb(tone(c[0]/cover*e)*darken) * 255))
			px.G = uint8(math.Round(srgb(tone(c[1]/cover*e)*darken) * 255))
			px.B = uint8(math.Round(srgb(tone(c[2]/cover*e)*darken) * 255))
			img.SetNRGBA(x, y, px)
		}
	}
	return Sprite{Image: img, Pivot: image.Point{X: int(math.Round(-f.minX)), Y: int(math.Round(-f.minY))}}
}

// tone maps linear light of any brightness into 0 to 1: true below a
// knee, so textures keep their colours, and eased above it toward white,
// so bright sun on light stone does not clip.
func tone(c float64) float64 {
	const knee = 0.6
	if c <= knee {
		return max(c, 0)
	}
	return knee + (1-knee)*(1-math.Exp(-(c-knee)/(1-knee)))
}

// texel is a normal map's texel, and how strongly it bends a normal.
type texel struct {
	s        [4]float64
	strength float64
}

// bend is the normal n bent by a normal map's texel, which leans it along
// the tangent tg and its bitangent.
func bend(n [3]float64, tg [4]float64, tx texel) [3]float64 {
	s, strength := tx.s, tx.strength
	t := [3]float64{tg[0], tg[1], tg[2]}
	t = normalize(sub(t, scale(n, dot(n, t))))
	if t == ([3]float64{}) {
		return n
	}
	b := scale(cross(n, t), map[bool]float64{true: -1, false: 1}[tg[3] < 0])
	x, y, z := (s[0]*2-1)*strength, (s[1]*2-1)*strength, s[2]*2-1
	return normalize(add3(scale(t, x), scale(b, y), scale(n, z)))
}
