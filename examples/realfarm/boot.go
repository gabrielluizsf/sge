package main

import (
	"fmt"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/examples/realfarm/kit/bakes"
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/loader"
)

// bootGroup is the loader group of what the game cannot start without.
const bootGroup = "boot"

// boot bakes what the farm cannot start without, tasks, in a window of
// title that shows how far it has come, and returns when it is all baked
// or the window is closed (an error then).
func boot(b *bakes.Baker, tasks []bakes.Task) error {
	var l loader.Loader
	l.Budget = 1 // one task a frame: each is long, the window is drawn between
	failed := make(chan error, 1)
	current := ""
	for _, t := range tasks {
		l.Request(loader.Job{Name: t.Name, Group: bootGroup, Make: func() (any, error) {
			current = t.Name
			err := b.Do(t)
			if err != nil {
				select {
				case failed <- err:
				default:
				}
			}
			return nil, err
		}}, 0)
	}
	g := game.New()
	screen := &loader.Screen{Loader: &l, Group: bootGroup,
		Back: canvas.RGB(20, 24, 22), Track: canvas.RGB(52, 60, 54), Fill: canvas.RGB(150, 200, 110),
		Label: func(done, asked int) *canvas.Canvas {
			return bootLabel(done, asked, current)
		}}
	o := game.NewObject("boot")
	o.Screen = true
	o.AddScript(&l)
	o.AddScript(screen)
	o.AddScript(game.ScriptFunc(func(_ *game.Object, _ float64) {
		if !screen.Loading() || len(failed) > 0 {
			g.Quit()
		}
	}))
	g.Add(o)
	if err := g.Run("sge realfarm", 960, 540); err != nil {
		return err
	}
	select {
	case err := <-failed:
		return err
	default:
	}
	if screen.Loading() {
		return fmt.Errorf("closed while the world was being prepared")
	}
	return b.Shelf.Save()
}

// bootLabel is the loading screen's line: how many of the tasks asked are
// done and which one is baking, in the HUD's own serif, never the
// stylized farm's face.
func bootLabel(done, asked int, current string) *canvas.Canvas {
	return serifText(fmt.Sprintf("Preparando o mundo %d/%d - %s", done, asked, current), 'w')
}
