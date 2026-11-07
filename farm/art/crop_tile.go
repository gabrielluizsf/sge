package art

// CropNames are the crops drawn, in the order the farm lists them.
var CropNames = []string{"carrot", "lettuce", "potato", "tomato", "beans", "corn", "wheat", "pumpkin", "cassava",
	"strawberry", "watermelon", "pepper", "sugarcane", "sunflower", "coffee"}

// cropTiles draw each crop at a growth.
var cropTiles = map[string]func(growth) string{
	"carrot":     func(g growth) string { return carrotTile(g.stage, g.lean) },
	"tomato":     func(g growth) string { return tomatoTile(g.stage, g.lean) },
	"lettuce":    lettuceTile,
	"potato":     potatoTile,
	"beans":      beansTile,
	"corn":       cornTile,
	"wheat":      wheatTile,
	"pumpkin":    pumpkinTile,
	"cassava":    cassavaTile,
	"strawberry": strawberryTile,
	"watermelon": watermelonTile,
	"pepper":     pepperTile,
	"sugarcane":  sugarcaneTile,
	"sunflower":  sunflowerTile,
	"coffee":     coffeeTile,
}

// CropTile is the SVG of the crop name at stage, 0 the seed to 3 ripe,
// leant by lean in the wind, standing on its cell.
func CropTile(name string, stage int, lean float64) string {
	return cropTiles[name](growth{stage, lean})
}
