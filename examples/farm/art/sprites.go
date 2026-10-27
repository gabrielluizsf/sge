package art

import (
	"fmt"

	"github.com/gabrielluizsf/sge/terrain"
)

// Sprites is every sprite of the farm, by the name of its SVG file. The
// tiles are IsoW by IsoH, the ground laid on a cell's diamond and what
// stands on it upright; the stall is a picture of its own, the cursor a
// cell's diamond, and the farmer and the icons Size square.
func Sprites() map[string]Grid {
	s := map[string]Grid{
		"path":     isoGround(path),
		"soil_dry": isoGround(dry),
		"soil_wet": isoGround(wet),
		"fence":    isoFence(),
		// What stands where nothing shows: the stall's cells, solid.
		"block":  blank(IsoW, IsoH),
		"stall":  Stall(),
		"cursor": IsoCursor(),
	}
	for i, lean := range []int{0, 1, 0, -1} {
		s[fmt.Sprint("grass", i)] = isoGround(grass(lean))
		s[fmt.Sprint("flowers", i)] = isoGround(flowers(lean))
		s[fmt.Sprint("tree", i)] = standing(append(crown(lean).copyOf(), treeTrunk...), true)
	}
	for _, mask := range terrain.Masks() {
		for f := range 4 {
			s[fmt.Sprint("shore", mask, "_", f)] = isoGround(shore(mask, f))
		}
	}
	for i := range 4 {
		s[fmt.Sprint("carrot", i)] = standing(carrot[i], false)
		s[fmt.Sprint("tomato", i)] = standing(tomato[i], false)
	}
	s["carrot4"] = standing(sway(carrot[3], 9), false)
	s["tomato4"] = standing(sway(tomato[3], 6), false)
	for name, g := range Farmer {
		s["farmer_"+name] = g
	}
	for name, g := range Icons {
		s["icon_"+name] = g
	}
	return s
}
