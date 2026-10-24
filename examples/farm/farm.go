package main

import "github.com/gabrielluizsf/sge/terrain"

// The farmer's hands: a slot of the hotbar holds seeds of a crop or the
// watering can.
const (
	slotCarrot = iota
	slotTomato
	slotCan
	slots
)

// farm is the game's state: the money, the seeds and harvest of each
// crop, what the farmer holds, and the plots of soil.
type farm struct {
	layers
	wild    *terrain.World
	money   int
	seeds   []int
	harvest []int
	slot    int
	plots   []*plot
	shop    *shop
	floats  []*floater
	// sold is everything sold over the game, for the demo to report.
	sold int
}

// newFarm is a farm with some money and nothing growing on the soil of
// the ground layer.
func newFarm(l layers, wild *terrain.World) *farm {
	f := &farm{layers: l, wild: wild, shop: &shop{}, money: 20, seeds: make([]int, len(crops)), harvest: make([]int, len(crops))}
	for row := 0; row < l.ground.Rows; row++ {
		for col := 0; col < l.ground.Cols; col++ {
			if l.ground.At(col, row) == tileID["soil_dry"] {
				f.plots = append(f.plots, &plot{col: col, row: row, crop: -1})
			}
		}
	}
	return f
}
