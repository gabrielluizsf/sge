package farm

import (
	"fmt"

	"github.com/gabrielluizsf/sge/tilemap"
	"github.com/gabrielluizsf/sge/vector"
)

// animalKinds are the kinds of animal a farm may keep, with how many coats
// each may come in.
var animalKinds = map[string]int{"cow": 1, "sheep": 1, "horse": 2}

// oldWorld reports whether the snapshot is from before the farm's animals
// were saved, when the towns had no stables: what it says of the wild may
// lie on ground those have since taken, and is let go, not refused.
func (c *checker) oldWorld() bool { return len(c.s.Herd) == 0 }

// takeHerd is the farm's animals and pens as plain data; the stables' stay
// out, being the same in any game. Nil when its art has no animals.
func takeHerd(f *farm) ([]creatureState, []penState) {
	if f.herd == nil {
		return nil, nil
	}
	pens := []penState{}
	index := map[*pen]int{}
	for _, p := range f.herd.pens {
		if p.town == nil {
			index[p] = len(pens)
			at := f.savedPen(p)
			pens = append(pens, penState{Kind: p.kind, Col: at.col, Row: at.row})
		}
	}
	out := make([]creatureState, 0, len(f.herd.crits))
	for _, c := range f.herd.crits {
		if c.pen.town == nil {
			out = append(out, creatureState{Kind: c.kind, Coat: c.coat, At: point{c.pos.X, c.pos.Y}, Milked: c.milked, Pen: index[c.pen]})
		}
	}
	return out, pens
}

// savedPen is where p is: as it stands, whether it is being moved or not.
func (f *farm) savedPen(p *pen) rect { return p.r }

// checkHerd makes the pens and animals of the snapshot: each pen on ground
// where one could stand, each animal of a kind a farm keeps, in its pen,
// its last milking no later than today's.
func (c *checker) checkHerd() error {
	if len(c.s.Herd) > maxHerd || len(c.s.Pens) > maxHerd {
		return fmt.Errorf("save: %d animals in %d pens", len(c.s.Herd), len(c.s.Pens))
	}
	for _, p := range c.s.Pens {
		r := rect{p.Col, p.Row, penCols, penRows}
		_, known := animalKinds[p.Kind]
		if err := checkCell("pen", p.Col, p.Row); err != nil || !known || !c.penGround(r) {
			return fmt.Errorf("save: a pen of %q at %d, %d: %v", p.Kind, p.Col, p.Row, err)
		}
		c.out.pens = append(c.out.pens, &pen{kind: p.Kind, r: r})
	}
	for _, a := range c.s.Herd {
		coats, ok := animalKinds[a.Kind]
		if !ok {
			return fmt.Errorf("save: an animal of kind %q", a.Kind)
		}
		err := firstError(inRange("coat", a.Coat, 0, coats-1), checkPoint("animal", a.At), inRange("milked", a.Milked, 0, c.s.Day+1), inRange("pen", a.Pen, 0, len(c.out.pens)-1))
		if err != nil {
			return err
		}
		p := c.out.pens[a.Pen]
		pos := keepIn(p.r.inside(), vector.Vector3{X: a.At.X, Y: a.At.Y})
		c.out.herd = append(c.out.herd, &creature{kind: a.Kind, coat: a.Coat, pos: pos, to: pos, milked: a.Milked, wait: 1, heading: 2, pen: p})
	}
	return nil
}

// penGround reports whether a pen may have stood on r: dry and with nothing
// on it but a fence, the first pen's standing where it does in a new game.
func (c *checker) penGround(r rect) bool {
	f := c.f
	for row := r.row; row < r.row+r.rows; row++ {
		for col := r.col; col < r.col+r.cols; col++ {
			id := f.thingAt(col, row)
			fence := id == tileID["fence"] || id == tileID[fenceDown]
			if id != tilemap.Empty && !fence {
				return false
			}
			if !fence && !f.open(col, row) {
				return false
			}
		}
	}
	return true
}

// applyHerd puts the saved pens and animals in the farm's herd, in place of
// its own, their fences with them; the stables' stay. A save with none is
// an old one, and keeps what a new game has.
func applyHerd(f *farm, l *loaded) {
	if f.herd == nil || len(l.herd) == 0 {
		return
	}
	kept := f.herd.crits[:0]
	for _, c := range f.herd.crits {
		if c.pen.town != nil {
			kept = append(kept, c)
		} else if c.obj != nil {
			c.obj.Destroy()
		}
	}
	f.herd.crits = append(kept, l.herd...)
	pens := f.herd.pens[:0]
	for _, p := range f.herd.pens {
		if p.town != nil {
			pens = append(pens, p)
		} else {
			f.fence(p.r, false)
		}
	}
	f.herd.pens = append(pens, l.pens...)
	for _, p := range l.pens {
		f.fence(p.r, true)
	}
}
