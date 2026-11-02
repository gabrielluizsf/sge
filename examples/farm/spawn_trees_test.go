package farm

import "testing"

// reachArt is bigArt with its trees grown to their real size, as the
// realistic look has them.
type reachArt struct{ bigArt }

func (reachArt) Reach() int { return 480 }

// treesNear counts the trees within r cells of col, row, on the farm or in
// the wild.
func treesNear(f *farm, col, row, r int) (n int) {
	for c := col - r; c <= col+r; c++ {
		for k := row - r; k <= row+r; k++ {
			id := f.wild.TileAt(wildThings, c, k)
			if c >= 0 && k >= 0 && c < f.things.Cols && k < f.things.Rows {
				id = f.things.At(c, k)
			}
			if isTree(id) {
				n++
			}
		}
	}
	return n
}

func TestSpawnIsNotCrowdedWithTrees(t *testing.T) {
	g, f := build(reachArt{bigArt{SVGArt()}}, idle{})
	where := siteFor(reachArt{bigArt{SVGArt()}})
	for range 3 {
		g.Step(frame)
	}
	const r = 14
	near := treesNear(f, where.spawn.Col, where.spawn.Row, r)
	area := (2*r + 1) * (2*r + 1)
	elsewhere, samples := 0, 0
	for _, p := range [][2]int{{-40, -40}, {-60, 30}, {90, 90}, {40, -70}, {-30, 80}, {120, 0}} {
		f.player.Transform.Position = f.ground.CellCenter(p[0], p[1])
		for range 3 {
			g.Step(frame)
		}
		elsewhere += treesNear(f, p[0], p[1], r)
		samples++
	}
	t.Logf("trees within %d cells of spawn %v: %d of %d; elsewhere avg %d", r, where.spawn, near, area, elsewhere/samples)
	if near > elsewhere/samples {
		t.Errorf("spawn has %d trees, more than the wild's %d on average", near, elsewhere/samples)
	}
}
