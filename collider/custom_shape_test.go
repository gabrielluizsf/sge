package collider

import "testing"

// dots is a shape of its own: points that overlap what contains them.
type dots struct{ pts [][2]float64 }

func (d *dots) Bounds() Box { return Box{0, 0, 100, 100} }

func (d *dots) Overlaps(other Collider) bool {
	for _, p := range d.pts {
		if BoxBox(Box{p[0], p[1], 0, 0}, other.Bounds()) {
			return true
		}
	}
	return false
}

func (d *dots) Moved(dx, dy float64) Collider {
	out := &dots{}
	for _, p := range d.pts {
		out.pts = append(out.pts, [2]float64{p[0] + dx, p[1] + dy})
	}
	return out
}

// A shape that knows what it overlaps decides it, on either side, rather
// than being tested by its bounds.
func TestOverlappersDecideTheirCollisions(t *testing.T) {
	d := &dots{pts: [][2]float64{{10, 10}}}
	near, far := Box{5, 5, 10, 10}, Box{50, 50, 10, 10}
	if !Collides(d, near) || !Collides(near, d) {
		t.Error("the dot inside the box does not collide with it")
	}
	if Collides(d, far) || Collides(far, d) {
		t.Error("a box inside the bounds but away from every dot collides")
	}
}

// Offset moves a shape that knows how to move, rather than turning it
// into its bounds.
func TestOffsetMovesMovers(t *testing.T) {
	moved, ok := Offset(&dots{pts: [][2]float64{{10, 10}}}, 40, 40).(*dots)
	if !ok {
		t.Fatalf("Offset turned the shape into %T", moved)
	}
	if !Collides(moved, Box{45, 45, 10, 10}) || Collides(moved, Box{5, 5, 10, 10}) {
		t.Errorf("moved dots = %v", moved.pts)
	}
}
