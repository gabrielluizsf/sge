package art

import (
	"fmt"

	"github.com/gabrielluizsf/sge/terrain"
)

// Sprites is every sprite of the farm, by the name of its SVG file.
func Sprites() map[string]Grid {
	s := map[string]Grid{
		"path":            path,
		"soil_dry":        dry,
		"soil_wet":        wet,
		"fence":           fence,
		"stall_roof_l":    stallRoofLeft,
		"stall_roof_r":    stallRoofRight,
		"stall_counter_l": stallCounterLeft,
		"stall_counter_r": stallCounterRight,
		"tree_trunk":      treeTrunk,
		"cursor":          Cursor,
	}
	for i, lean := range []int{0, 1, 0, -1} {
		s[fmt.Sprint("grass", i)] = grass(lean)
		s[fmt.Sprint("flowers", i)] = flowers(lean)
		s[fmt.Sprint("tree_crown", i)] = crown(lean)
		s[fmt.Sprint("water", i)] = water(i)
	}
	for _, mask := range terrain.Masks() {
		for f := range 4 {
			s[fmt.Sprint("shore", mask, "_", f)] = shore(mask, f)
		}
	}
	for i := range 4 {
		s[fmt.Sprint("carrot", i)] = carrot[i]
		s[fmt.Sprint("tomato", i)] = tomato[i]
	}
	s["carrot4"] = sway(carrot[3], 9)
	s["tomato4"] = sway(tomato[3], 6)
	for name, g := range Farmer {
		s["farmer_"+name] = g
	}
	for name, g := range Icons {
		s["icon_"+name] = g
	}
	return s
}
