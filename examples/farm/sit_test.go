package farm

import (
	"testing"

	"github.com/gabrielluizsf/sge/vector"
)

// E beside a chair sits the player on it, facing the way it is turned; time
// goes on while seated; any step stands them up where they were.
func TestSittingOnAChairAndStandingUp(t *testing.T) {
	g, f := inHouse(t)
	f.home = append(f.home, newPiece("chair", 4, 3, 2))
	f.relay()
	f.player.Transform.Position = f.ground.CellCenter(4, 4)
	stood := f.player.Transform.Position
	words, do, ok := f.interaction(f.player)
	if !ok || words != "E: sit" {
		t.Fatalf("%q %v", words, ok)
	}
	do()
	if f.seat == nil || f.player.Transform.Position != f.ground.CellCenter(4, 3) {
		t.Fatalf("seat %+v at %v", f.seat, f.player.Transform.Position)
	}
	if f.hands.facing != (vector.Vector3{Y: -1}) { // heading 2 faces up
		t.Errorf("facing %v", f.hands.facing)
	}
	hour := f.clock.Hour
	f.hands.ctl = idle{}
	for range 120 {
		g.Step(frame)
	}
	if f.clock.Hour == hour {
		t.Error("time stopped while seated")
	}
	if f.player.Transform.Position != f.ground.CellCenter(4, 3) {
		t.Error("the seated player moved")
	}
	if words, _, _ := f.interaction(f.player); words != "Move to stand up" {
		t.Errorf("hint %q", words)
	}
	f.hands.ctl = walkRight{}
	for range 60 {
		g.Step(frame)
	}
	if f.seat == nil || f.player.Transform.Position != f.ground.CellCenter(4, 3) {
		t.Error("the seated player was pushed")
	}
	f.seatStep(false)
	if f.seat == nil {
		t.Error("stood up without a step")
	}
	f.seatStep(true)
	if f.seat != nil || f.player.Transform.Position != stood {
		t.Errorf("stood at %v, not %v", f.player.Transform.Position, stood)
	}
}

// A seated player is not stood up by anything but their own step, and
// can not arrange furniture or sleep from the chair.
func TestTheSeatedPlayerStaysPut(t *testing.T) {
	_, f := inHouse(t)
	f.home = append(f.home, newPiece("chair", 4, 3, 0))
	f.relay()
	f.player.Transform.Position = f.ground.CellCenter(4, 4)
	f.sitDown(1)
	if f.startEdit("chair") {
		t.Error("arranging while seated")
	}
	f.standUp()
	if f.seat != nil || f.player.Transform.Position != f.ground.CellCenter(4, 4) {
		t.Errorf("stood at %v", f.player.Transform.Position)
	}
}

// bakedSeatArt is an art that draws the chair with the person on it.
type bakedSeatArt struct{ longBed }

func (bakedSeatArt) SeatBaked() bool { return true }

// In every one of its four rotations a chair turns the player to look the
// way its front does, on a body of eight facings, as the realistic look's.
func TestTheSeatedPlayerFacesTheChairsFrontInEveryRotation(t *testing.T) {
	want := []string{"downleft", "upleft", "upright", "downright"}
	for h, dir := range want {
		g, f := inHouse(t)
		f.home = append(f.home, newPiece("chair", 4, 3, h))
		f.relay()
		g.Step(frame)
		f.hands.facings = 8
		f.sitDown(len(f.home) - 1)
		if got := f.hands.pose(f.player); got != dir {
			t.Errorf("heading %d faces %s, want %s", h, got, dir)
		}
	}
}

// With a look that draws the chair with the sitter, its own tile goes while
// they sit, in every rotation, and comes back solid when they stand.
func TestABakedSeatTakesTheChairsTileWhileSitting(t *testing.T) {
	for h := range headings {
		_, f := inHouse(t)
		f.art = bakedSeatArt{}
		f.home = append(f.home, newPiece("chair", 4, 3, h))
		f.relay()
		tile := f.room.things.At(4, 3)
		if tile < 0 {
			t.Fatal("no chair on the floor")
		}
		f.sitDown(len(f.home) - 1)
		if got := f.room.things.At(4, 3); got >= 0 {
			t.Errorf("heading %d: tile %d stays under the sitter", h, got)
		}
		f.standUp()
		if got := f.room.things.At(4, 3); got != tile {
			t.Errorf("heading %d: tile %d after standing, want %d", h, got, tile)
		}
	}
}

// A look that draws no chair with the sitter keeps its tile.
func TestAPlainSeatKeepsTheChairsTile(t *testing.T) {
	_, f := inHouse(t)
	f.home = append(f.home, newPiece("chair", 4, 3, 0))
	f.relay()
	tile := f.room.things.At(4, 3)
	f.sitDown(len(f.home) - 1)
	if got := f.room.things.At(4, 3); got != tile {
		t.Errorf("tile %d, want %d", got, tile)
	}
}
