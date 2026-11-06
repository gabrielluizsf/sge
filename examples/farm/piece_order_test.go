package farm

import (
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
)

// block is an art whose every tile is a solid slab standing on its cell, a
// piece of furniture a tall box of the colour of its name.
type block struct {
	Art
	box *canvas.Canvas
}

func (b block) Tile(name string) *canvas.Canvas {
	if _, ok := pieceRefs[pieceTileName(name, 0, 0, 0)]; ok {
		return b.box
	}
	return b.Art.Tile(name)
}

func (block) Furniture(name string) (int, int) {
	if name == "bed" {
		return 2, 1
	}
	return 1, 1
}

// solidBox is a w by h picture of one colour all over.
func solidBox(w, h int, c canvas.Color) *canvas.Canvas {
	cv, _ := canvas.NewCanvas(w, h)
	for i := range cv.Pixels {
		cv.Pixels[i] = c
	}
	return cv
}

// The parts of a piece's picture, one for each cell, make the picture again
// with nothing lost and nothing drawn twice, in every heading.
func TestSlicesMakeThePictureAgain(t *testing.T) {
	full := solidBox(isoW+20, isoH, canvas.RGB(200, 100, 50))
	for _, dims := range [][2]int{{2, 1}, {1, 2}, {2, 2}} {
		cols, rows := dims[0], dims[1]
		whole, _ := canvas.NewCanvas(8*isoW, 8*isoH)
		parts, _ := canvas.NewCanvas(8*isoW, 8*isoH)
		ax, ay := 4*isoW, 6*isoH                                                // the back cell's bottom corner
		whole.Blit(ax-full.Width/2, ay+pieceDrop(cols, rows)-full.Height, full) // its bottom on the front cell's corner
		for j := range rows {
			for i := range cols {
				s := sliceOf(full, cols, rows, [2]int{i, j})
				sx, sy := (i-j)*isoW/2, (i+j)*isoW/4
				parts.Blit(ax+sx-s.Width/2, ay+sy-s.Height, s) // placed by its bottom middle on its cell's corner
			}
		}
		for i := range whole.Pixels {
			if (whole.Pixels[i].A() != 0) != (parts.Pixels[i].A() != 0) {
				t.Fatalf("%dx%d: pixel %d,%d differs", cols, rows, i%whole.Width, i/whole.Width)
			}
		}
	}
}

// Wherever a pixel falls, it is on a part whose frame reaches it: a part
// is placed by its bottom on its own cell's corner and can draw nothing
// lower, so what lies under the back cell's corner, along the piece's long
// side, must go with a nearer cell.
func TestNoPixelFallsBelowItsCellsFrame(t *testing.T) {
	for _, dims := range [][2]int{{2, 1}, {1, 2}, {2, 2}, {3, 1}} {
		cols, rows := dims[0], dims[1]
		for y := -isoH; y < pieceDrop(cols, rows); y++ {
			for x := -2 * isoW; x <= 2*isoW; x++ {
				at := cellOf(x, y, cols, rows)
				if y >= (at[0]+at[1])*isoW/4 {
					t.Fatalf("%dx%d: pixel %d,%d goes to cell %v, whose frame ends above it", cols, rows, x, y, at)
				}
			}
		}
	}
}

// Each cell's part lies where its cell is: no part is drawn at a cell it
// does not belong to. A part of a two-cell piece has pixels of its own.
func TestEachCellHasItsOwnPart(t *testing.T) {
	full := solidBox(isoW+20, isoH, canvas.RGB(200, 100, 50))
	for _, at := range [][2]int{{0, 0}, {1, 0}} {
		s := sliceOf(full, 2, 1, at)
		n := 0
		for _, p := range s.Pixels {
			if p.A() != 0 {
				n++
			}
		}
		if n == 0 {
			t.Errorf("cell %v has no part", at)
		}
	}
}

// A person standing on any side of a two-cell piece, turned any way, in
// any place on the floor, is drawn over the parts that are farther from
// the viewer than they are and under those that are nearer: each part is
// a tile of its own cell, at its cell's depth, among the people. This is
// what once let someone walk under the back half of a bed.
func TestPeopleAreDrawnInFrontOrBehindAnyPartOfAPiece(t *testing.T) {
	g, f := build(block{SVGArt(), solidBox(isoW+20, isoH, canvas.RGB(200, 100, 50))}, idle{})
	f.enter(houseDoor{at: homeDoor, faceOut: "down", seed: seed, own: true})
	g.Step(frame)
	depth := func(c, r int) float64 { return g.ToScreen(f.ground.CellCenter(c, r)).Y }
	checked := 0
	for h := range headings {
		for r := 1; r <= roomRows-2; r++ {
			for c := 1; c <= roomCols-2; c++ {
				f.home = []piece{newPiece("bed", c, r, h)}
				if !houseFits(f.art, f.home) {
					continue
				}
				f.relay()
				s := f.home[0].cells(f.art)
				for j := range s.rows {
					for i := range s.cols {
						if got := f.room.things.At(c+i, r+j); got != pieceTileID(f.home[0], i, j) {
							t.Fatalf("heading %d at %d,%d: cell %d,%d shows tile %d, not its part", h, c, r, i, j, got)
						}
					}
				}
				for pr := range roomRows {
					for pc := range roomCols {
						if s.has(pc, pr) {
							continue
						}
						for j := range s.rows {
							for i := range s.cols {
								part, person := depth(c+i, r+j), depth(pc, pr)
								nearer := pc+pr > c+i+r+j
								switch {
								case pc+pr == c+i+r+j:
								case nearer != (person > part):
									t.Errorf("heading %d, bed at %d,%d: person at %d,%d vs part %d,%d drawn the wrong way", h, c, r, pc, pr, c+i, r+j)
								}
								checked++
							}
						}
					}
				}
			}
		}
	}
	if checked < 1000 {
		t.Fatalf("only %d cases checked", checked)
	}
}

// Someone standing behind the right of a bed (the cell that tied with the
// bed's back cell when the whole bed was drawn from it) is drawn under the
// bed's second cell.
func TestThePersonBesideTheBedsFarEndIsBehindIt(t *testing.T) {
	g, f := build(block{SVGArt(), solidBox(isoW+20, isoH, canvas.RGB(200, 100, 50))}, idle{})
	f.enter(houseDoor{at: homeDoor, faceOut: "down", seed: seed, own: true})
	g.Step(frame)
	f.home = []piece{newPiece("bed", 2, 3, 0)}
	f.relay()
	depth := func(c, r int) float64 { return g.ToScreen(f.ground.CellCenter(c, r)).Y }
	if !(depth(3, 2) < depth(3, 3)) { // behind the second cell: drawn first
		t.Fatal("the cell behind the bed's far end is not farther than its second cell")
	}
	if f.room.things.At(3, 3) != pieceTileID(f.home[0], 1, 0) {
		t.Error("the bed's second cell is no tile of its own")
	}
}
