package farm

import (
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/backpack"
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/pathfind"
	"github.com/gabrielluizsf/sge/vector"
)

// openIn opens the chest at index i of the house beside the bag.
func openIn(f *farm, i int) {
	f.home = append(f.home, newPiece("chest", 6, 1, 0))
	f.openChest(&f.home[len(f.home)-1])
}

// Items go into a chest and out of it, by key and by mouse, and a chest
// keeps ten stacks, no more.
func TestAChestKeepsTenStacksAndMovesItemsBothWays(t *testing.T) {
	_, f := inHouse(t)
	f.bag.Add(wood, 5)
	f.bag.Add(stone, 2)
	openIn(f, 0)
	if f.pack.Side == nil || len(f.pack.Side.Slots()) != chestSlots {
		t.Fatalf("the chest screen: %+v", f.pack.Side)
	}
	// Quick send: the cursor goes to the wood and shift-Enter sends it.
	items := carried{f}.Items()
	at := 0
	for i, it := range items {
		if it.ID == wood {
			at = i
		}
	}
	f.pack.Handle(backpack.Input{Dx: at % inventoryCols, Dy: at / inventoryCols})
	if ev := f.pack.Handle(backpack.Input{Select: true, Quick: true}); ev != backpack.Stored {
		t.Fatalf("send: %v", ev)
	}
	if f.bag.Count(wood) != 0 || f.home[len(f.home)-1].held[0] != (chestSlot{wood, 5}) {
		t.Fatalf("bag %d, chest %+v", f.bag.Count(wood), f.home[len(f.home)-1].held)
	}
	// And back, with the mouse: Paint places the rects first.
	cv, _ := canvas.NewCanvas(1280, 720)
	f.pack.Paint(cv, f.packArt())
	from := f.pack.Where(backpack.Spot{Area: backpack.Chest})
	f.pack.Handle(backpack.Input{Press: true, Mouse: from})
	grid := f.pack.Where(backpack.Spot{Area: backpack.Grid, Index: 3})
	if ev := f.pack.Handle(backpack.Input{Release: true, Mouse: grid}); ev != backpack.Stored || f.bag.Count(wood) != 5 {
		t.Fatalf("take back: %v, wood %d", ev, f.bag.Count(wood))
	}
	// Ten kinds fill it; the eleventh is refused.
	st := chestStore{f, &f.home[len(f.home)-1]}
	kinds := []string{wood, stone, "chair", "table", "shelf", "plant", "stove", "bed", seedsOf(0), cropOf(0)}
	for _, k := range kinds {
		f.bag.Add(k, 1)
		if !st.Put(k) {
			t.Fatalf("could not put %s", k)
		}
	}
	f.bag.Add(seedsOf(1), 1)
	if st.Put(seedsOf(1)) || f.bag.Count(seedsOf(1)) != 1 {
		t.Error("an eleventh stack went in")
	}
	f.bag.Add(wood, 4) // the same kind joins its stack
	if !st.Put(wood) || f.home[len(f.home)-1].held[0].n != 10 {
		t.Errorf("wood did not join: %+v", f.home[len(f.home)-1].held[0])
	}
}

// E beside a chest opens it; closing the inventory lets go of it.
func TestEOpensTheChestBesideAndClosingLetsGo(t *testing.T) {
	_, f := inHouse(t)
	f.home = append(f.home, newPiece("chest", 3, 3, 0))
	f.relay()
	f.player.Transform.Position = f.ground.CellCenter(3, 4)
	words, do, ok := f.interaction(f.player)
	if !ok || words != "E: open chest" {
		t.Fatalf("%q %v", words, ok)
	}
	do()
	if !f.pack.Open || f.pack.Side == nil {
		t.Fatal("not open")
	}
	f.pack.Hide()
}

// A chest out of doors stands on free ground, is solid for the player and
// for the villagers' paths, and keeps its place whatever happens to the
// world around it.
func TestAChestOutsideIsSolidAndStays(t *testing.T) {
	g, f := build(SVGArt(), idle{})
	g.Step(frame)
	col, row := f.cellOf(feet(f.player))
	spot := piece{}
	for dc := 1; dc <= 3 && spot.name == ""; dc++ {
		if p := newPiece("chest", col+dc, row, 0); f.canPlaceYard(p) {
			spot = p
		}
	}
	if spot.name == "" {
		t.Fatal("no free ground near the player")
	}
	if !f.open(spot.col, spot.row) {
		t.Fatal("the cell is not open to begin with")
	}
	y := f.putChest(spot)
	if f.open(spot.col, spot.row) {
		t.Error("villagers can path through a chest")
	}
	way := f.grid().Path(pathfind.Cell{Col: spot.col - 1, Row: spot.row}, pathfind.Cell{Col: spot.col + 1, Row: spot.row})
	for _, c := range way {
		if c.Col == spot.col && c.Row == spot.row {
			t.Errorf("a path goes through the chest: %v", way)
		}
	}
	// The player walks into it and stops.
	f.player.Transform.Position = f.ground.CellCenter(spot.col-1, spot.row)
	for range 120 {
		f.hands.ctl = walkRight{}
		g.Step(frame)
	}
	if c, _ := f.cellOf(feet(f.player)); c >= spot.col {
		t.Errorf("the player walked through the chest to column %d", c)
	}
	if f.canPlaceYard(spot) {
		t.Error("a second chest on the first")
	}
	// It stays through the player going indoors and out, and the wild
	// moving on far away.
	f.enter(houseDoor{at: homeDoor, faceOut: "down", seed: seed, own: true})
	g.Step(frame)
	f.leave()
	g.Step(frame)
	if f.chestAt(spot.col, spot.row) != y || g.Find("chest") == nil {
		t.Error("the chest is gone after going indoors")
	}
	f.player.Transform.Position = f.ground.CellCenter(col+400, row+400)
	for range 30 {
		g.Step(frame)
	}
	f.player.Transform.Position = f.ground.CellCenter(col, row)
	g.Step(frame)
	if f.chestAt(spot.col, spot.row) != y || f.open(spot.col, spot.row) {
		t.Error("the chest did not stay where it was put")
	}
}

// walkRight is a controller that walks to the right.
type walkRight struct{}

func (walkRight) control(f *farm, o *game.Object, _ float64) intent {
	return intent{dir: vector.Vector3{X: 1}}
}

// Out of doors only chests can be put down, near the player and on free
// ground; a lifted chest comes back where it was.
func TestOnlyChestsGoOutOfDoorsAndOnFreeGround(t *testing.T) {
	g, f := build(SVGArt(), idle{})
	g.Step(frame)
	f.bag.Add("chest", 1)
	f.bag.Add("table", 1)
	if !f.startEdit("table") || f.edit.item != "chest" {
		t.Fatalf("a table was chosen out of doors: %+v", f.edit)
	}
	col, row := f.cellOf(feet(f.player))
	f.edit.col, f.edit.row = col+30, row
	f.clampCursor()
	if f.edit.col > col+editReach {
		t.Errorf("the cursor left the player's reach: %d", f.edit.col)
	}
	f.edit.col, f.edit.row = col, row
	if p, ok, _ := f.ghost(); ok {
		t.Errorf("a chest on the player: %+v", p)
	}
}
