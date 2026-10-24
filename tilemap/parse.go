package tilemap

import (
	"fmt"
	"strings"
)

// Parse is a map drawn as text, a line a row and a character a cell, with
// legend saying which tile each character is. A character it does not
// name is Empty, and rows shorter than the longest are filled with Empty.
func Parse(name string, tiles *TileSet, rows string, legend map[rune]int) (*Map, error) {
	lines := strings.Split(strings.Trim(rows, "\n"), "\n")
	cols := 0
	for _, l := range lines {
		cols = max(cols, len([]rune(l)))
	}
	m, err := New(name, tiles, cols, len(lines))
	if err != nil {
		return nil, fmt.Errorf("tilemap: parsing %s: %w", name, err)
	}
	for row, l := range lines {
		for col, ch := range []rune(l) {
			id, ok := legend[ch]
			if !ok {
				id = Empty
			}
			m.Set(col, row, id)
		}
	}
	return m, nil
}
