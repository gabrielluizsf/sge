package tilemap

import (
	"math"

	"github.com/gabrielluizsf/sge/collider"
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/vector"
)

// Move moves o by dx then by dy, as far as it goes before its collider
// meets a solid tile, and reports whether each move was stopped: walking
// into a wall slides along it. An object without a collider moves freely.
func (m *Map) Move(o *game.Object, dx, dy float64) (hitX, hitY bool) {
	return Move(o, dx, dy, m)
}

// Move moves o by dx then by dy through all of solids at once, the layers
// of a level say, stopping at the first wall of any: walking into a wall
// slides along it. It reports whether each move was stopped. An object
// without a collider moves freely.
func Move(o *game.Object, dx, dy float64, solids ...Solid) (hitX, hitY bool) {
	if o.Collider == nil {
		o.Transform.Translate(vector.Vector3{X: dx, Y: dy})
		return false, false
	}
	stride := math.Inf(1)
	for _, s := range solids {
		stride = min(stride, s.Stride())
	}
	blocked := func() bool {
		c := collider.Offset(o.Collider, o.Transform.Position.X, o.Transform.Position.Y)
		for _, s := range solids {
			if s.Blocks(c) {
				return true
			}
		}
		return false
	}
	hitX = slide(o, dx, 0, stride, blocked)
	hitY = slide(o, 0, dy, stride, blocked)
	return hitX, hitY
}

// slide moves o by dx, dy along one axis, in steps of stride so it cannot
// jump over a thin wall, backing off at the first step that is blocked to
// where it last was not, to within a hundredth of a pixel.
func slide(o *game.Object, dx, dy, stride float64, blocked func() bool) bool {
	d := max(math.Abs(dx), math.Abs(dy))
	if d == 0 {
		return false
	}
	steps := 1
	if stride > 0 && !math.IsInf(stride, 1) {
		steps = int(math.Ceil(d / stride))
	}
	p := &o.Transform.Position
	from := *p
	at := func(f float64) {
		p.X, p.Y = from.X+dx*f, from.Y+dy*f
	}
	for i := 1; i <= steps; i++ {
		hi := float64(i) / float64(steps)
		at(hi)
		if !blocked() {
			continue
		}
		lo := float64(i-1) / float64(steps) // clear: the step before was
		for (hi-lo)*d > 0.01 {
			mid := (lo + hi) / 2
			at(mid)
			if blocked() {
				hi = mid
			} else {
				lo = mid
			}
		}
		at(lo)
		return true
	}
	return false
}
