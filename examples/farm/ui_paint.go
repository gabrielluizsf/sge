package farm

import "github.com/gabrielluizsf/antui/canvas"

// paint draws the menu that is open: the pause menu, or the credits over
// it, or a licence over those.
func (u *ui) paint(cv *canvas.Canvas) {
	switch u.open {
	case pausing:
		u.pause.Paint(cv, u.f.art)
	case commanding:
		u.pause.Paint(cv, u.f.art)
		u.reader.Paint(cv, u.f.art)
	case crediting:
		u.credits.Paint(cv, u.f.art)
	case licensing:
		u.credits.Paint(cv, u.f.art)
		u.reader.Paint(cv, u.f.art)
	}
}
