package gl

import (
	"testing"
	"time"

	"github.com/gabrielluizsf/sge/gltf"
	"github.com/gabrielluizsf/sge/iso"
	"github.com/gabrielluizsf/sge/render3d"
	"github.com/gabrielluizsf/sge/render3d/internal/shapes"
)

// character is a walking rig of bones bones, skinned to boxes of about
// triangles triangles in all: a character of a game's detail.
func character(bones, triangles int) *gltf.Model {
	m := shapes.Rig(bones)
	p := &m.Meshes[0].Primitives[0]
	for b := 0; len(p.Indices) < 3*triangles; b++ {
		j := b % bones
		x := float32(b/bones%5)*0.02 - 0.04
		shapes.Box(p, shapes.Cuboid{Center: [3]float32{x, float32(j) * 0.01, 0.03}, Half: [3]float32{0.01, 0.01, 0.01}, Joint: uint16(j)})
	}
	return m
}

// A frame of 40 walking characters of 20,000 triangles and 163 bones
// each, drawn at 1280 by 720, with the sun's shadow and without: its
// time on the CPU and the GPU, waited for to the end.
func BenchmarkDraw40Skinned(b *testing.B) {
	for _, c := range []struct {
		name   string
		shadow int
	}{{"sun", 2048}, {"noshadow", 0}} {
		b.Run(c.name, func(b *testing.B) { draw40(b, c.shadow) })
	}
}

// draw40 is BenchmarkDraw40Skinned with a shadow map shadow pixels
// square, 0 for none.
func draw40(b *testing.B, shadow int) {
	dev := offscreen(b, 1280, 720)
	r := render3d.NewRenderer(dev)
	m, err := r.Upload(character(163, 20000))
	if err != nil {
		b.Fatal(err)
	}
	s := &render3d.Scene{
		Camera: render3d.Iso(render3d.IsoDesc{Projection: iso.New(32, 64, 32), Width: 1280, Height: 720, Depth: 5000}),
		Sun: render3d.Sun{Toward: [3]float64{1, 0.5, 2}, Color: [3]float64{1, 1, 1}, Ambient: [3]float64{0.3, 0.3, 0.3},
			Shadow: shadow, Area: render3d.Sphere{Radius: 250}, Bias: 0.002},
		Clear: sky,
	}
	for i := range 40 {
		s.Instances = append(s.Instances, &render3d.Instance{Model: render3d.One(m),
			Place: render3d.At([3]float64{float64(i%8)*40 - 160, float64(i/8%5) * 40, 0}, float64(i*37), 32), Animation: "walk", Time: float64(i) * 0.1})
	}
	if err := r.Draw(s); err != nil {
		b.Fatal(err)
	}
	if r.Stats.Shown != 40 || r.Stats.Triangles < 40*20000 {
		b.Fatalf("stats %+v", r.Stats)
	}
	dev.Present()
	b.ResetTimer()
	start := time.Now()
	for range b.N {
		for _, in := range s.Instances {
			in.Time += 1.0 / 60
		}
		if err := r.Draw(s); err != nil {
			b.Fatal(err)
		}
		dev.Present()
	}
	b.ReportMetric(float64(time.Since(start).Microseconds())/1000/float64(b.N), "ms/frame")
}
