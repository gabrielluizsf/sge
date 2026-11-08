package pathfind

import "testing"

// A search asks about a cell once, however many neighbours it has: what
// it costs to tell whether a cell is open is paid once.
func TestASearchAsksAboutACellOnce(t *testing.T) {
	asked := map[Cell]int{}
	g := Grid{Limit: 2500, Open: func(c, r int) bool {
		asked[Cell{c, r}]++
		return c%7 != 3 || r%5 == 0
	}}
	if way := g.Path(Cell{0, 0}, Cell{20, 20}); way == nil {
		t.Fatal("no way")
	}
	for c, n := range asked {
		if n > 1 {
			t.Fatalf("cell %v asked about %d times", c, n)
		}
	}
}

func BenchmarkPathAcrossARoughField(b *testing.B) {
	g := Grid{Limit: 2500, Open: func(c, r int) bool { return c%7 != 3 || r%5 == 0 }}
	for b.Loop() {
		g.Path(Cell{0, 0}, Cell{40, 40})
	}
}
