package sprite

import (
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
)

// An SVG drawn at twice its size keeps its pixels square: pixel art made
// as SVG scales up sharp.
func TestFromSVGRendersAtTheSizeAsked(t *testing.T) {
	src := `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 2 2" width="2" height="2" shape-rendering="crispEdges">` +
		`<rect x="0" y="0" width="1" height="1" fill="#ff0000"/><rect x="1" y="1" width="1" height="1" fill="#0000ff"/></svg>`
	s, err := FromSVG(src, 4, 4)
	if err != nil {
		t.Fatal(err)
	}
	img := s.Image
	if img.Width != 4 || img.Height != 4 || s.PivotX != 0.5 || s.PivotY != 0.5 {
		t.Fatalf("%dx%d pivot %v,%v", img.Width, img.Height, s.PivotX, s.PivotY)
	}
	red, blue := canvas.RGB(255, 0, 0), canvas.RGB(0, 0, 255)
	for _, p := range [][2]int{{0, 0}, {1, 1}} {
		if img.At(p[0], p[1]) != red {
			t.Errorf("pixel %v is %v, not red", p, img.At(p[0], p[1]))
		}
	}
	if img.At(3, 3) != blue || img.At(2, 2) != blue || img.At(3, 0).A() != 0 {
		t.Errorf("corners %v %v %v", img.At(3, 3), img.At(2, 2), img.At(3, 0))
	}
}

func TestFromSVGRefusesWhatIsNotSVG(t *testing.T) {
	if _, err := FromSVG("<html>", 4, 4); err == nil {
		t.Error("no error")
	}
	if _, err := FromSVG(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1 1"/>`, 0, 4); err == nil {
		t.Error("no error for a sprite 0 pixels wide")
	}
}
