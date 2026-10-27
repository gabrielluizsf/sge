package art

import (
	"fmt"

	"github.com/gabrielluizsf/sge/terrain"
)

// Sprites is every sprite of the farm, as SVG, by the name of its file.
// The tiles are IsoW by IsoH, the ground laid on a cell's diamond and what
// stands on it upright; the stall is StallW by StallH, the cursor a cell's
// diamond, the farmer FarmerW by FarmerH and the icons Size square.
func Sprites() map[string]string {
	s := map[string]string{
		"path":     pathTile(),
		"soil_dry": soilTile(soilLight, soilBase, soilShade, false),
		"soil_wet": soilTile(wetLight, wetBase, wetShade, true),
		"fence":    fenceTile(),
		// What stands where nothing shows: the stall's cells, solid.
		"block":  newDoc(IsoW, IsoH).String(),
		"stall":  Stall(),
		"cursor": Cursor(),
	}
	for i, lean := range []float64{0, 0.6, 0, -0.6} {
		s[fmt.Sprint("grass", i)] = grassTile(lean)
		s[fmt.Sprint("flowers", i)] = flowersTile(lean)
		s[fmt.Sprint("tree", i)] = treeTile(lean * 0.8)
	}
	for _, mask := range terrain.Masks() {
		for f := range 4 {
			s[fmt.Sprint("shore", mask, "_", f)] = shoreTile(mask, f)
		}
	}
	for i := range 4 {
		s[fmt.Sprint("carrot", i)] = carrotTile(i, 0)
		s[fmt.Sprint("tomato", i)] = tomatoTile(i, 0)
	}
	s["carrot4"] = carrotTile(3, 0.9)
	s["tomato4"] = tomatoTile(3, 0.9)
	for name, svg := range FarmerFrames() {
		s["farmer_"+name] = svg
	}
	for name, svg := range Icons() {
		s["icon_"+name] = svg
	}
	return s
}
