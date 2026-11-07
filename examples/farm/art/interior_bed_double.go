package art

// Double bed is a piece of furniture of two cells across and two deep,
// drawn from its back cell: it reaches forward over the cells in front of
// it. Its picture is two cells wide and a cell and a half taller than a
// cell's, so its head is at the back, up to the left, as a real bed's is.
const (
	DoubleBedW = 2 * IsoW
	DoubleBedH = IsoH + IsoW/2
)

// bedDoubleTile is the couple's bed: a frame with a high headboard and a
// low footboard, a mattress, two pillows and a blanket over the foot half.
func bedDoubleTile(blanket [3]string) string {
	d := newDoc(DoubleBedW, DoubleBedH)
	d.origin = pt{DoubleBedW / 2, IsoH - 16}
	p := d.point(2*Size/2, 2*Size/2, 0)
	d.add(`<ellipse cx="%.2f" cy="%.2f" rx="%g" ry="%g" fill="%s" opacity="0.35"/>`, p.x, p.y, 15.0, 7.0, shadowInk)
	d.box([4]float64{0, 3, 1.6, 29}, [2]float64{0, 14}, wood3)
	d.box([4]float64{1.6, 3, 31, 29}, [2]float64{0, 5}, wood3)
	d.box([4]float64{1.6, 3.5, 30.5, 28.5}, [2]float64{5, 8}, [3]string{white, shirtBase, shirtShade})
	d.box([4]float64{2.5, 5, 8.5, 14}, [2]float64{8, 10}, [3]string{white, white, shirtShade})
	d.box([4]float64{2.5, 18, 8.5, 27}, [2]float64{8, 10}, [3]string{white, white, shirtShade})
	d.box([4]float64{13, 3.8, 30.5, 28.2}, [2]float64{8, 8.7}, blanket)
	d.box([4]float64{29.4, 3, 31, 29}, [2]float64{0, 8.5}, wood3)
	return d.String()
}
