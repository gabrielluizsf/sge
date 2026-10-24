package art

import "strings"

// Grid is a sprite: Size rows of Size palette letters.
type Grid []string

// grid is the Grid a block of text draws, one line a row.
func grid(text string) Grid {
	return Grid(strings.Split(strings.Trim(text, "\n"), "\n"))
}

// mirror is g flipped left to right: a sprite walking left is one walking
// right, turned round.
func (g Grid) mirror() Grid {
	out := make(Grid, len(g))
	for i, row := range g {
		r := []rune(row)
		for a, b := 0, len(r)-1; a < b; a, b = a+1, b-1 {
			r[a], r[b] = r[b], r[a]
		}
		out[i] = string(r)
	}
	return out
}

// set is g with the pixel x, y made letter c.
func (g Grid) set(x, y int, c rune) Grid {
	r := []rune(g[y])
	r[x] = c
	g[y] = string(r)
	return g
}

// copyOf is a copy of g to change.
func (g Grid) copyOf() Grid {
	return append(Grid(nil), g...)
}
