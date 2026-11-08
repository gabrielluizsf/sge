package particle

import (
	"math"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/vector"
)

// fadesFrom is how much of its rest a speck lies at full colour before it
// starts to fade.
const fadesFrom = 0.7

// Paint draws the specks in the air, as drops. Those on the ground are
// drawn by Lying, which goes on an object under what stands there, so that
// nothing lying on the ground is drawn over the people on it.
func (s *System) Paint(o *game.Object, cv *canvas.Canvas) {
	g := o.Game()
	for i := range s.specks {
		if p := &s.specks[i]; p.live && p.z > 0 {
			at := g.ToScreen(vector.Vector3{X: p.x, Y: p.y})
			disc(cv, [2]int{int(at.X), int(at.Y - p.z)}, [2]int{p.size / 2, p.size / 2}, p.color)
		}
	}
}

// Lying is the script that draws a system's specks on the ground, each a
// round stain flat on it, as wide as twice its height, which holds its
// colour and then fades as its rest runs out.
type Lying struct{ System *System }

// Update does nothing: the system moves its specks.
func (Lying) Update(*game.Object, float64) {}

// Paint draws the specks on the ground.
func (l Lying) Paint(o *game.Object, cv *canvas.Canvas) {
	g := o.Game()
	s := l.System
	for i := range s.specks {
		p := &s.specks[i]
		if !p.live || p.z > 0 {
			continue
		}
		at := g.ToScreen(vector.Vector3{X: p.x, Y: p.y})
		disc(cv, [2]int{int(at.X), int(at.Y)}, [2]int{p.size, p.size / 2}, canvas.Fade(p.color, p.opacity()))
	}
}

// opacity is how opaque the speck lying on the ground is, out of 255.
func (p *speck) opacity() int {
	left := (1 - p.lain/p.rest) / (1 - fadesFrom)
	return int(255*min(1, max(0, left)) + 0.5)
}

// disc fills the ellipse of radii r about the middle at, a pixel or more
// across, row by row.
func disc(cv *canvas.Canvas, at, r [2]int, c canvas.Color) {
	rx, ry := max(r[0], 1), max(r[1], 1)
	for dy := -ry + 1; dy < ry; dy++ {
		n := float64(dy) / float64(ry)
		half := int(float64(rx)*math.Sqrt(1-n*n) + 0.5)
		cv.FillRect(canvas.Area{X: at[0] - half, Y: at[1] + dy, Width: 2 * half, Height: 1}, c)
	}
}
