// Command collect is a small sge game: move with the arrow keys or WASD
// and pick up the coins. It is the same program in Go and, transpiled by
// gots, on Deno.
//
// COLLECT_DEMO=n plays n frames by itself, at a fixed step and seed, and
// prints the score, so two runs can be compared.
package main

import (
	"fmt"
	"math/rand/v2"
	"os"
	"strconv"

	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/collider"
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/sprite"
	"github.com/gabrielluizsf/sge/vector"
)

const width, height = 480, 320

// player moves by the keys, or towards the nearest coin in a demo.
type player struct {
	speed float64
	demo  bool
	score int
}

func (p *player) Update(o *game.Object, dt float64) {
	var dir vector.Vector3
	if p.demo {
		if c := nearest(o); c != nil {
			dir = c.Transform.Position.Sub(o.Transform.Position)
		}
	} else if win := o.Game().Window(); win != nil {
		dir = keys(win)
	}
	if dir.SqrLength() > 0 {
		o.Transform.Translate(dir.Normalized().Scale(p.speed * dt))
	}
	pos := &o.Transform.Position
	pos.X = min(max(pos.X, 16), width-16)
	pos.Y = min(max(pos.Y, 16), height-16)
}

func keys(win *antui.Window) vector.Vector3 {
	var d vector.Vector3
	if win.KeyDown(antui.KeyLeft) || win.KeyDown(antui.KeyA) {
		d.X--
	}
	if win.KeyDown(antui.KeyRight) || win.KeyDown(antui.KeyD) {
		d.X++
	}
	if win.KeyDown(antui.KeyUp) || win.KeyDown(antui.KeyW) {
		d.Y--
	}
	if win.KeyDown(antui.KeyDown) || win.KeyDown(antui.KeyS) {
		d.Y++
	}
	return d
}

func nearest(o *game.Object) *game.Object {
	var best *game.Object
	for _, c := range o.Game().Objects() {
		if c.Name != "coin" {
			continue
		}
		if best == nil || c.Transform.Position.Distance(o.Transform.Position) < best.Transform.Position.Distance(o.Transform.Position) {
			best = c
		}
	}
	return best
}

// Collision picks up a coin, which turns up somewhere else.
func (p *player) Collision(o, other *game.Object) {
	if other.Name == "coin" {
		p.score++
		place(other)
	}
}

// Paint shows the score.
func (p *player) Paint(o *game.Object, cv *canvas.Canvas) {
	cv.Text(canvas.Point{X: 8, Y: 8}, fmt.Sprintf("score %d", p.score), canvas.RGB(240, 240, 240))
}

var rng = rand.New(rand.NewPCG(1, 2))

func place(coin *game.Object) {
	coin.Transform.Position = vector.Vector3{X: 20 + rng.Float64()*(width-40), Y: 30 + rng.Float64()*(height-50)}
}

// spin turns a coin around, for something to see.
type spin struct{}

func (spin) Update(o *game.Object, dt float64) { o.Transform.Rotate(vector.Vector3{Z: 1}, 2*dt) }

func main() {
	frames, _ := strconv.Atoi(os.Getenv("COLLECT_DEMO"))
	g := game.New()
	g.Background = canvas.RGB(20, 24, 32)

	body, err := sprite.Rect(24, 24, canvas.RGB(80, 170, 255))
	if err != nil {
		panic(err)
	}
	coinLook, err := sprite.Rect(12, 12, canvas.RGB(250, 200, 40))
	if err != nil {
		panic(err)
	}
	p := &player{speed: 220, demo: frames > 0}
	hero := g.Add(game.NewObject("player"))
	hero.Sprite, hero.Z = body, 1
	hero.Collider = collider.Box{X: -12, Y: -12, Width: 24, Height: 24}
	hero.Transform.Position = vector.Vector3{X: width / 2, Y: height / 2}
	hero.AddScript(p)
	for range 5 {
		c := g.Add(game.NewObject("coin"))
		c.Sprite, c.Collider = coinLook, collider.Circle{Radius: 7}
		c.AddScript(spin{})
		place(c)
	}
	if frames > 0 {
		g.FixedStep = 1.0 / 60
		n := 0
		g.Add(game.NewObject("clock")).AddScript(game.ScriptFunc(func(o *game.Object, dt float64) {
			if n++; n >= frames {
				o.Game().Quit()
			}
		}))
	}
	if err := g.Run("sge collect", width, height); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println("score:", p.score)
}
