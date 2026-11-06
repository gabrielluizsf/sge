package game

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/vector"
)

// RevealOpacity is how much of a tall thing shows, of 1, while it stands
// over a Reveal.
const RevealOpacity = 0.35

// Reveal is something a projected world keeps in sight: a standing thing
// at a world point, Height pixels tall on the screen, a player, say, or
// the cell the cursor is on. A tall tile or an object that Fades, nearer
// than it and covering it, is drawn faded, as isometric games do, so a
// tree never hides the farmer walking behind it.
type Reveal struct {
	At     vector.Vector3
	Height int
}

// Cover is a Reveal on the screen: its feet, and how far up it reaches.
type Cover struct{ X, Y, Height int }

// Covers are where the Reveals of the frame being drawn stand on the
// screen.
func (g *Game) Covers() []Cover { return g.covers }

// findCovers finds where each Reveal stands on the screen, once for a
// frame.
func (g *Game) findCovers() {
	g.covers = g.covers[:0]
	if g.Reveals == nil || g.Projection == nil {
		return
	}
	for _, r := range g.Reveals() {
		s := g.ToScreen(r.At)
		g.covers = append(g.covers, Cover{X: int(s.X + 0.5), Y: int(s.Y + 0.5), Height: r.Height})
	}
}

// Placed is a picture on the screen: its top-left pixel is at (X, Y).
type Placed struct {
	Img  *canvas.Canvas
	X, Y int
}

// Hides reports whether the picture p, standing at screen depth
// nearness, hides one of covers behind it: one that stands further back,
// a shown pixel of the picture lying over its body.
func Hides(p Placed, nearness float64, covers []Cover) bool {
	img, x, y := p.Img, p.X, p.Y
	for _, c := range covers {
		if float64(c.Y) >= nearness {
			continue // it stands in front of the picture
		}
		px := c.X - x
		if px < 0 || px >= img.Width {
			continue
		}
		for py := max(c.Y-c.Height-y, 0); py <= min(c.Y-y, img.Height-1); py++ {
			if img.At(px, py).A() > 0 {
				return true
			}
		}
	}
	return false
}

// BlendFaded blends the picture p over cv at RevealOpacity: the picture
// faded once, and kept for the next frames, is blended as any is.
func BlendFaded(cv *canvas.Canvas, p Placed) {
	cv.Blit(p.X, p.Y, fadedOf(p.Img))
}
