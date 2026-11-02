package farm

import (
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/tilemap"
	"github.com/gabrielluizsf/sge/town"
)

// bigArt is the stylized art with buildings as large as a realistic look
// has them: a three by two stall, a two by two well, houses eight by six.
type bigArt struct{ Art }

func (bigArt) Footprint(name string) (cols, rows int) {
	switch name {
	case "stall":
		return 3, 2
	case "well":
		return 2, 2
	}
	return 8, 6
}

// pacer is a game, and the farm in it, to pacer the player about in.
type pacer struct {
	g *game.Game
	f *farm
}

// at puts the player on cell col, row and lets the game run a few frames
// so the world near them is made, as it is where the camera is.
func (w pacer) at(col, row int) {
	w.f.player.Transform.Position = w.f.ground.CellCenter(col, row)
	for range 3 {
		w.g.Step(frame)
	}
}

// bumps walks the player into cell col, row from each side that is open
// and reports whether they never got into it, and whether there was a side
// to pacer from.
func (w pacer) bumps(col, row int) (blocked, tried bool) {
	for _, d := range [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
		sc, sr := col+d[0], row+d[1]
		w.at(sc, sr)
		if !w.f.open(sc, sr) {
			continue
		}
		tried = true
		w.at(sc, sr)
		o := w.f.player
		tilemap.Move(o, -float64(d[0])*1.4*tile, -float64(d[1])*1.4*tile, w.f.solids()...)
		if c, r := w.f.cellOf(feet(o)); c == col && r == row {
			return false, true
		}
	}
	return true, tried
}

// mustBump fails unless walking into cell col, row is stopped.
func (w pacer) mustBump(t *testing.T, at town.Cell, what string) {
	col, row := at.Col, at.Row
	t.Helper()
	if blocked, tried := w.bumps(col, row); !tried || !blocked {
		t.Errorf("%s at %d,%d: the player walks in (open side: %v)", what, col, row, tried)
	}
}

// mustBumpAll fails unless every cell of r is a wall.
func (w pacer) mustBumpAll(t *testing.T, r rect, what string) {
	t.Helper()
	for row := r.row; row < r.row+r.rows; row++ {
		for col := r.col; col < r.col+r.cols; col++ {
			if w.f.open(col, row) {
				t.Errorf("%s: cell %d,%d is open", what, col, row)
			}
		}
	}
	for _, c := range []town.Cell{{Col: r.col, Row: r.row}, {Col: r.col + r.cols - 1, Row: r.row + r.rows - 1}, {Col: r.col + r.cols/2, Row: r.row}} {
		w.mustBump(t, c, what)
	}
}

// Everything in the game is a wall to pacer into, whatever the size of the
// buildings: the farm's house and stall, a town's houses, stall, anvil
// and forge, the fences, trees and stumps; and nothing stands over
// another.
func TestEverythingBlocksTheWay(t *testing.T) {
	for name, a := range map[string]Art{"stylized": SVGArt(), "large": bigArt{SVGArt()}} {
		g, f := build(a, idle{})
		cv, _ := canvas.NewCanvas(320, 240)
		g.Draw(cv)
		w := pacer{g, f}
		s := siteFor(a)
		w.mustBumpAll(t, s.home, name+" home")
		w.mustBumpAll(t, s.stall, name+" farm stall")
		w.mustBump(t, town.Cell{Col: 5}, name+" fence")
		trees := 0
		for row := range f.things.Rows {
			for col := range f.things.Cols {
				if isTree(f.things.At(col, row)) {
					w.mustBump(t, town.Cell{Col: col, Row: row}, name+" tree")
					trees++
				}
			}
		}
		if trees == 0 {
			t.Error("no trees on the farm")
		}
		tn := f.towns.Of(1, 0)
		for _, h := range tn.Houses {
			w.mustBumpAll(t, houseRect(h), name+" house")
		}
		w.mustBumpAll(t, stallRect(a, tn.Stall), name+" town stall")
		if r := wellRect(a, tn); r.cols > 0 {
			w.mustBumpAll(t, r, name+" well")
		}
		w.mustBump(t, tn.Anvil, name+" anvil")
		w.mustBump(t, tn.Forge, name+" forge")
		// A felled tree leaves a stump, which is solid too.
		f.setThing(tn.Center.Col+20, tn.Center.Row, tileID["stump"])
		w.mustBump(t, town.Cell{Col: tn.Center.Col + 20, Row: tn.Center.Row}, name+" stump")
		wild := 0
		for col := tn.Center.Col + 12; col < tn.Center.Col+40 && wild < 5; col++ {
			if isTree(f.thingAt(col, tn.Center.Row-12)) {
				w.mustBump(t, town.Cell{Col: col, Row: tn.Center.Row - 12}, name+" wild tree")
				wild++
			}
		}
	}
}

// The farm's stall and house stand apart, and so do a town's houses and
// stall; the player starts on open ground in front of their house.
func TestBuildingsStandApartAndThePlayerStartsInTheOpen(t *testing.T) {
	for name, a := range map[string]Art{"stylized": SVGArt(), "large": bigArt{SVGArt()}} {
		_, f := build(a, idle{})
		s := siteFor(a)
		for row := s.home.row - 1; row <= s.home.row+s.home.rows; row++ {
			for col := s.home.col - 1; col <= s.home.col+s.home.cols; col++ {
				if s.stall.has(col, row) {
					t.Errorf("%s: the stall touches the house at %d,%d", name, col, row)
				}
			}
		}
		if c, r := f.cellOf(feet(f.player)); !f.open(c, r) || s.home.has(c, r) || s.stall.has(c, r) {
			t.Errorf("%s: the player starts at %d,%d, not on open ground", name, c, r)
		}
		if name == "large" && (feet(f.player).Y/tile < float64(s.home.row+s.home.rows) || s.spawn != s.door) {
			t.Errorf("the player does not start in front of their house: %v, door %v", s.spawn, s.door)
		}
		tn := f.towns.Of(1, 0)
		st := stallRect(a, tn.Stall)
		for _, h := range tn.Houses {
			for row := h.Row - 1; row <= h.Row+h.Depth; row++ {
				for col := h.Col - 1; col <= h.Col+h.Width; col++ {
					if st.has(col, row) {
						t.Errorf("%s: the stall touches house %+v", name, h)
					}
				}
			}
		}
	}
}
