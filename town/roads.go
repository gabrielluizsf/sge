package town

// roadPad is how far, in cells, from a region the roads that reach it are
// kept: the most a query may ask for.
const roadPad = 8

// roadsOf is the rectangles of the roads within roadPad cells of region
// rx, ry, worked out once. The towns that road may start from lie within
// linkRadius regions, plus one for the pad.
func (l *Layout) roadsOf(rx, ry int) []Rect {
	if l.lastOK && l.lastRx == rx && l.lastRy == ry {
		return l.last
	}
	k, ok := key(Cell{rx, ry})
	if !ok {
		return nil
	}
	rects, done := l.roads[k]
	if !done {
		rects = l.gather(rx, ry)
		if l.roads == nil {
			l.roads = map[int][]Rect{}
		}
		l.roads[k] = rects
	}
	l.lastRx, l.lastRy, l.last, l.lastOK = rx, ry, rects, true
	return rects
}

// gather works out the roads near region rx, ry.
func (l *Layout) gather(rx, ry int) []Rect {
	box := Rect{rx * l.Region, ry * l.Region, (rx+1)*l.Region - 1, (ry+1)*l.Region - 1}
	var out []Rect
	for _, r := range l.Roads {
		if r.Within(box, roadPad) {
			out = append(out, r)
		}
	}
	reach := linkRadius + 1
	for y := ry - reach; y <= ry+reach; y++ {
		for x := rx - reach; x <= rx+reach; x++ {
			if t := l.Of(x, y); t != nil {
				for _, r := range l.links(t) {
					if r.Within(box, roadPad) && !has(out, r) {
						out = append(out, r)
					}
				}
			}
		}
	}
	return out
}

// OnRoad reports whether cell col, row is on a road between towns, or
// within margin cells of one, at most roadPad.
func (l *Layout) OnRoad(col, row, margin int) bool {
	for _, r := range l.roadsOf(floorDiv(col, l.Region), floorDiv(row, l.Region)) {
		if r.Dist(col, row) <= margin {
			return true
		}
	}
	return false
}

// has reports whether rects holds r: two towns that lead to each other
// lay one road.
func has(rects []Rect, r Rect) bool {
	for _, o := range rects {
		if o == r {
			return true
		}
	}
	return false
}
