package farm

// What villagers may take of the player's farm: how many plots each may
// keep, the share of all the farm's soil they may hold together, and how
// many come to work it, or to fell its trees, a day.
const (
	fieldEach  = 3
	fieldShare = 0.4
	fieldHands = 2
)

// farmPlots are the plots of the farm's own soil.
func (f *farm) farmPlots() []*plot {
	var out []*plot
	for _, p := range f.plots {
		if p.town == nil {
			out = append(out, p)
		}
	}
	return out
}

// ownedBy is how many plots of the farm v keeps; held how many all
// villagers do.
func (f *farm) ownedBy(v *villager) (n int) {
	for _, p := range f.farmPlots() {
		if p.owner == v {
			n++
		}
	}
	return n
}

// held is how many plots of the farm villagers keep.
func (f *farm) held() (n int) {
	for _, p := range f.farmPlots() {
		if p.owner != nil {
			n++
		}
	}
	return n
}

// roomInField reports whether v may take one plot more: they keep fewer
// than fieldEach, and all together fewer than fieldShare of the soil.
func (f *farm) roomInField(v *villager) bool {
	return f.ownedBy(v) < fieldEach && float64(f.held()+1) <= fieldShare*float64(len(f.farmPlots()))
}
