package farm

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/daytime"
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/random"
	"github.com/gabrielluizsf/sge/request"
	"github.com/gabrielluizsf/sge/scene"
	"github.com/gabrielluizsf/sge/terrain"
	"github.com/gabrielluizsf/sge/tilemap"
	"github.com/gabrielluizsf/sge/town"
)

// The hands: a slot of the hotbar holds seeds of a crop, the watering
// can, or the axe.
const (
	slotCarrot = iota
	slotTomato
	slotCan
	slotAxe
	slots
)

// farm is the game's state: the player's kit, the plots of soil on the
// farm and in the towns' gardens, the towns and their people, the trees
// being cut and the stumps growing back, and what is open on screen.
type farm struct {
	layers
	*kit
	// art draws every picture of the game.
	art    Art
	tiles  *tilemap.TileSet
	wild   *terrain.World
	towns  *town.Layout
	keeper *town.Keeper
	plots  []*plot
	shop   *shop
	talk   *talk
	floats []*floater
	// people are the villagers of every town woken so far, asleep or not.
	people  []*villager
	settled map[*town.Town]bool
	// hits are the strokes each tree being cut has taken, by cellKey.
	hits   map[int]int
	stumps []*stump
	// errand is what the smith asks for the stone axe; done once given.
	errand     request.Order
	errandDone bool
	luck       *random.Source
	player     *game.Object
	hands      *farmer
	// clock is the time of day; fainted the last day the player fainted
	// for staying up.
	clock   *daytime.Clock
	fainted int
	// stage switches between the world outside and the room the player is
	// in, with whoever lives there.
	stage     *scene.Stage
	room      *room
	residents []resident
	notices   []*notice
	// speaker plays what happens on the farm; nil plays nothing.
	speaker speaker
}

// newFarm is a farm with some money and nothing growing on the soil of
// the ground layer, in a world of towns.
func newFarm(l layers, w world) *farm {
	f := &farm{layers: l, kit: &kit{money: 20, axe: noAxe}, art: w.art, tiles: w.tiles, wild: w.wild, towns: w.towns, shop: &shop{},
		settled: map[*town.Town]bool{}, hits: map[int]int{}, luck: random.New(seed),
		clock: &daytime.Clock{Hour: startHour, Speed: 24.0 / dayLength}, fainted: -1}
	f.errand = errandBook.Draw(f.luck)
	for row := 0; row < l.ground.Rows; row++ {
		for col := 0; col < l.ground.Cols; col++ {
			if l.ground.At(col, row) == tileID["soil_dry"] {
				f.plots = append(f.plots, &plot{col: col, row: row, crop: -1})
			}
		}
	}
	return f
}

// The day starts at startHour, and a whole day lasts dayLength seconds.
const (
	startHour = 7.0
	dayLength = 12 * 60
)

// roomDark is what is behind a room's walls.
var roomDark = canvas.RGB(28, 20, 30)

// resident is a villager at home, seen in their room.
type resident struct {
	v   *villager
	obj *game.Object
}
