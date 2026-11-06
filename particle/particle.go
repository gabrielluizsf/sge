// Package particle is a pool of small specks that fly, fall and come to
// rest on the ground: sparks, dust, blood, whatever the game configures it
// to be. It knows nothing of what its specks are. Everything it draws by
// chance comes from package random, and nothing it does allocates after
// it is made, so the same seed plays the same in Go and on Deno, frame by
// frame, pixel by pixel.
package particle

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/random"
)

// Config is what a kind of speck is: how fast it leaves, how hard it is
// pulled down, how long it lies on the ground, how big it is and which
// colours it comes in. Speeds are in pixels a second, sizes in pixels, a
// speck's height Z in pixels over the ground.
type Config struct {
	// Speed is the spread over the ground, up to this much a second in any
	// way (plus what the emitter's direction adds); Lift the upward speed,
	// from half this to all of it.
	Speed, Lift float64
	// Gravity pulls a speck down, pixels a second each second.
	Gravity float64
	// Rest is how many seconds a speck lies where it landed before it has
	// faded away, and Spread how many seconds more or less, up to, each
	// varies by.
	Rest, Spread float64
	// Size is a speck's side in pixels, from Size to Size+Grow.
	Size, Grow int
	// Colors are the colours a speck is, one drawn each.
	Colors []canvas.Color
}

// speck is one particle: where it is over the ground, how fast it moves,
// and how long it has lain.
type speck struct {
	x, y, z    float64
	vx, vy, vz float64
	rest, lain float64
	size       int
	color      canvas.Color
	live       bool
}

// System is a pool of specks of one kind. It holds as many as it was made
// for; when it is full the oldest is taken for the next.
type System struct {
	cfg    Config
	rng    *random.Source
	specks []speck
	next   int
	live   int
}

// New is a system of capacity specks of kind cfg, drawing by chance from
// seed.
func New(capacity int, cfg Config, seed uint32) *System {
	return &System{cfg: cfg, rng: random.New(seed), specks: make([]speck, capacity)}
}

// Live is how many specks are showing: in the air or lying on the ground.
func (s *System) Live() int { return s.live }
