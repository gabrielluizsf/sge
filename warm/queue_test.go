package warm

import (
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/sprite"
)

func lazy(made *int) *sprite.Sprite {
	return &sprite.Sprite{Make: func() *canvas.Canvas {
		*made++
		cv, _ := canvas.NewCanvas(1, 1)
		return cv
	}}
}

// Sprites are made PerFrame a frame, in the order added, and the sprite is
// made once, whoever draws it after.
func TestSpritesAreMadeAFewAFrame(t *testing.T) {
	made := 0
	a, b, c := lazy(&made), lazy(&made), lazy(&made)
	q := &Queue{PerFrame: 2}
	q.Add(a, b, c)
	q.Update(nil, 0)
	if made != 2 || !a.IsReady() || !b.IsReady() || c.IsReady() || q.Waiting() != 1 {
		t.Fatalf("after one frame: %d made, %d waiting", made, q.Waiting())
	}
	q.Update(nil, 0)
	q.Update(nil, 0)
	if made != 3 || q.Waiting() != 0 {
		t.Fatalf("after three: %d made, %d waiting", made, q.Waiting())
	}
	c.Picture()
	if made != 3 {
		t.Fatal("made again when drawn")
	}
}

// A sprite made already is not counted against the frame's share, and
// none set means one.
func TestAMadeSpriteCostsNothing(t *testing.T) {
	made := 0
	a, b := lazy(&made), lazy(&made)
	q := &Queue{}
	q.Add(a, b)
	a.Ready()
	q.Update(nil, 0)
	if made != 2 || q.Waiting() != 0 {
		t.Fatalf("%d made, %d waiting", made, q.Waiting())
	}
}

// A long queue is made twice as fast, until it is short again.
func TestALongQueueIsMadeFaster(t *testing.T) {
	made := 0
	q := &Queue{PerFrame: 1, Backlog: 3}
	for range 6 {
		q.Add(lazy(&made))
	}
	q.Update(nil, 0) // 6 waiting, more than 3: two made
	if made != 2 {
		t.Fatalf("%d made with a long queue", made)
	}
	q.Update(nil, 0) // 4 waiting: two again
	q.Update(nil, 0) // 2 waiting: one
	if made != 5 {
		t.Fatalf("%d made after three frames", made)
	}
}
