package terrain

import "testing"

func TestSmoothWearsAwayThreadsAndFillsNotches(t *testing.T) {
	shape := []string{
		"..........",
		".#####....",
		".#####.#..",
		".##.##.#..",
		".#####.#..",
		"..........",
	}
	in := func(c, r int) bool {
		return r >= 0 && r < len(shape) && c >= 0 && c < len(shape[r]) && shape[r][c] == '#'
	}
	s := Smooth(in, 5)
	if s(7, 3) {
		t.Error("a thread a cell wide was kept")
	}
	if !s(3, 3) {
		t.Error("a notch inside was not filled")
	}
	if !s(3, 2) || s(8, 0) {
		t.Error("the solid part or the outside changed")
	}
}
