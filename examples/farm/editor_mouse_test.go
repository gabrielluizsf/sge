package farm

import (
	"testing"

	"github.com/gabrielluizsf/sge/pathfind"
)

// Pointing moves the cursor to the cell pointed at, no further than a piece
// may be put, and a click then does what Enter does there.
func TestThePointerMovesTheCursor(t *testing.T) {
	_, f := newGame(idle{})
	f.bag.Add("chest", 1)
	r, ok := f.freePenSpot()
	if !ok {
		t.Fatal("no free ground")
	}
	at := pathfind.Cell{Col: r.col + 2, Row: r.row + 2}
	f.player.Transform.Position = f.center(at)
	if !f.startEdit("chest") {
		t.Fatal("no arranging")
	}
	f.editPoint(at.Col+2, at.Row+1)
	if f.edit.col != at.Col+2 || f.edit.row != at.Row+1 {
		t.Fatalf("cursor at %d, %d", f.edit.col, f.edit.row)
	}
	f.editPoint(at.Col+200, at.Row)
	if f.edit.col != at.Col+editReach {
		t.Errorf("the cursor went to %d", f.edit.col)
	}
	f.editPoint(at.Col+2, at.Row+1)
	f.editAct()
	f.editAct()
	if f.chestAt(at.Col+2, at.Row+1) == nil || f.bag.Count("chest") != 0 {
		t.Errorf("no chest placed by the pointer: %d in the bag", f.bag.Count("chest"))
	}
}
