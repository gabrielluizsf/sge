package farm

import (
	"fmt"

	"github.com/gabrielluizsf/antui/canvas"
)

// A piece of furniture is laid on the cells of a room one tile for each
// cell it covers: the picture cut along the cells' diagonals, so each cell's
// part is drawn among the people at that cell's depth. The tiles' names
// are the piece's, the heading and the cell's place in the piece.

// pieceSlice is the most cells a piece may cover across or deep.
const pieceSlice = 2

// pieceTileName is the tile of the cell i, j of the piece name turned
// heading.
func pieceTileName(name string, heading, i, j int) string {
	return fmt.Sprint(name, "~", heading, "~", i, j)
}

// pieceTileNames are all the pieces' tiles.
func pieceTileNames() []string {
	var out []string
	for _, n := range furnishings {
		for h := range headings {
			for i := range pieceSlice {
				for j := range pieceSlice {
					out = append(out, pieceTileName(n, h, i, j))
				}
			}
		}
	}
	return out
}

// pieceTileID is the tile of cell i, j of piece p.
func pieceTileID(p piece, i, j int) int {
	return tileID[pieceTileName(p.name, p.heading, i, j)]
}

// pieceRefs says, by tile name, which cell of which piece a tile is.
var pieceRefs = pieceRefsOf()

type pieceRef struct {
	name    string
	heading int
	cell    [2]int
}

func pieceRefsOf() map[string]pieceRef {
	m := map[string]pieceRef{}
	for _, n := range furnishings {
		for h := range headings {
			for i := range pieceSlice {
				for j := range pieceSlice {
					m[pieceTileName(n, h, i, j)] = pieceRef{n, h, [2]int{i, j}}
				}
			}
		}
	}
	return m
}

// drawTile is the picture of tile name in a: the art's, or a cell of a
// piece of furniture cut from its.
func drawTile(a Art, name string) *canvas.Canvas {
	if r, ok := pieceRefs[name]; ok {
		return pieceTile(a, r.name, r.heading, r.cell[0], r.cell[1])
	}
	return a.Tile(name)
}
