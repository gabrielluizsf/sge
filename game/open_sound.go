package game

import "github.com/gabrielluizsf/antui"

// openSound starts playing the game's mixer, if it has one, at the rate
// the device runs at; nil where there is none, or no device to play on.
func (g *Game) openSound() *antui.Audio {
	if g.Sound == nil {
		return nil
	}
	a, err := antui.OpenAudio("sge", 44100, g.Sound.Fill)
	if err != nil {
		return nil
	}
	g.Sound.SetRate(a.Rate())
	return a
}
