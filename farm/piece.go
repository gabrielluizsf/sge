package farm

// piece is a piece of furniture standing in the player's house: what it
// is, the cell it stands on (the back one, of the least column and row,
// if it covers more than one) and which of four ways it is turned.
type piece struct {
	name     string
	col, row int
	heading  int
	// held is what a chest keeps, a slot each; nil for other pieces.
	held []chestSlot
}

// headings is how many ways a piece may be turned: each a quarter turn
// on from the one before.
const headings = 4

// furnishings are the pieces of furniture there are, which are items to
// carry: a piece in the bag is the item of its name.
var furnishings = []string{"bed", "table", "chair", "shelf", "stove", "plant", "chest"}

// isFurnishing reports whether id is a piece of furniture.
func isFurnishing(id string) bool {
	for _, n := range furnishings {
		if n == id {
			return true
		}
	}
	return false
}

// footprint is how many cells across and deep piece name covers when
// turned heading: the art's, with the two swapped on a quarter turn.
func footprint(a Art, name string, heading int) (cols, rows int) {
	cols, rows = 1, 1
	if fa, ok := a.(FurnitureArt); ok {
		cols, rows = fa.Furniture(name)
	}
	cols, rows = max(cols, 1), max(rows, 1)
	if heading%2 == 1 {
		cols, rows = rows, cols
	}
	return cols, rows
}

// cells is the span of cells p covers.
func (p piece) cells(a Art) span {
	cols, rows := footprint(a, p.name, p.heading)
	return span{c0: p.col, r0: p.row, cols: cols, rows: rows}
}

// turned is p turned a quarter turn on, about its back cell.
func (p piece) turned() piece {
	p.heading = (p.heading + 1) % headings
	return p
}

// defaultHouse is what the player's house has to begin with: the bed.
func defaultHouse() []piece { return []piece{{name: "bed", col: 1, row: 1}} }
