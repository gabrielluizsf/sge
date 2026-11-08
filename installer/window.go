package installer

import (
	"sync/atomic"

	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/game"
)

// Window is how an installer's window looks: its title and size, a
// backdrop, and a bar filling as the jobs are done, under Label's picture
// of the step under way (nil shows the bar alone), as a game's loading
// screen does.
type Window struct {
	Title             string
	Width, Height     int
	Back, Track, Fill canvas.Color
	Label             func(step string, done, asked int) *canvas.Canvas
}

// show runs jobs while the window win shows how far they have come: they
// run beside the window, so it is drawn and answers all along. It returns
// the jobs' error, or one for a window closed before they were done.
func show(win *antui.Window, look Window, s Setup, jobs []job) error {
	var at atomic.Int32
	result := make(chan error, 1)
	go func() { result <- runJobs(s, jobs, func(i int) { at.Store(int32(i)) }) }()
	g := game.New()
	bar := &progressBar{look: look, jobs: jobs, at: &at, shown: -1}
	o := game.NewObject("installer")
	o.Screen = true
	o.AddScript(bar)
	var err error
	o.AddScript(game.ScriptFunc(func(*game.Object, float64) {
		select {
		case err = <-result:
			bar.finished = true
			g.Quit()
		default:
		}
	}))
	g.Add(o)
	g.Show(win)
	if !bar.finished {
		return errClosed
	}
	return err
}
