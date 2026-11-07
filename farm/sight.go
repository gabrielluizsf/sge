package farm

import (
	"github.com/gabrielluizsf/sge/pathfind"
	"github.com/gabrielluizsf/sge/vector"
)

// sightMargin is how many cells beyond what the view shows still count
// as in sight: those who walk there are ready to be seen when they come
// in, and need no catching up.
const sightMargin = 6

// inSight reports whether world point p is in the view or just beyond it.
// While the view is not known, nothing is out of sight.
func (f *farm) inSight(p vector.Vector3) bool {
	v := f.player.Game().VisibleWorld()
	if v.Width <= 0 || v.Height <= 0 {
		return true
	}
	m := float64(sightMargin * tile)
	return p.X >= v.X-m && p.X <= v.X+v.Width+m && p.Y >= v.Y-m && p.Y <= v.Y+v.Height+m
}

// plotInSight reports whether plot p's cell is in sight.
func (f *farm) plotInSight(p *plot) bool {
	return f.inSight(f.center(pathfind.Cell{Col: p.col, Row: p.row}))
}
