package town

import "testing"

// The way from a town leads to the town its road leads to, and back the
// same way reversed, turning only at right angles.
func TestARouteGoesThereAndBackTheSameWay(t *testing.T) {
	l := roomy()
	found := 0
	for ry := range 4 {
		for rx := range 4 {
			a := l.Of(rx, ry)
			if a == nil {
				continue
			}
			b, there := l.Route(a)
			if b == nil {
				continue
			}
			found++
			back, home := l.Route(b)
			if back != nil && back == a {
				if len(home) != len(there) {
					t.Fatalf("%d points there, %d back", len(there), len(home))
				}
				for i := range there {
					if there[i] != home[len(home)-1-i] {
						t.Fatalf("point %d differs", i)
					}
				}
			}
			if there[0] != a.Center || there[len(there)-1] != b.Center {
				t.Errorf("route %v does not run plaza to plaza", there)
			}
			for i := 1; i < len(there); i++ {
				if there[i].Col != there[i-1].Col && there[i].Row != there[i-1].Row {
					t.Errorf("leg %v to %v is not straight", there[i-1], there[i])
				}
			}
		}
	}
	if found == 0 {
		t.Fatal("no town has a road")
	}
}
