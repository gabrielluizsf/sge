package pose

import (
	"math"
	"testing"

	"github.com/gabrielluizsf/sge/gltf"
)

// A channel of evenly spaced keys is sampled by division as Sample does
// by search, before, between, on and after its keys, linear or stepped.
func TestEvenKeysSampleAsSearchedOnes(t *testing.T) {
	for _, interp := range []string{"LINEAR", "STEP"} {
		ch := gltf.Channel{Path: "rotation", Interpolation: interp}
		for k := range 24 {
			a := float64(k) * 0.07
			ch.Times = append(ch.Times, 0.5+float32(k)/24)
			ch.Values = append(ch.Values, [4]float32{float32(math.Sin(a)), 0, 0, float32(math.Cos(a))})
		}
		tr := trackOf(&ch)
		if !tr.even {
			t.Fatal("even keys not seen as even")
		}
		for i := -10; i < 200; i++ {
			at := float64(i) / 128
			got, want := tr.sample(&ch, at), Sample(&ch, at)
			for k := range got {
				if math.Abs(got[k]-want[k]) > 1e-6 {
					t.Fatalf("%s at %v: %v, want %v", interp, at, got, want)
				}
			}
		}
	}
	uneven := gltf.Channel{Path: "translation", Times: []float32{0, 0.1, 0.5}, Values: make([][4]float32, 3)}
	if trackOf(&uneven).even {
		t.Error("uneven keys seen as even")
	}
}

// Rotations close enough to be mixed straight stray from the arc by less
// than a hundredth of a degree.
func TestCloseRotationsStayOnTheArc(t *testing.T) {
	th := math.Acos(0.99) * 2 // the widest turn mixed straight
	a := [4]float64{0, 0, 0, 1}
	b := [4]float64{0, 0, math.Sin(th / 2), math.Cos(th / 2)}
	for i := range 11 {
		f := float64(i) / 10
		q := Slerp(a, b, f)
		got := 2 * math.Atan2(q[2], q[3])
		if math.Abs(got-f*th)*180/math.Pi > 0.01 {
			t.Errorf("at %v: %v degrees, want %v", f, got*180/math.Pi, f*th*180/math.Pi)
		}
	}
}
