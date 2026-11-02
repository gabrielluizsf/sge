package farm

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/game"
)

// dim is the shade laid over the world behind the main menu, to read the
// menu on it.
var dim = canvas.RGBA(canvas.RGBAParams{R: 10, G: 12, B: 24, A: 110})

// Paint draws the screen that is open over the world, shaded.
func (m *mainMenu) Paint(_ *game.Object, cv *canvas.Canvas) {
	if m.open == creating {
		return
	}
	cv.FillRect(canvas.Area{Width: cv.Width, Height: cv.Height}, dim)
	switch m.open {
	case titling:
		m.main.Paint(cv, m.f.art)
	case confirming:
		m.confirm.Paint(cv, m.f.art)
	case crediting:
		m.credits.Paint(cv, m.f.art)
	case licensing:
		m.credits.Paint(cv, m.f.art)
		m.reader.Paint(cv, m.f.art)
	}
}
