package installer

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"testing"
)

// testPNG is a w by h PNG, opaque red.
func testPNG(t *testing.T, w, h int) []byte {
	img := image.NewNRGBA(image.Rect(0, 0, w, h))
	for y := range h {
		for x := range w {
			img.SetNRGBA(x, y, color.NRGBA{R: 200, A: 255})
		}
	}
	var b bytes.Buffer
	if err := png.Encode(&b, img); err != nil {
		t.Fatal(err)
	}
	return b.Bytes()
}
