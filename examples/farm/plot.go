package farm

// plot is a cell of soil: what grows on it, how far, and whether it has
// been watered since its last stage.
type plot struct {
	col, row int
	crop     int // an index of crops, or -1
	stage    int
	watered  bool
	growing  float64 // seconds watered at this stage
}

// thirsty reports whether watering would make the plot's crop grow.
func (p *plot) thirsty() bool {
	return p.crop >= 0 && p.stage < ripe && !p.watered
}
