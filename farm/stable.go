package farm

import (
	"github.com/gabrielluizsf/sge/i18n"
	"github.com/gabrielluizsf/sge/pathfind"
	"github.com/gabrielluizsf/sge/tilemap"
	"github.com/gabrielluizsf/sge/vector"
)

// offer is an animal a stable sells.
type offer struct {
	kind, label string
	price       int
}

// stableOffers are what a stable sells: sheep first, the cheaper, then
// horses.
var stableOffers = []offer{
	{kind: "sheep", label: "Sheep", price: 120},
	{kind: "horse", label: "Horse", price: 420},
	{kind: "cow", label: "Cow", price: 300},
}

// penRoom is how many of a kind one pen holds.
var penRoom = map[string]int{"cow": 3, "sheep": 6, "horse": 3}

// buyAnimal buys k the animal o, if the money is there and a pen has or
// can be given room for it on the farm, and says so.
func (f *farm) buyAnimal(k *kit, o offer) bool {
	at := f.player.Transform.Position
	at.Y -= tile
	if k.money < o.price {
		f.float(at, i18n.T("Not enough money"), 'r')
		f.sound("deny", nil)
		return false
	}
	p := f.herd.penFor(o.kind)
	if p == nil {
		f.notify(i18n.T("There is no room for another pen"))
		f.sound("deny", nil)
		return false
	}
	k.money -= o.price
	f.herd.add(o.kind, f.luck.Intn(animalKinds[o.kind]), p)
	f.float(at, i18n.T("-$%d", o.price), 'r')
	f.notify(i18n.T("Your %s is in a pen by the farm: press B there to move it", i18n.T(o.label)))
	f.sound("coin", &at)
	return true
}

// penFor is a pen of the farm with room for kind: one that has it, or a new
// one, or nil when no ground is free for it.
func (h *herd) penFor(kind string) *pen {
	for _, p := range h.pens {
		if p.town == nil && p.kind == kind && h.count(p) < penRoom[kind] {
			return p
		}
	}
	r, ok := h.f.freePenSpot()
	if !ok {
		return nil
	}
	p := &pen{kind: kind, r: r}
	h.pens = append(h.pens, p)
	h.f.fence(r, true)
	return p
}

// count is how many animals live in p.
func (h *herd) count(p *pen) int {
	n := 0
	for _, c := range h.crits {
		if c.pen == p {
			n++
		}
	}
	return n
}

// freePenSpot is a place near the player's house for a new pen.
func (f *farm) freePenSpot() (rect, bool) {
	d := homeDoor
	for reach := 4; reach < 40; reach++ {
		for dr := -reach; dr <= reach; dr++ {
			for dc := -reach; dc <= reach; dc++ {
				if max(abs(dc), abs(dr)) != reach {
					continue
				}
				r := rect{d.Col + dc, d.Row + dr, penCols, penRows}
				if f.penFits(r, nil) && f.reachable(r) {
					return r, true
				}
			}
		}
	}
	return rect{}, false
}

// penFits reports whether a pen may stand on r: free dry ground with nothing
// on it, away from the towns, the plots, the doors and the stalls, leaving
// its own cells out of it, when it is the pen moving.
func (f *farm) penFits(r rect, moving *pen) bool {
	pc, pr := f.cellOf(feet(f.player))
	for row := r.row - 1; row <= r.row+r.rows; row++ {
		for col := r.col - 1; col <= r.col+r.cols; col++ {
			inside := r.has(col, row)
			if inside && !f.penCellFree(col, row, moving) {
				return false
			}
			if _, door := f.doorAt(col, row); door {
				return false
			}
			if _, stall := f.stallAt(col, row); stall || f.plotAt(col, row) != nil {
				return false
			}
			if inside && fenceAt(r, col, row) >= 0 && col == pc && row == pr {
				return false
			}
		}
	}
	gate := r.col + r.cols/2
	if !f.open(gate, r.row+r.rows) {
		return false // the gate must not open on a fence or the water
	}
	return f.towns.At(r.col, r.row) == nil || !f.towns.Grounds(r.col+r.cols/2, r.row+r.rows/2, 3)
}

// penCellFree reports whether nothing stands on the cell but moving's fence.
func (f *farm) penCellFree(col, row int, moving *pen) bool {
	if moving != nil && moving.r.has(col, row) {
		return true
	}
	return f.open(col, row) && f.thingAt(col, row) == tilemap.Empty
}

// fence puts up, or takes down, the fence of a pen on r.
func (f *farm) fence(r rect, up bool) {
	for row := r.row; row < r.row+r.rows; row++ {
		for col := r.col; col < r.col+r.cols; col++ {
			if id := fenceAt(r, col, row); id >= 0 {
				if !up {
					id = tilemap.Empty
				}
				f.setThing(col, row, id)
			}
		}
	}
}

// movePen stands p on r, with its fence and its animals, as they stood in it.
func (f *farm) movePen(p *pen, r rect) {
	d := vector.Vector3{X: float64(r.col-p.r.col) * tile, Y: float64(r.row-p.r.row) * tile}
	f.fence(p.r, false)
	p.r = r
	f.fence(r, true)
	for _, c := range f.herd.crits {
		if c.pen == p {
			c.pos, c.to = c.pos.Add(d), c.to.Add(d)
			c.walking = false
		}
	}
}

// reachable reports whether one can walk from the player's door to pen r's gate.
func (f *farm) reachable(r rect) bool {
	front := pathfind.Cell{Col: r.col + r.cols/2, Row: r.row + r.rows}
	door := pathfind.Cell{Col: homeDoor.Col, Row: homeDoor.Row}
	return door == front || len(f.grid().Path(door, front)) > 0
}
