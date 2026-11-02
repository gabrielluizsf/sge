package farm

// marks reports whether the cursor shows on cell col, row for what the
// player holds: a plot, a tree with the axe, water or a well with the can.
func (f *farm) marks(col, row int) bool {
	switch {
	case f.plotAt(col, row) != nil:
		return true
	case f.slot == slotAxe:
		return isTree(f.thingAt(col, row))
	case f.slot == slotCan:
		return f.isWater(col, row)
	}
	return false
}
