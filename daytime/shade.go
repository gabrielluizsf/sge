package daytime

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/game"
)

// Shade is the script of an object fixed to the screen, over the world,
// that darkens it as the clock's light goes: by night the screen is
// covered in Night, whose alpha is how dark the night is, and at dusk and
// dawn in Glow, warm, half as strong.
type Shade struct {
	Clock *Clock
	Night canvas.Color
	Glow  canvas.Color
}

func (s Shade) Update(*game.Object, float64) {}

func (s Shade) Paint(_ *game.Object, cv *canvas.Canvas) {
	if c := s.Tint(); c.A() > 0 {
		cv.FillRect(canvas.Area{Width: cv.Width, Height: cv.Height}, c)
	}
}

// Tint is the colour the screen is covered in at the clock's hour.
func (s Shade) Tint() canvas.Color {
	dark := 1 - s.Clock.Light()
	if dark == 0 {
		return canvas.Transparent
	}
	// The glow is strongest half way through dusk or dawn.
	glow := float32(1 - (2*dark-1)*(2*dark-1))
	c := canvas.Mix(canvas.Fade(s.Night, 255), canvas.Fade(s.Glow, 255), glow*0.5)
	a := float64(s.Night.A())*dark + float64(s.Glow.A())*float64(glow)*0.5
	return canvas.Fade(c, int(min(255, a)))
}
