package farm

// slotState is a slot of a chest: an item and how many, none when n is 0.
type slotState struct {
	Item string `json:"item"`
	N    int    `json:"n"`
}

// pieceState is a piece of furniture or a chest out of doors: what it is,
// the cell it stands on, which way it is turned, and, for a chest, what it
// keeps.
type pieceState struct {
	Name    string      `json:"name"`
	Col     int         `json:"col"`
	Row     int         `json:"row"`
	Heading int         `json:"heading"`
	Held    []slotState `json:"held"`
}

// takeHouse is the house's furniture as plain data, with the piece being
// moved, if the house's, back where it was.
func takeHouse(f *farm) []pieceState {
	ps := f.home
	if p, ok := f.lifting(); ok && f.room != nil {
		ps = append(append([]piece(nil), ps...), p)
	}
	return takePieces(ps)
}

// takePieces are pieces as plain data, never nil.
func takePieces(pieces []piece) []pieceState {
	out := make([]pieceState, 0, len(pieces))
	for _, p := range pieces {
		s := pieceState{Name: p.name, Col: p.col, Row: p.row, Heading: p.heading, Held: []slotState{}}
		for _, h := range p.held {
			s.Held = append(s.Held, slotState{Item: h.item, N: h.n})
		}
		out = append(out, s)
	}
	return out
}

// takeYard are the chests out of doors as plain data.
func takeYard(f *farm) []pieceState {
	ps := make([]piece, 0, len(f.yard)+1)
	for _, y := range f.yard {
		ps = append(ps, y.p)
	}
	if p, ok := f.lifting(); ok && f.room == nil {
		ps = append(ps, p)
	}
	return takePieces(ps)
}

// piecesOf are s as pieces, which were checked.
func piecesOf(s []pieceState) []piece {
	out := make([]piece, 0, len(s))
	for _, st := range s {
		p := newPiece(st.Name, st.Col, st.Row, st.Heading)
		for i, h := range st.Held {
			p.held[i] = chestSlot{item: h.Item, n: h.N}
		}
		out = append(out, p)
	}
	return out
}

// applyHouse sets the player's furniture and the chests out of doors, the
// chests the save had replacing those that were. A save with no furniture
// is an old one, and its house has only the bed.
func applyHouse(f *farm, s *snapshot, l *loaded) {
	f.home = defaultHouse()
	if len(l.home) > 0 {
		f.home = l.home
	}
	for len(f.yard) > 0 {
		f.takeChest(f.yard[0])
	}
	for _, p := range l.yard {
		f.putChest(p)
	}
	f.edit, f.seat = nil, nil
}
