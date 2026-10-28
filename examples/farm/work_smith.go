package main

import (
	"fmt"

	"github.com/gabrielluizsf/sge/pathfind"
	"github.com/gabrielluizsf/sge/town"
)

// toolPrice is what a smith earns for each tool forged, and forgeStrokes
// how many blows of the hammer one takes, burning a load of wood.
const (
	toolPrice    = 6
	forgeStrokes = 3
)

// smithWork is the smith's work at the smithy: forging at the anvil while
// the forge burns; feeding it the wood they have when it burns low; and
// when they have none, buying the lumberjack's, or felling a tree.
func (m *mind) smithWork(f *farm) *job {
	v := m.v
	k := v.kit
	t := v.town
	switch {
	case v.fuel >= 1:
		return &job{to: pathfind.Cell{Col: t.Anvil.Col, Row: t.Anvil.Row + 1}, face: "up", slot: slotAxe, uses: forgeStrokes, rest: 1, then: func(*farm) {
			v.fuel--
			k.money += toolPrice
		}}
	case k.bag.Count(wood) > 0:
		return &job{to: pathfind.Cell{Col: t.Forge.Col, Row: t.Forge.Row + 1}, face: "up", slot: -1, rest: 1, then: func(f *farm) {
			n := k.bag.Count(wood)
			k.bag.Take(wood, n)
			v.fuel += float64(n)
			at := f.center(pathfind.Cell(t.Forge))
			f.float(at, fmt.Sprint("+", n, " LENHA"), 'y')
		}}
	}
	if lj := f.lumberjackWith(t, 2); lj != nil && k.money >= 2*woodPrice {
		col, row := f.cellOf(lj.pos)
		return &job{to: pathfind.Cell{Col: col, Row: row}, slot: -1, rest: 1, then: func(f *farm) { f.buyWood(v, lj) }}
	}
	return m.chopJob(f)
}

// lumberjackWith is a lumberjack of town t, awake and out, with at least
// n wood, or nil.
func (f *farm) lumberjackWith(t *town.Town, n int) *villager {
	for _, v := range f.people {
		if v.town == t && v.role == roleLumberjack && !v.inside && v.kit.bag.Count(wood) >= n {
			return v
		}
	}
	return nil
}

// buyWood is smith s buying what wood they can afford from lumberjack lj,
// if they are still near: up to what lj has.
func (f *farm) buyWood(s, lj *villager) {
	if s.pos.Sub(lj.pos).Length() > 2.5*tile {
		return
	}
	n := min(lj.kit.bag.Count(wood), s.kit.money/woodPrice)
	if n <= 0 || !lj.kit.bag.Take(wood, n) {
		return
	}
	s.kit.bag.Add(wood, n)
	s.kit.money -= n * woodPrice
	lj.kit.money += n * woodPrice
	f.float(lj.pos, fmt.Sprint(n, " MADEIRA $", n*woodPrice), 'y')
	f.sound("coin", &lj.pos)
}
