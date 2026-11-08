package human

import (
	"image"
	"image/color"
	"math"
	"testing"
)

// noisy is a size by size image of a base colour with a deterministic
// pattern of light and dark across it, and alpha 255.
func noisy(size int, base color.NRGBA) *image.NRGBA {
	img := image.NewNRGBA(image.Rect(0, 0, size, size))
	for y := range size {
		for x := range size {
			f := 0.55 + 0.9*float64((x*7+y*13)%11)/10
			img.SetNRGBA(x, y, color.NRGBA{uint8(min(255, float64(base.R)*f)), uint8(min(255, float64(base.G)*f)), uint8(min(255, float64(base.B)*f)), 255})
		}
	}
	return img
}

// stats are the mean colour, sRGB 0..1, and the coefficient of variation
// of the light of img's pixels.
func stats(img *image.NRGBA) (mean [3]float64, cv float64) {
	n := float64(len(img.Pix) / 4)
	var ys []float64
	for i := 0; i < len(img.Pix); i += 4 {
		var y float64
		for c := range 3 {
			v := float64(img.Pix[i+c]) / 255
			mean[c] += v / n
			y += v * [3]float64{0.2126, 0.7152, 0.0722}[c]
		}
		ys = append(ys, y)
	}
	var m, v float64
	for _, y := range ys {
		m += y / n
	}
	for _, y := range ys {
		v += (y - m) * (y - m) / n
	}
	return mean, math.Sqrt(v) / m
}

func dist(a, b [3]float64) float64 {
	return math.Sqrt((a[0]-b[0])*(a[0]-b[0]) + (a[1]-b[1])*(a[1]-b[1]) + (a[2]-b[2])*(a[2]-b[2]))
}

func TestTintMovesTheColourAndKeepsTheDetail(t *testing.T) {
	src := noisy(64, color.NRGBA{120, 80, 50, 255})
	srcMean, srcCV := stats(src)
	for _, to := range []Colour{RGB(0.8, 0.6, 0.25), RGB(0.1, 0.1, 0.12), RGB(0.55, 0.2, 0.1), RGB(0.2, 0.4, 0.8)} {
		out := tinted(src, tint{to: to.linear(), keep: 0})
		mean, cv := stats(out)
		want := [3]float64{to.R, to.G, to.B}
		if dist(mean, want) > 0.08 || dist(mean, want) >= dist(srcMean, want) {
			t.Errorf("tinted toward %v: mean %v (was %v)", want, mean, srcMean)
		}
		if cv < 0.6*srcCV {
			t.Errorf("tinted toward %v: variation %.3f of %.3f is flat", want, cv, srcCV)
		}
		if out.Pix[3] != 255 || src.Pix[0] == 0 {
			t.Error("alpha lost")
		}
	}
}

func TestTintKeepsAlphaAndLeavesTheSourceAlone(t *testing.T) {
	src := noisy(8, color.NRGBA{100, 90, 80, 255})
	src.SetNRGBA(1, 1, color.NRGBA{100, 90, 80, 0})
	before := append([]uint8(nil), src.Pix...)
	out := tinted(src, tint{to: RGB(0.9, 0.2, 0.2).linear()})
	if out.NRGBAAt(1, 1).A != 0 {
		t.Error("alpha changed")
	}
	for i := range before {
		if before[i] != src.Pix[i] {
			t.Fatal("the source was changed")
		}
	}
}

func TestSkinToneRampGoesFromLightToDark(t *testing.T) {
	prev := 2.0
	for i := 0; i <= 20; i++ {
		c := skinColour(0.001 + float64(i)/20*0.999)
		l := 0.2126*c.R + 0.7152*c.G + 0.0722*c.B
		if l >= prev {
			t.Fatalf("tone step %d is not darker: %v >= %v", i, l, prev)
		}
		prev = l
	}
}

func TestEyeTintLeavesTheWhitesAlone(t *testing.T) {
	img := image.NewNRGBA(image.Rect(0, 0, 100, 100))
	for i := 0; i < len(img.Pix); i += 4 {
		img.Pix[i], img.Pix[i+1], img.Pix[i+2], img.Pix[i+3] = 230, 220, 215, 255
	}
	out := tinted(img, tint{to: RGB(0.2, 0.4, 0.9).linear(), mask: irises})
	if out.NRGBAAt(5, 5) != img.NRGBAAt(5, 5) || out.NRGBAAt(50, 95) != img.NRGBAAt(50, 95) {
		t.Error("the whites of the eyes were tinted")
	}
	if c := out.NRGBAAt(72, 30); c.B <= c.R {
		t.Errorf("the iris was not tinted: %v", c)
	}
}
