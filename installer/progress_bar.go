package installer

import (
	"errors"
	"sync/atomic"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/game"
)

// errClosed is the error of a window closed while installing.
var errClosed = errors.New("installer: the window was closed before the installing was done")

// progressBar is the script painting the window: the backdrop, the bar
// as far as the job under way, at, and its label.
type progressBar struct {
	look     Window
	jobs     []job
	at       *atomic.Int32
	shown    int32
	finished bool
}

func (p *progressBar) Update(o *game.Object, _ float64) {
	if at := p.at.Load(); at != p.shown {
		p.shown = at
		if g := o.Game(); g != nil {
			g.Invalidate()
		}
	}
}

func (p *progressBar) Paint(_ *game.Object, cv *canvas.Canvas) {
	done, asked := int(p.at.Load()), len(p.jobs)
	cv.FillRect(canvas.Area{Width: cv.Width, Height: cv.Height}, p.look.Back)
	w, h := cv.Width*2/3, max(6, cv.Height/40)
	x, y := (cv.Width-w)/2, cv.Height*3/5
	cv.FillRect(canvas.Area{X: x, Y: y, Width: w, Height: h}, p.look.Track)
	cv.FillRect(canvas.Area{X: x, Y: y, Width: w * done / max(asked, 1), Height: h}, p.look.Fill)
	if p.look.Label == nil || done >= asked {
		return
	}
	if t := p.look.Label(p.jobs[done].name, done, asked); t != nil {
		cv.Blit((cv.Width-t.Width)/2, y-t.Height-h, t)
	}
}
