package collider

import "testing"

func TestCollides(t *testing.T) {
	cases := []struct {
		name string
		a, b Collider
		want bool
	}{
		{"overlapping boxes", Box{0, 0, 10, 10}, Box{5, 5, 10, 10}, true},
		{"distant boxes", Box{0, 0, 10, 10}, Box{20, 20, 5, 5}, false},
		{"boxes touching an edge", Box{0, 0, 10, 10}, Box{10, 0, 10, 10}, false},
		{"circle in a box", Circle{5, 5, 1}, Box{0, 0, 10, 10}, true},
		{"box around a circle", Box{0, 0, 10, 10}, Circle{5, 5, 1}, true},
		{"circle beside a box corner", Circle{12, 12, 2}, Box{0, 0, 10, 10}, false},
		{"overlapping circles", Circle{0, 0, 5}, Circle{8, 0, 5}, true},
		{"distant circles", Circle{0, 0, 5}, Circle{20, 0, 5}, false},
		{"circles whose bounds overlap", Circle{0, 0, 5}, Circle{8, 8, 5}, false},
	}
	for _, c := range cases {
		if got := Collides(c.a, c.b); got != c.want {
			t.Errorf("%s: Collides = %v, want %v", c.name, got, c.want)
		}
	}
}

func TestCircleBounds(t *testing.T) {
	if got, want := (Circle{5, 6, 2}).Bounds(), (Box{3, 4, 4, 4}); got != want {
		t.Errorf("Bounds = %v, want %v", got, want)
	}
}
