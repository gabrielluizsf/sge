package farm

import "github.com/gabrielluizsf/sge/lineage"

// takePlace is where the player is: where they stand, or, in a room, at
// its door outside.
func takePlace(f *farm) point {
	if r := f.room; r != nil {
		p := f.ground.CellCenter(r.door.Col, r.door.Row)
		return point{X: p.X, Y: p.Y}
	}
	p := f.player.Transform.Position
	return point{X: p.X, Y: p.Y}
}

// takeSelf is who the player is.
func takeSelf(f *farm) selfState {
	s := f.self
	return selfState{
		ID: s.id, Name: s.name, Woman: s.woman, Genes: s.genes,
		Life: lifeState{Born: s.life.Born, Span: s.life.Span}, Spouse: indexOfPerson(f, s.spouse),
		Expecting: takeExpecting(s.expecting),
	}
}

// takeExpecting is e as plain data, nil for none.
func takeExpecting(e *lineage.Expecting) *expectState {
	if e == nil {
		return nil
	}
	return &expectState{Parents: e.Parents, Due: e.Due}
}

// indexOfPerson is where v is among the farm's people, or -1.
func indexOfPerson(f *farm, v *villager) int {
	if v == nil {
		return -1
	}
	for i, p := range f.people {
		if p == v {
			return i
		}
	}
	return -1
}

// takeFamily is the family tree, each person after their parents.
func takeFamily(f *farm) []kinState {
	es := f.family.Entries()
	out := make([]kinState, 0, len(es))
	for _, e := range es {
		out = append(out, kinState{ID: e.ID, Parents: e.Parents, Born: e.Born, Dead: e.Dead})
	}
	return out
}
