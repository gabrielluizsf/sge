package label

import (
	"os"
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
)

func face(t *testing.T) *canvas.Face {
	data, err := os.ReadFile("../examples/farm/fonts/Sniglet-ExtraBold.ttf")
	if err != nil {
		t.Fatal(err)
	}
	f, err := canvas.ParseFace(data, 18)
	if err != nil {
		t.Fatal(err)
	}
	return f
}

var (
	cream = canvas.RGB(255, 246, 228)
	cocoa = canvas.RGB(58, 42, 46)
)

// A glyph's edge is part transparent in the ink, never an opaque dark
// fringe; inside it is the ink, and around the text, beyond the outline,
// nothing.
func TestALabelsEdgesAreSoftNotDark(t *testing.T) {
	l := New(face(t), "o", Style{Ink: cream})
	solid, soft := 0, 0
	for y := range l.Height {
		for x := range l.Width {
			c := l.At(x, y)
			switch a := c.A(); {
			case a == 255:
				solid++
				if c != cream {
					t.Fatalf("an opaque pixel %v is not the ink", c)
				}
			case a > 0:
				soft++
				if c.R() != cream.R() || c.B() != cream.B() {
					t.Fatalf("a soft pixel %v is not the ink, less opaque", c)
				}
			}
		}
	}
	if solid == 0 || soft == 0 {
		t.Errorf("%d solid and %d soft pixels", solid, soft)
	}
}

// An outline surrounds the ink, and widens the picture by its width.
func TestALabelIsOutlined(t *testing.T) {
	f := face(t)
	plain := New(f, "Oi!", Style{Ink: cream})
	l := New(f, "Oi!", Style{Ink: cream, Outline: cocoa, Width: 2})
	if l.Width != plain.Width+4 || l.Height != plain.Height+4 {
		t.Fatalf("%dx%d outlined, %dx%d plain", l.Width, l.Height, plain.Width, plain.Height)
	}
	outline := 0
	for y := range l.Height {
		for x := range l.Width {
			if l.At(x, y) == cocoa {
				outline++
			}
		}
	}
	if outline == 0 || l.At(0, 0).A() != 0 {
		t.Errorf("%d outline pixels, corner %v", outline, l.At(0, 0))
	}
	if e := New(f, "", Style{Ink: cream}); e.Width < 1 {
		t.Error("an empty label has no picture")
	}
}
