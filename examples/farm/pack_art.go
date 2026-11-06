package farm

import (
	"strings"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/backpack"
)

// pieceIcon starts the icon name of a piece of furniture: the art's picture
// of its tile, made small.
const pieceIcon = "icon_piece_"

// iconPx is the side of an icon, in pixels of art.
const iconPx = 16

// packArt is the art the inventory is drawn with: f's, with the icons of
// the furniture made from its tiles, whichever look it is, kept once made.
type packArt struct {
	Art
	icons map[string]*canvas.Canvas
}

var _ backpack.Art = packArt{}

// Icon is the art's icon, or, for a piece of furniture, its tile shrunk.
func (p packArt) Icon(name string) *canvas.Canvas {
	piece, ok := strings.CutPrefix(name, pieceIcon)
	if !ok {
		return p.Art.Icon(name)
	}
	if cv, ok := p.icons[piece]; ok {
		return cv
	}
	cv := shrunkIcon(pieceFull(p.Art, piece, 0))
	if cv != nil {
		p.icons[piece] = cv
	}
	return cv
}

// shrunkIcon is cv's picture, what shows of it, fitted to an icon in the
// middle of its square.
func shrunkIcon(cv *canvas.Canvas) *canvas.Canvas {
	x0, y0, x1, y1 := cv.Width, cv.Height, 0, 0
	for y := range cv.Height {
		for x := range cv.Width {
			if cv.At(x, y).A() != 0 {
				x0, y0, x1, y1 = min(x0, x), min(y0, y), max(x1, x+1), max(y1, y+1)
			}
		}
	}
	side := iconPx * scale
	out, err := canvas.NewCanvas(side, side)
	if err != nil || x1 <= x0 {
		return out
	}
	w, h := x1-x0, y1-y0
	k := max(w, h)
	dw, dh := w*(side-2*scale)/k, h*(side-2*scale)/k
	out.BlitScaled(canvas.Area{X: (side - dw) / 2, Y: (side - dh) / 2, Width: max(dw, 1), Height: max(dh, 1)}, cropped(cv, x0, y0, w, h))
	return out
}

// cropped is the w by h part of cv from x, y.
func cropped(cv *canvas.Canvas, x, y, w, h int) *canvas.Canvas {
	out, err := canvas.NewCanvas(w, h)
	if err != nil {
		return cv
	}
	for j := range h {
		for i := range w {
			out.Pixels[j*out.Stride+i] = cv.At(x+i, y+j)
		}
	}
	return out
}

// packArt is the art the inventory is drawn with.
func (f *farm) packArt() packArt {
	if f.icons == nil {
		f.icons = map[string]*canvas.Canvas{}
	}
	return packArt{f.art, f.icons}
}
