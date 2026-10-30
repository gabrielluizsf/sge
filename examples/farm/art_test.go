package farm

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/sprite"
	"github.com/gabrielluizsf/sge/vector"
)

// recording is an Art that draws as another does and notes what it is
// asked for, its people facing eight ways.
type recording struct {
	Art
	asked map[string]bool
}

func (r *recording) Tile(name string) *canvas.Canvas {
	r.asked["tile "+name] = true
	return r.Art.Tile(name)
}

func (r *recording) Building(name string) *sprite.Sprite {
	r.asked["building "+name] = true
	return r.Art.Building(name)
}

func (r *recording) Person(who Who) Body {
	r.asked[fmt.Sprint("person ", who.Player)] = true
	return eightWay{r.Art.Person(who)}
}

// eightWay is a four-way body facing eight ways, the ways between drawn
// as down or up.
type eightWay struct{ Body }

func (eightWay) Facings() int { return 8 }

func (b eightWay) Clip(anim, dir string) Clip {
	for _, way := range []string{"down", "up"} {
		if strings.HasPrefix(dir, way) {
			dir = way
		}
	}
	return b.Body.Clip(anim, dir)
}

// A game started with an Art of its own is drawn by it: its tiles, its
// buildings and its people, who face as many ways as their bodies do.
func TestNewDrawsWithItsArt(t *testing.T) {
	r := &recording{Art: SVGArt(), asked: map[string]bool{}}
	s := New(r, idle{})
	if _, err := s.Play(Demo{Frames: 1, Width: 320, Height: 240}); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"tile grass0", "tile tree0", "building stall", "building house0", "person true"} {
		if !r.asked[want] {
			t.Errorf("the art was not asked for %s", want)
		}
	}
	o, p := s.f.player, s.f.hands
	p.facing = vector.Vector3{X: 1} // east in the world, down and right on screen
	if got := p.pose(o); got != "downright" {
		t.Fatalf("facing east, the farmer poses %q", got)
	}
	if !o.Animator.Play("downright_walk") {
		t.Error("no walk facing down and right")
	}
}

// Four facings are the side the way leans to most, eight the nearest of
// the ways between too.
func TestFacingOf(t *testing.T) {
	for _, c := range []struct {
		x, y float64
		n    int
		want string
	}{
		{1, 0.5, 4, "right"}, {1, 1, 4, "down"}, {-1, -1, 4, "up"}, {-1, 0.2, 4, "left"},
		{1, 0.5, 8, "downright"}, {1, 0.1, 8, "right"}, {-1, -1, 8, "upleft"}, {0, 1, 8, "down"}, {-1, 0.9, 8, "downleft"},
	} {
		if got := facingOf(c.x, c.y, c.n); got != c.want {
			t.Errorf("facingOf(%g, %g, %d) = %q, want %q", c.x, c.y, c.n, got, c.want)
		}
	}
}
