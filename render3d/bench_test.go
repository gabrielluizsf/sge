package render3d

import (
	"fmt"
	"testing"

	"github.com/gabrielluizsf/sge/iso"
	"github.com/gabrielluizsf/sge/render3d/internal/shapes"
)

// A frame's work on the CPU for 40 walking characters, and 40 more off
// the screen: culling, choosing detail, posing, palettes and sorting; the
// characters of 163 bones, as MakeHuman's default skeleton has, or of 53,
// as its skeleton for games has.
func BenchmarkPrepare40Skinned(b *testing.B) {
	for _, bones := range []int{163, 53} {
		b.Run(fmt.Sprint(bones, "bones"), func(b *testing.B) { prepare40(b, bones) })
	}
}

// prepare40 is BenchmarkPrepare40Skinned for characters of bones bones.
func prepare40(b *testing.B, bones int) {
	r := NewRenderer(&null{w: 1280, h: 720})
	m, err := r.Upload(shapes.Rig(bones))
	if err != nil {
		b.Fatal(err)
	}
	lod := &LOD{Levels: []Level{{Model: m, Pixels: 40}, {Model: m}}, Hysteresis: 0.1}
	s := &Scene{Camera: Iso(IsoDesc{Projection: iso.New(32, 64, 32), Width: 1280, Height: 720, Depth: 5000})}
	for i := range 80 {
		x := float64(i%8)*40 - 160
		if i >= 40 {
			x += 100000
		}
		s.Instances = append(s.Instances, &Instance{Model: lod, Place: At([3]float64{x, float64(i/8%5) * 40, 0}, float64(i*37), 32), Animation: "walk", Time: float64(i) * 0.1})
	}
	r.Prepare(s)
	if r.Stats.Shown != 40 || r.Stats.Culled != 40 {
		b.Fatalf("stats %+v", r.Stats)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		for _, in := range s.Instances {
			in.Time += 1.0 / 60
		}
		r.Prepare(s)
	}
}
