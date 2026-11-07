package farm

import (
	"strings"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/examples/farm/art"
	"github.com/gabrielluizsf/sge/sprite"
)

// animalFrames is how many frames the stylized animals walk through.
const animalFrames = 4

// Animal is the animal as look says, drawn from SVG, facing the way its
// heading points on the screen: right or left, as the drawings are of the
// side.
func (a *svgArt) Animal(look AnimalLook) *sprite.Sprite {
	left := look.Heading >= 3 && look.Heading <= 5
	frame := 0
	if look.Motion == "walk" {
		frame = 1 + (look.Frame%animalFrames+animalFrames)%animalFrames
	}
	src, w, h := art.Animal(look.Kind, look.Coat, frame, look.Motion == "graze", false)
	if left {
		src = faceLeft(src, w)
	}
	return &sprite.Sprite{Image: animalPicture(src, w, h), PivotX: 0.5, PivotY: 1}
}

// faceLeft is the SVG src, w wide, turned to face the other way.
func faceLeft(src string, w int) string {
	src = strings.Replace(src, "</defs>\n", "</defs>\n<g transform=\"translate("+itoa(w)+" 0) scale(-1 1)\">\n", 1)
	return strings.Replace(src, "</svg>", "</g></svg>", 1)
}

// animalPicture is the SVG src drawn w by h units at the game's scale.
func animalPicture(src string, w, h int) *canvas.Canvas { return picture(src, w, h) }

// Pace is a slow walk, a frame in each of a step's four.
func (a *svgArt) Pace(kind string) (speed, fps float64) {
	return map[string]float64{"cow": 0.6, "sheep": 0.7, "horse": 1.2, "wolf": 1.8}[kind], 5
}
