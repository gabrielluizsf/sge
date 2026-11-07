package farm

import "github.com/gabrielluizsf/antui/canvas"

// TurnedArt is what an Art may also have whose furniture has a picture
// for each of the four ways it is turned, as the realistic look's models
// are baked at four headings: the piece name as Tile draws it, standing
// on its back cell, turned heading quarter turns. It is nil when that
// heading is not drawn. Without it the stylized look keeps the picture
// of heading 0 for 0 and 2, and mirrors it for 1 and 3.
type TurnedArt interface {
	Turned(name string, heading int) *canvas.Canvas
}

// pieceTile is the picture of cell i, j of the piece name turned heading,
// standing as a tile does, to be drawn at that cell.
func pieceTile(a Art, name string, heading, i, j int) *canvas.Canvas {
	full := pieceFull(a, name, heading)
	cols, rows := footprint(a, name, heading)
	if i >= cols || j >= rows {
		return clearTile()
	}
	if cols*rows == 1 {
		return full
	}
	return sliceOf(full, cols, rows, [2]int{i, j})
}

// pieceFull is the whole picture of the piece, drawn from its back cell,
// its bottom pieceDrop below that cell's bottom corner.
func pieceFull(a Art, name string, heading int) *canvas.Canvas {
	if t, ok := a.(TurnedArt); ok {
		if cv := t.Turned(name, heading); cv != nil {
			return cv
		}
	}
	cv := a.Tile(name)
	if cols, rows := footprint(a, name, 0); heading%2 == 1 && cols*rows == 1 {
		return mirrored(cv)
	}
	return cv
}

// clearTile is a tile with nothing drawn on it.
func clearTile() *canvas.Canvas {
	cv, err := canvas.NewCanvas(isoW, isoH)
	if err != nil {
		panic(err)
	}
	return cv
}

// mirrored is cv turned about its vertical middle.
func mirrored(cv *canvas.Canvas) *canvas.Canvas {
	out, err := canvas.NewCanvas(cv.Width, cv.Height)
	if err != nil {
		return cv
	}
	for y := range cv.Height {
		for x := range cv.Width {
			out.Pixels[y*out.Stride+cv.Width-1-x] = cv.At(x, y)
		}
	}
	return out
}
