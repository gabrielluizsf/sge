package terrain

import "testing"

func TestOpenLeavesNothingThinnerThanThreeCells(t *testing.T) {
	shape := []string{
		"............",
		".####.......",
		".####...#...",
		".####...##..",
		".####.......",
		"..#.........",
		"..#..##.....",
		".....##.....",
	}
	in := func(c, r int) bool {
		return r >= 0 && r < len(shape) && c >= 0 && c < len(shape[r]) && shape[r][c] == '#'
	}
	o := Open(in)
	for r := range shape {
		for c := range shape[r] {
			want := c >= 1 && c <= 4 && r >= 1 && r <= 4
			if o(c, r) != want {
				t.Errorf("cell %d,%d: %v", c, r, o(c, r))
			}
		}
	}
	// Opened noise: every cell in is in a whole three by three square.
	lakes := Open(Blob(5, 6, 0.55))
	for r := -30; r < 30; r++ {
		for c := -30; c < 30; c++ {
			if !lakes(c, r) {
				continue
			}
			ok := false
			for dy := -2; dy <= 0 && !ok; dy++ {
				for dx := -2; dx <= 0 && !ok; dx++ {
					all := true
					for y := r + dy; y < r+dy+3; y++ {
						for x := c + dx; x < c+dx+3; x++ {
							all = all && lakes(x, y)
						}
					}
					ok = all
				}
			}
			if !ok {
				t.Fatalf("cell %d,%d of the lakes is in no whole 3x3 square", c, r)
			}
		}
	}
}

func TestMemoAnswersAsItsAreaAndForgets(t *testing.T) {
	calls := 0
	in := func(c, r int) bool { calls++; return (c+r)%2 == 0 }
	m := Memo(in, 4)
	for range 3 {
		if !m(2, 2) || m(1, 2) {
			t.Fatal("wrong answer")
		}
	}
	if calls != 2 {
		t.Errorf("%d calls for two cells", calls)
	}
	for i := range 10 {
		m(i, 9)
	}
	if calls > 2+10 {
		t.Errorf("%d calls", calls)
	}
}

func TestCellKeysDifferAndStopAtTheirReach(t *testing.T) {
	seen := map[int][2]int{}
	for r := -3; r <= 3; r++ {
		for c := -3; c <= 3; c++ {
			for _, p := range [][2]int{{c, r}, {c + reach - 4, r}, {c, r - reach + 4}} {
				k, ok := cellKey(p[0], p[1])
				if !ok {
					t.Fatalf("%v out of reach", p)
				}
				if q, dup := seen[k]; dup {
					t.Fatalf("%v and %v share a key", p, q)
				}
				seen[k] = p
			}
		}
	}
	if _, ok := cellKey(reach, 0); ok {
		t.Error("a cell out of reach was keyed")
	}
}
