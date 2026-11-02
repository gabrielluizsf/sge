package tilemap

// refresh paints, in the cells of r, the frame of the map's time of each
// that turns with it.
func (m *Map) refresh(r cellRange) {
	for row := r.r0; row <= r.r1; row++ {
		for col := r.c0; col <= r.c1; col++ {
			if m.Tiles.turns(m.cells[row*m.Cols+col]) {
				m.paint(col, row)
			}
		}
	}
}

// catchUp brings the cells of r that the last step did not reach, as the
// camera has moved since, to the map's time: a cell shows the frame of the
// time of the step it is drawn in, whenever it came into view.
func (m *Map) catchUp(r cellRange) {
	if m.Tiles.turning == 0 || m.fresh.covers(r) {
		return
	}
	for row := r.r0; row <= r.r1; row++ {
		for col := r.c0; col <= r.c1; col++ {
			if !m.fresh.has(col, row) && m.Tiles.turns(m.cells[row*m.Cols+col]) {
				m.paint(col, row)
			}
		}
	}
	if r.covers(m.fresh) {
		m.fresh = r
	}
}
