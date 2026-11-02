package farm

import (
	"strings"
	"testing"

	"github.com/gabrielluizsf/sge/examples/farm/art"
)

// A bed two cells long is laid on two free cells, drawn from the front
// one, the other a block, in every room; one cell each without a span.
func TestFurnitureCoversTheCellsItSpans(t *testing.T) {
	long := func(p byte) (int, int) {
		if p == 'b' || p == 'B' {
			return 2, 1
		}
		return 1, 1
	}
	for seed := range uint32(200) {
		_, things := roomLayout(seed, long)
		rows := strings.Split(strings.TrimRight(things, "\n"), "\n")
		beds := 0
		for r, row := range rows {
			for c := range row {
				if row[c] != 'b' && row[c] != 'B' {
					continue
				}
				beds++
				if c == 0 || rows[r][c-1] != 'x' {
					t.Fatalf("seed %d: the bed at %d,%d has no block behind it:\n%s", seed, c, r, things)
				}
			}
		}
		if beds != 1 {
			t.Fatalf("seed %d: %d beds:\n%s", seed, beds, things)
		}
	}
	one := func(byte) (int, int) { return 1, 1 }
	_, a := roomLayout(7, one)
	if strings.Count(a, "x") != roomRows+roomCols-1-1 { // the open sides, the mat's cell not
		t.Errorf("blocks beyond the open sides with no span:\n%s", a)
	}
}

// fit finds no room where none is free, and leaves t as it was.
func TestFitLeavesAFullRoomAlone(t *testing.T) {
	grid := [][]byte{[]byte(" x"), []byte("xx")}
	if fit(grid, [2]int{0, 0}, 2, 1, 'b') {
		t.Fatal("a two-cell bed fit in one free cell")
	}
	if string(grid[0]) != " x" {
		t.Errorf("the room changed: %q", grid[0])
	}
}

// Every kind of house the realistic look has has its walls inside.
func TestEveryHouseKindHasWalls(t *testing.T) {
	for k := range art.WallKinds {
		for _, w := range []string{"wall_n", "wall_nwin", "wall_w", "wall_nw"} {
			found := false
			for _, n := range TileNames() {
				found = found || n == w+itoa(k)
			}
			if !found {
				t.Errorf("no tile %s%d", w, k)
			}
		}
	}
}
