package loader

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/game"
)

// Screen is the script of an object fixed to the screen that covers the
// game while a group loads: a backdrop, and a bar filling as the group's
// jobs are done, under Label's picture of the progress when it is set.
type Screen struct {
	Loader            *Loader
	Group             string
	Back, Track, Fill canvas.Color
	Label             func(done, asked int) *canvas.Canvas
}

func (s *Screen) Update(*game.Object, float64) {}

// Loading reports whether the group has jobs not done.
func (s *Screen) Loading() bool {
	done, asked := s.Loader.Progress(s.Group)
	return done < asked
}

func (s *Screen) Paint(_ *game.Object, cv *canvas.Canvas) {
	if !s.Loading() {
		return
	}
	done, asked := s.Loader.Progress(s.Group)
	cv.FillRect(canvas.Area{Width: cv.Width, Height: cv.Height}, s.Back)
	w, h := cv.Width/2, max(4, cv.Height/60)
	x, y := (cv.Width-w)/2, cv.Height*2/3
	cv.FillRect(canvas.Area{X: x, Y: y, Width: w, Height: h}, s.Track)
	cv.FillRect(canvas.Area{X: x, Y: y, Width: w * done / max(asked, 1), Height: h}, s.Fill)
	if s.Label != nil {
		if t := s.Label(done, asked); t != nil {
			cv.Blit((cv.Width-t.Width)/2, y-t.Height-h, t)
		}
	}
}
