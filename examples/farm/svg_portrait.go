package farm

import (
	"fmt"

	"github.com/gabrielluizsf/antui/canvas"
)

// whoKey tells people apart as their portraits go: all but their age in
// years.
func whoKey(who Who) string {
	who.Years = 0
	return fmt.Sprint(who)
}

// bustShare is how much of a person's picture, from the top of their head,
// a portrait shows: the head and shoulders.
const bustShare = 0.46

// bustZoom is how many times a portrait is enlarged: pixel art, by whole
// pixels.
const bustZoom = 5

// Portrait is who's head and shoulders, from their standing picture,
// enlarged.
func (a *svgArt) Portrait(who Who) *canvas.Canvas {
	key := "portrait" + whoKey(who)
	if cv, ok := a.hud[key]; ok {
		return cv
	}
	cv := Bust(a.Person(who).Clip(AnimStand, "down").Frames[0].Make())
	a.hud[key] = cv
	return cv
}

// Bust is the top of the figure in body, its head and shoulders, cut to
// the figure's edges and enlarged by whole pixels.
func Bust(body *canvas.Canvas) *canvas.Canvas {
	x0, y0, x1, y1 := body.Width, body.Height, 0, 0
	for y := range body.Height {
		for x := range body.Width {
			if body.At(x, y).A() > 0 {
				x0, y0, x1, y1 = min(x0, x), min(y0, y), max(x1, x+1), max(y1, y+1)
			}
		}
	}
	if x1 <= x0 {
		return body
	}
	h := max(1, int(float64(y1-y0)*bustShare))
	out, err := canvas.NewCanvas((x1-x0)*bustZoom, h*bustZoom)
	if err != nil {
		panic(err)
	}
	for y := range out.Height {
		for x := range out.Width {
			out.Put(x, y, body.At(x0+x/bustZoom, y0+y/bustZoom))
		}
	}
	return out
}
