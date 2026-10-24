package tilemap

import (
	"fmt"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/game"
	"github.com/gabrielluizsf/sge/sprite"
)

// New is an empty map of cols by rows cells of tiles, an object named
// name to add to a game.
func New(name string, tiles *TileSet, cols, rows int) (*Map, error) {
	w, h := cols*tiles.Width, rows*tiles.Height
	if cols <= 0 || rows <= 0 || w <= 0 || h <= 0 || w > canvas.MaxPixels/h {
		return nil, fmt.Errorf("tilemap: a map of %dx%d cells", cols, rows)
	}
	m := &Map{Object: game.NewObject(name), Tiles: tiles, Cols: cols, Rows: rows, cells: make([]int, cols*rows), painted: make([]int, cols*rows)}
	for i := range m.cells {
		m.cells[i] = Empty
		m.painted[i] = Empty
	}
	m.AddScript(clock{m})
	m.AddScript(sparse{m})
	m.picture = &sprite.Sprite{}
	m.Collider = &solidTiles{m}
	return m, nil
}
