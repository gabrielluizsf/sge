package farm

import (
	"fmt"

	"github.com/gabrielluizsf/antui/canvas"
)

// Demo is a run of the game by itself, without a window or sound: Frames
// frames at a fixed step of a sixtieth of a second, seen in a view Width
// by Height pixels.
type Demo struct {
	Frames        int
	Width, Height int
}

// Play plays d on s and is the view it was drawn on, as Shot draws it.
func (s *Session) Play(d Demo) (*canvas.Canvas, error) {
	cv, err := canvas.NewCanvas(d.Width, d.Height)
	if err != nil {
		return nil, err
	}
	s.Game.Draw(cv) // the view's size, as the window gives it to Run
	for range d.Frames {
		s.Game.Step(1.0 / 60)
	}
	return cv, nil
}

// Report is how the farm is doing: its money, what it has sold, its
// seeds and how many plots have a crop on them.
func (s *Session) Report() string {
	f := s.f
	return fmt.Sprintf("money: %d, sold: %d, seeds: [%d %d], growing: %d", f.money, f.sold, f.bag.Count(seedsOf(0)), f.bag.Count(seedsOf(1)), growing(f))
}

// Shot draws the game on cv and writes it to the PNG file path.
func (s *Session) Shot(cv *canvas.Canvas, path string) error {
	s.Game.Draw(cv)
	return writePNG(cv, path)
}

// growing is how many plots have a crop on them.
func growing(f *farm) int {
	n := 0
	for _, p := range f.plots {
		if p.crop >= 0 {
			n++
		}
	}
	return n
}
