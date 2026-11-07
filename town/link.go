package town

// linkRadius is how many regions away a town looks for the town its road
// leads to.
const linkRadius = 3

// minRun is the fewest cells a road keeps between the ends of two streets.
const minRun = 6

// links returns the rectangles of the road leaving town t for the nearest
// town it can reach, worked out once: along t's main street out of its
// end, down or up between the two towns, and along the other's main street
// in, meeting the street's end. With no town to reach it is empty.
func (l *Layout) links(t *Town) []Rect {
	if t.linked {
		return t.road
	}
	t.linked = true
	best := 0
	for dy := -linkRadius; dy <= linkRadius; dy++ {
		for dx := -linkRadius; dx <= linkRadius; dx++ {
			n := l.Of(t.Region.Col+dx, t.Region.Row+dy)
			if n == nil || n == t {
				continue
			}
			d := abs(n.Center.Col-t.Center.Col) + abs(n.Center.Row-t.Center.Row)
			if t.road != nil && d >= best {
				continue
			}
			if p := l.path(t, n); p != nil && l.clear(p, t, n) {
				t.road, t.to, best = p, n, d
			}
		}
	}
	return t.road
}

// path is the road between towns t and n: out of the end of the western
// town's street, to halfway between the two, across to the other's row,
// and in along its street; the same road whichever of them it is asked
// from. It is nil where the streets' ends are too close to take it between.
func (l *Layout) path(t, n *Town) []Rect {
	if n.Center.Col < t.Center.Col {
		t, n = n, t
	}
	a, b := t.Center.Col+t.Reach, n.Center.Col-n.Reach
	if b-a < minRun {
		return nil
	}
	if t.Center.Row == n.Center.Row {
		return []Rect{span(a, t.Center.Row, b, t.Center.Row)}
	}
	mid := a + (b-a)/2
	return []Rect{
		span(a, t.Center.Row, mid, t.Center.Row),
		span(mid, t.Center.Row, mid, n.Center.Row),
		span(mid, n.Center.Row, b, n.Center.Row),
	}
}
