package soft

import (
	"image"
	"image/color"
	"testing"

	"github.com/gabrielluizsf/sge/pose"
	"github.com/gabrielluizsf/sge/render3d"
	"github.com/gabrielluizsf/sge/render3d/internal/shapes"
)

// A material whose texture and colour were changed in the source model,
// then retextured, draws the new ones without the model being uploaded
// again.
func TestRetextureShowsTheNewColours(t *testing.T) {
	dev := New(64, 64)
	r := render3d.NewRenderer(dev)
	src := shapes.Cube()
	m, err := r.Upload(src)
	if err != nil {
		t.Fatal(err)
	}
	cam := render3d.Perspective(render3d.PerspectiveDesc{Eye: [3]float64{0, -5, 0}, Up: [3]float64{0, 0, 1}, FovY: 40, Aspect: 1, Near: 0.1, Far: 20})
	in := &render3d.Instance{Model: render3d.One(m), Place: pose.Identity}
	green := image.NewNRGBA(image.Rect(0, 0, 2, 2))
	for i := range 4 {
		green.SetNRGBA(i%2, i/2, color.NRGBA{0, 255, 0, 255})
	}
	src.Images[0] = green
	if err := r.Retexture(m, 0); err != nil {
		t.Fatal(err)
	}
	if err := r.Draw(front(cam, in)); err != nil {
		t.Fatal(err)
	}
	img, _ := dev.Read(render3d.Screen)
	if c := img.NRGBAAt(27, 32); !shapes.Near(c, [3]int{0, 255, 0}) {
		t.Errorf("left of the face is %v, not green", c)
	}
	src.Materials[0].Texture, src.Materials[0].BaseColor = -1, [4]float64{1, 0, 0, 1}
	if err := r.Retexture(m, 0); err != nil {
		t.Fatal(err)
	}
	if err := r.Draw(front(cam, in)); err != nil {
		t.Fatal(err)
	}
	img, _ = dev.Read(render3d.Screen)
	if c := img.NRGBAAt(37, 32); !shapes.Near(c, [3]int{255, 0, 0}) {
		t.Errorf("right of the face is %v, not red", c)
	}
}
