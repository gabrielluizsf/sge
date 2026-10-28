package terrain

// In reports whether cell col, row is in the area called name, so one
// area can keep away from another: trees from water.
func (g *Generator) In(name string, col, row int) bool {
	for _, a := range g.Areas {
		if a.Name == name && a.In(col, row) {
			return true
		}
	}
	return false
}
