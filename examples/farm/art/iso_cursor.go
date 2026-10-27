package art

// IsoCursor marks the cell the farmer works: its diamond, outlined, open
// in the middle so the tile shows. It is a cell's diamond, IsoW by Size.
func IsoCursor() Grid {
	r := &raster{g: blank(IsoW, Size), ox: Size, oy: 0}
	corners := []point{{0.5, 0.5, 0}, {Size - 0.5, 0.5, 0}, {Size - 0.5, Size - 0.5, 0}, {0.5, Size - 0.5, 0}}
	for i, c := range corners {
		r.line('w', c, corners[(i+1)%4])
	}
	return r.g
}
