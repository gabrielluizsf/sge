package farm

// FurnitureArt is what an Art may also have whose furniture is larger
// than a cell, as a real bed is: how many cells across and deep the piece
// name, one of the houses' insides, covers. A piece is drawn from its
// front cell, the one of the most column and row, and its picture reaches
// back over the others. Without it each piece covers its own cell.
type FurnitureArt interface {
	Furniture(name string) (cols, rows int)
}

// pieceNames are the furniture's names by their letters in a room's
// layout.
var pieceNames = map[byte]string{'b': "bed", 'B': "bed_blue", 't': "table", 'c': "chair", 's': "shelf", 'o': "stove", 'p': "plant", 'h': "chest"}

// spanOf is how many cells each piece of furniture covers in a, by its
// letter.
func spanOf(a Art) func(piece byte) (cols, rows int) {
	fa, ok := a.(FurnitureArt)
	return func(piece byte) (int, int) {
		if !ok {
			return 1, 1
		}
		return fa.Furniture(pieceNames[piece])
	}
}
