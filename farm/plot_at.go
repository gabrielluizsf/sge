package farm

// plotAt is the plot at cell col, row, or nil where there is no soil.
func (f *farm) plotAt(col, row int) *plot {
	for _, p := range f.plots {
		if p.col == col && p.row == row {
			return p
		}
	}
	return nil
}
