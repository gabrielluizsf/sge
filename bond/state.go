package bond

// State is a bond's numbers as plain data, to keep and bring back: the
// affection, how far the courting went, the days after the last gift and
// talk, and the wedding day.
type State struct {
	Affection int
	Courted   Stage
	Gifted    int
	Talked    int
	Wedding   int
}

// State is b as plain data.
func (b Bond) State() State {
	return State{Affection: b.Affection, Courted: b.Courted, Gifted: b.gifted, Talked: b.talked, Wedding: b.wedding}
}

// Bond is the bond s was taken from.
func (s State) Bond() Bond {
	return Bond{Affection: s.Affection, Courted: s.Courted, gifted: s.Gifted, talked: s.Talked, wedding: s.Wedding}
}
