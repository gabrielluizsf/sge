package town

// Route is the way on foot from the middle of town t's plaza to the
// middle of the plaza of the town its road leads to, as the cells where
// it turns, and that town; nil when t has no road. It is the same way,
// reversed, from the other town.
func (l *Layout) Route(t *Town) (*Town, []Cell) {
	if l.links(t) == nil || t.to == nil {
		return nil, nil
	}
	n := t.to
	w, e := t, n
	if n.Center.Col < t.Center.Col {
		w, e = n, t
	}
	a, b := w.Center.Col+w.Reach, e.Center.Col-e.Reach
	pts := []Cell{w.Center, {a, w.Center.Row}}
	if w.Center.Row != e.Center.Row {
		mid := a + (b-a)/2
		pts = append(pts, Cell{mid, w.Center.Row}, Cell{mid, e.Center.Row})
	}
	pts = append(pts, Cell{b, e.Center.Row}, e.Center)
	if t != w {
		for i, j := 0, len(pts)-1; i < j; i, j = i+1, j-1 {
			pts[i], pts[j] = pts[j], pts[i]
		}
	}
	return n, pts
}
