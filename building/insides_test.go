package building

import (
	"fmt"
	"testing"
)

// Each piece of a house's inside is a sound model on its one cell: the
// floors level under the ground's, the walls a storey tall along the
// cell's back edges, a pane in the window's.
func TestInsidesStandOnTheirCell(t *testing.T) {
	names := map[string]bool{}
	for _, k := range Insides() {
		names[k.Name] = true
		for _, palette := range []Palette{Plain(), textured()} {
			m := k.Model(palette)
			for _, p := range m.Meshes[0].Primitives {
				sound(t, soundOf{k.Name, m, p})
			}
		}
		lo, hi := bounds(k.Model(Plain()), "")
		if lo[0] < -0.5-1e-9 || lo[2] < -0.5-1e-9 || hi[0] > 0.5+1e-9 || hi[2] > 0.5+1e-9 {
			t.Errorf("%s reaches off its cell: %v to %v", k.Name, lo, hi)
		}
		if lo[1] < -floorThick-1e-9 {
			t.Errorf("%s sinks to %v", k.Name, lo[1])
		}
	}
	for _, n := range []string{"floor", "mat", "rug0", "rug1", "rug2", "rug3"} {
		if !names[n] {
			t.Errorf("no %s", n)
		}
	}
	for k, h := range houseStyles {
		for _, w := range []string{"wall_n", "wall_nwin", "wall_w", "wall_nw"} {
			if !names[fmt.Sprint(w, k)] {
				t.Errorf("no %s%d", w, k)
			}
		}
		walls := innerWalls(k, h)
		_, hi := bounds(walls[0].Model(Plain()), "")
		if want := max(h.Storey, 2.7); hi[1] < want-1e-9 || hi[1] > 3+1e-9 {
			t.Errorf("house %d's walls %v tall", k, hi[1])
		}
		if lo, _ := bounds(walls[1].Model(Plain()), "glass"); lo[1] < paneLo-1e-6 || lo[1] > paneHi {
			t.Errorf("house %d's window has no pane", k)
		}
	}
}
