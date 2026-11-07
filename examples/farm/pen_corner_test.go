package farm

import "testing"

// A pen's fence closes: a corner at each corner, whose rails go the way its
// two sides do, runs between them, and only its gate open.
func TestAPensFenceMeetsAtItsCorners(t *testing.T) {
	newTiles(SVGArt())
	r := rect{10, 20, penCols, penRows}
	want := map[[2]int]string{
		{r.col, r.row}: "fence_nw", {r.col + r.cols - 1, r.row}: "fence_ne",
		{r.col, r.row + r.rows - 1}: "fence_sw", {r.col + r.cols - 1, r.row + r.rows - 1}: "fence_se",
	}
	for at, name := range want {
		if got := fenceAt(r, at[0], at[1]); got != tileID[name] {
			t.Errorf("%v is %d, not %s", at, got, name)
		}
	}
	open := 0
	for row := r.row; row < r.row+r.rows; row++ {
		for col := r.col; col < r.col+r.cols; col++ {
			edge := row == r.row || row == r.row+r.rows-1 || col == r.col || col == r.col+r.cols-1
			switch id := fenceAt(r, col, row); {
			case !edge && id >= 0:
				t.Errorf("a fence inside the pen at %d, %d", col, row)
			case edge && id < 0:
				open++
			case edge && !isFence(id):
				t.Errorf("%d, %d is no fence", col, row)
			}
		}
	}
	if open != 1 {
		t.Errorf("%d cells of the fence are open, not the gate alone", open)
	}
}

// A pen's corners and sides stop whoever walks into them.
func TestAPensCornersAreSolid(t *testing.T) {
	_, f := newGame(idle{})
	r := penRect
	for _, at := range [][2]int{{r.col, r.row}, {r.col + r.cols - 1, r.row}, {r.col, r.row + r.rows - 1}, {r.col + r.cols - 1, r.row + r.rows - 1}} {
		if f.open(at[0], at[1]) {
			t.Errorf("the corner at %v may be walked through", at)
		}
	}
}
