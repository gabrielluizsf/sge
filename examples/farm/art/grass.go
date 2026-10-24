package art

// The tufts and specks of the grass, kept a pixel off every edge so tiles
// meet without a seam.
var (
	grassTufts  = [][2]int{{3, 2}, {11, 1}, {7, 6}, {13, 7}, {2, 10}, {9, 11}, {5, 13}, {12, 13}}
	grassSpecks = [][2]int{{8, 2}, {1, 5}, {14, 4}, {4, 8}, {14, 11}, {9, 14}, {1, 14}}
)

// grass is a grass tile with every tuft's tip leant by lean pixels: -1,
// 0 and 1 are the frames of the wind going through it.
func grass(lean int) Grid {
	g := fill('2')
	for _, p := range grassSpecks {
		g = g.set(p[0], p[1], '1')
	}
	for _, t := range grassTufts {
		x, y := t[0], t[1]
		g = g.set(x-1, y+1, '3').set(x+1, y+1, '3').set(x, y+2, '1').set(x+lean, y, '3')
	}
	return g
}

// flowers is grass with flowers in it, nodding with the grass.
func flowers(lean int) Grid {
	g := grass(lean)
	for _, f := range [][2]int{{4, 4}, {11, 9}, {6, 12}} {
		x, y := f[0]+lean, f[1]
		g = g.set(x, y-1, 'p').set(x-1, y, 'p').set(x+1, y, 'p').set(x, y+1, 'p').set(x, y, 'y')
	}
	return g
}
