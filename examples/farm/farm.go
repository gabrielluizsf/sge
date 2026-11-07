package farm

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/backpack"
	"github.com/gabrielluizsf/sge/daytime"
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/life"
	"github.com/gabrielluizsf/sge/lineage"
	"github.com/gabrielluizsf/sge/particle"
	"github.com/gabrielluizsf/sge/progress"
	"github.com/gabrielluizsf/sge/random"
	"github.com/gabrielluizsf/sge/request"
	"github.com/gabrielluizsf/sge/scene"
	"github.com/gabrielluizsf/sge/terrain"
	"github.com/gabrielluizsf/sge/tilemap"
	"github.com/gabrielluizsf/sge/toast"
	"github.com/gabrielluizsf/sge/town"
	"github.com/gabrielluizsf/sge/vitals"
	"github.com/gabrielluizsf/sge/warm"
)

// The hands: a slot of the hotbar holds seeds, of the crop chosen, the watering
// can, the axe, the pickaxe or the hoe.
const (
	slotSeeds = iota
	slotCan
	slotAxe
	slotHoe
	slotFood // a crop, eaten when used
	slotPick // the pickaxe, breaking rock
)

// farm is the game's state: the player's kit, the plots of soil on the
// farm and in the towns' gardens, the towns and their people, the trees
// being cut and the stumps growing back, and what is open on screen.
type farm struct {
	layers
	*kit
	// art draws every picture of the game.
	art   Art
	tiles *tilemap.TileSet
	// pack is the player's inventory screen and hotbar; ui the menus.
	pack   *backpack.Pack
	ui     *ui
	wild   *terrain.World
	towns  *town.Layout
	keeper *town.Keeper
	// warm draws the pictures of what is about to come into view, a few
	// a frame, ahead of the frame that shows them.
	warm   *warm.Queue
	plots  []*plot
	shop   *shop
	talk   *talk
	floats []*floater
	// vit is the player's health and food; flash how long is left of the
	// red edge of the screen after a blow; blood the specks of blood and
	// bloodObj and stainObj the objects that show them, in the air and on
	// the ground.
	vit      vitals.Vitals
	flash    float64
	blood    *particle.System
	bloodObj *game.Object
	stainObj *game.Object
	// cues are the sounds due in a moment, as a blow lands.
	cues []cue
	// people are the villagers of every town woken so far, asleep or not.
	people  []*villager
	settled map[*town.Town]bool
	// away is when, in hours since the first day began, each town that is
	// asleep was last seen, to settle what its people did meanwhile.
	away map[*town.Town]float64
	// roads are the roads between the towns, worked out as asked; trafficIn
	// is the seconds until the travellers are looked at again.
	roads     map[*town.Town]*road
	trafficIn float64
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
	stage *scene.Stage
	room  *room
	// home is the furniture of the player's house.
	home []piece
	// yard are the chests out of doors.
	yard []*yardChest
	// edit is the arranging of furniture, if on; seat the chair the player sits on.
	edit *editor
	seat *seat
	// herd is the farm's animals, if its art draws them; milking is the
	// player at a cow.
	herd    *herd
	milking *milking
	// wolves are those that come at night.
	wolves    *wolfPack
	icons     map[string]*canvas.Canvas
	tints     map[string]*canvas.Canvas
	residents []resident
	notices   []*toast.Toast
	// prog is the family's level and XP, kept by whoever plays; shown the
	// level it has reached as the bar draws it, and glow how long is left
	// of the glow after a level up.
	prog  progress.Progress
	shown float64
	glow  float64
	// speaker plays what happens on the farm; nil plays nothing.
	speaker speaker
	// self is who the player is now; family the tree of everyone's
	// parents, the player's line and the villagers'; fate draws who is
	// born; ids is how many ids it has given; seen the last day the farm
	// has lived through; question the yes or no being asked, if any;
	// drawn the time of life the player is drawn at.
	self     *self
	family   *lineage.Tree
	fate     *random.Source
	ids      int
	seen     int
	question *question
	drawn    life.Stage
}

// newFarm is a farm with some money and nothing growing on the soil of
// the ground layer, in a world of towns.
func newFarm(l layers, w world) *farm {
	f := &farm{layers: l, kit: &kit{money: 20, axe: noAxe, exact: true}, art: w.art, tiles: w.tiles, wild: w.wild, towns: w.towns, shop: &shop{},
		settled: map[*town.Town]bool{}, roads: map[*town.Town]*road{}, hits: map[int]int{}, luck: random.New(seed),
		clock: &daytime.Clock{Hour: startHour, Speed: 24.0 / dayLength}, fainted: -1, prog: progress.New(), vit: vitals.Full()}
	f.home = defaultHouse()
	f.pack = newPack(f)
	f.hold(0)
	f.errand = errandBook.Draw(f.luck)
	f.family, f.fate = &lineage.Tree{}, random.New(seed+101)
	f.self = f.found()
	f.self.name = playerName(w.art)
	f.drawn = f.self.who(f.clock.Day).Stage
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
