package town

// Road reports whether cell col, row is paved: the main street or the
// plaza.
func (t *Town) Road(col, row int) bool {
	dx, dy := col-t.Center.Col, row-t.Center.Row
	return dy == 0 && abs(dx) <= t.Reach || abs(dx) <= plazaW && abs(dy) <= plazaH
}

// House is the house standing on cell col, row, or nil.
func (t *Town) House(col, row int) *House {
	for i := range t.Houses {
		h := &t.Houses[i]
		if col >= h.Col && col < h.Col+h.Width && row >= h.Row && row < h.Row+h.Depth {
			return h
		}
	}
	return nil
}

// Grounds reports whether cell col, row is within margin cells of the
// town, its street's ends and its houses' backs: where the wild keeps
// out of it.
func (t *Town) Grounds(col, row, margin int) bool {
	dx, dy := col-t.Center.Col, row-t.Center.Row
	if abs(dx) <= t.Reach+margin && abs(dy) <= t.Depth+margin {
		return true
	}
	return t.Back != (Rect{}) && t.Back.Dist(col, row) <= margin
}

// Grounds reports whether cell col, row is on the grounds of the town of
// its region, within margin.
func (l *Layout) Grounds(col, row, margin int) bool {
	t := l.At(col, row)
	return t != nil && t.Grounds(col, row, margin)
}
