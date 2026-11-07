package farm

import "fmt"

// The most a save may hold of furniture.
const (
	maxPieces = roomCols * roomRows
	maxYard   = 500
)

// checkHouse is nil if the furniture and the chests of the save are ones
// the game could have: known pieces turned one of four ways, a house they
// fit in with a bed, and chests that keep known items, a stack a slot.
func (c *checker) checkHouse() error {
	s := c.s
	if len(s.Home) > maxPieces || len(s.Yard) > maxYard {
		return fmt.Errorf("save: %d pieces in the house, %d chests out of doors", len(s.Home), len(s.Yard))
	}
	home, yard, err := checkedPieces(s)
	if err != nil {
		return err
	}
	if len(home) > 0 {
		if !houseFits(c.f.art, home) || !hasBed(home) {
			return fmt.Errorf("save: the furniture does not fit the house, or it has no bed")
		}
	}
	seen := map[[2]int]bool{}
	for _, p := range yard {
		if p.name != "chest" || seen[[2]int{p.col, p.row}] {
			return fmt.Errorf("save: bad chest out of doors at %d,%d", p.col, p.row)
		}
		seen[[2]int{p.col, p.row}] = true
	}
	c.out.home, c.out.yard = home, yard
	return nil
}

// hasBed reports whether one of pieces is a bed.
func hasBed(pieces []piece) bool {
	for _, p := range pieces {
		if p.name == "bed" {
			return true
		}
	}
	return false
}

// checkedPieces are the pieces of the save, each checked alone.
func checkedPieces(s *snapshot) (home, yard []piece, err error) {
	for _, list := range []struct {
		what string
		in   []pieceState
		out  *[]piece
	}{{"furniture", s.Home, &home}, {"chest", s.Yard, &yard}} {
		for _, st := range list.in {
			if err := checkPiece(list.what, st); err != nil {
				return nil, nil, err
			}
		}
		*list.out = piecesOf(list.in)
	}
	return home, yard, nil
}

// checkPiece is nil if st is a piece the game could have.
func checkPiece(what string, st pieceState) error {
	if !isFurnishing(st.Name) {
		return fmt.Errorf("save: %s %q is no piece of furniture", what, st.Name)
	}
	if err := firstError(inRange(what+" heading", st.Heading, 0, headings-1), checkCell(what, st.Col, st.Row)); err != nil {
		return err
	}
	if st.Name != "chest" {
		if len(st.Held) != 0 {
			return fmt.Errorf("save: a %s keeps things", st.Name)
		}
		return nil
	}
	if len(st.Held) != chestSlots {
		return fmt.Errorf("save: a chest has %d slots, not %d", len(st.Held), chestSlots)
	}
	seen := map[string]bool{}
	for _, h := range st.Held {
		if h.N == 0 && h.Item == "" {
			continue
		}
		if !knownItem(h.Item) || seen[h.Item] {
			return fmt.Errorf("save: a chest keeps a bad item %q", h.Item)
		}
		seen[h.Item] = true
		if err := inRange("chest "+h.Item, h.N, 1, maxCount); err != nil {
			return err
		}
	}
	return nil
}
