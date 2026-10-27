package art

// A tile's picture: a cell's diamond, IsoW wide and half that tall, at
// the bottom of a picture IsoH tall, so what stands on the cell reaches
// up out of it. A cell is Size units square in the world's terms.
const (
	Size = 16
	IsoW = 2 * Size
	IsoH = 3 * Size
)
