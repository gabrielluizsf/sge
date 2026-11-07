package farm

import "testing"

// Pieces are put down from the bag on the cursor, one each, and the bag
// gives them up.
func TestPlacingAPieceTakesItFromTheBag(t *testing.T) {
	_, f := inHouse(t)
	f.bag.Take("chair", 2) // one left
	f.edit = nil
	if !f.startEdit("chair") || f.edit.item != "chair" {
		t.Fatalf("edit %+v", f.edit)
	}
	arrange(f, 5, 2)
	f.editAct() // take in hand
	f.editAct() // put down
	if !placed(f, "chair", 5, 2) || f.bag.Count("chair") != 0 {
		t.Fatalf("home %+v, chairs %d", f.home, f.bag.Count("chair"))
	}
	if f.room.things.At(5, 2) != pieceTileID(f.home[1], 0, 0) {
		t.Error("the room does not show it")
	}
}

// A piece does not go on another, outside the walls, on the player, or on
// the door's cell; the bag keeps it.
func TestPiecesGoOnlyOnFreeFloor(t *testing.T) {
	_, f := inHouse(t)
	f.startEdit("table")
	for _, at := range [][2]int{{1, 1}, {2, 1}, {0, 3}, {roomCols - 1, 2}, {entry[0], entry[1]}, {3, roomRows - 1}} {
		arrange(f, at[0], at[1])
		f.edit.hold = nil
		f.editAct()
		f.editAct()
		if len(f.home) != 1 || f.bag.Count("table") != 3 {
			t.Fatalf("a table on %v: home %+v", at, f.home)
		}
	}
	// The cursor stays on the floor even if told to leave it.
	f.editMove(-20, -20)
	if f.edit.col != 1 || f.edit.row != 1 {
		t.Errorf("cursor %d,%d", f.edit.col, f.edit.row)
	}
}

// The door's cell stays free and no row of furniture walls off a part of
// the floor: nobody is trapped.
func TestNothingBlocksTheDoorOrTrapsAnyone(t *testing.T) {
	_, f := inHouse(t)
	f.home = nil
	for c := 1; c <= roomCols-3; c++ { // a wall of chairs across row 3, all but the last cell
		f.home = append(f.home, newPiece("chair", c, 3, 0))
	}
	last := newPiece("chair", roomCols-2, 3, 0)
	if f.canStand(last) {
		t.Error("the last chair of a wall across the room was allowed")
	}
	if !f.canStand(newPiece("chair", 1, 4, 0)) {
		t.Error("a chair below the wall was refused")
	}
	if f.canStand(newPiece("chair", entry[0], entry[1], 0)) {
		t.Error("a chair on the door's cell was allowed")
	}
	// A bed in front of the door's cell would box the entry in.
	f.home = []piece{newPiece("chair", entry[0]-1, entry[1], 0), newPiece("chair", entry[0]+1, entry[1], 0), newPiece("chair", entry[0], entry[1]-1, 0)}
	if f.canStand(newPiece("chair", 1, 1, 0)) {
		t.Error("with the door boxed in, a piece was allowed")
	}
}

// Moving lifts a piece off the floor and puts it down elsewhere; turning
// it a quarter swaps the cells of a two-cell piece; Escape puts back what
// was lifted.
func TestMovingAndTurningAPiece(t *testing.T) {
	_, f := inHouse(t)
	arrange(f, 1, 1)
	f.editAct() // lift the bed
	if len(f.home) != 0 || f.edit.hold == nil || !f.edit.lifted {
		t.Fatalf("home %+v hold %+v", f.home, f.edit.hold)
	}
	f.dropHold()
	if len(f.home) != 1 || f.home[0].col != 1 || f.home[0].row != 1 {
		t.Fatalf("the bed was not put back: %+v", f.home)
	}
	f.editAct()
	f.editTurn()
	arrange(f, 5, 3)
	f.editAct()
	b := f.home[0]
	if b.col != 5 || b.row != 3 || b.heading != 1 {
		t.Fatalf("bed %+v", b)
	}
	if s := b.cells(f.art); s.cols != 1 || s.rows != 2 {
		t.Errorf("a bed turned a quarter covers %d by %d", s.cols, s.rows)
	}
	if f.pieceAt(5, 4) != 0 || f.room.things.At(5, 4) != pieceTileID(b, 0, 1) {
		t.Error("the turned bed does not cover its second cell")
	}
	// Turned in place, it still fits; where it does not fit, it stays.
	arrange(f, 5, 3)
	f.editTurn()
	if f.home[0].heading != 2 {
		t.Errorf("heading %d", f.home[0].heading)
	}
	f.home[0] = newPiece("bed", 6, 5, 0) // two cells, against the right wall
	f.relay()
	arrange(f, 6, 5)
	f.editTurn() // 1 by 2 down from row 5 would leave the floor
	if f.home[0].heading != 0 {
		t.Errorf("a bed turned where it does not fit: %+v", f.home[0])
	}
}

// Stored pieces go back in the bag, but never the last bed, and never a
// chest with something in it.
func TestStoringKeepsTheLastBedAndFullChests(t *testing.T) {
	_, f := inHouse(t)
	arrange(f, 1, 1)
	f.editStore()
	if len(f.home) != 1 || f.bag.Count("bed") != 3 {
		t.Fatalf("the last bed was stored: %+v, bag %d", f.home, f.bag.Count("bed"))
	}
	f.editAct() // lifted, it is still the last
	f.editStore()
	if f.edit.hold == nil || f.bag.Count("bed") != 3 {
		t.Error("a lifted last bed was stored")
	}
	f.dropHold()
	f.home = append(f.home, newPiece("bed", 1, 3, 0), newPiece("chest", 6, 1, 0))
	f.relay()
	arrange(f, 1, 3)
	f.editStore()
	if len(f.home) != 2 || f.bag.Count("bed") != 4 {
		t.Errorf("a second bed was not stored: %+v", f.home)
	}
	f.home[1].held[0] = chestSlot{item: wood, n: 2}
	arrange(f, 6, 1)
	f.editStore()
	if len(f.home) != 2 {
		t.Error("a chest with wood in it was stored")
	}
	f.home[1].held[0] = chestSlot{}
	f.editStore()
	if len(f.home) != 1 || f.bag.Count("chest") != 4 {
		t.Errorf("an empty chest was not stored: %+v", f.home)
	}
}

// Sleeping needs a bed: after any arranging there is still one to sleep in.
func TestSleepingIsAlwaysPossible(t *testing.T) {
	_, f := inHouse(t)
	if !hasBed(f.home) {
		t.Fatal("no bed")
	}
	arrange(f, 1, 1)
	f.editStore()
	f.stopEdit()
	if !hasBed(f.home) {
		t.Error("no bed after arranging")
	}
}

// Another family's house can not be arranged.
func TestOthersHousesCannotBeArranged(t *testing.T) {
	g, f := build(SVGArt(), idle{})
	d := houseDoor{at: homeDoor, faceOut: "down", seed: seed + 5, own: false}
	f.enter(d)
	g.Step(frame)
	f.bag.Add("chair", 1)
	if f.startEdit("chair") {
		t.Error("arranging in someone else's house")
	}
}
