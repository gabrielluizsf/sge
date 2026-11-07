package farm

// fillAction is the action of filling a can from the water on cell col,
// row: drawn from a well's edge on a town's, and crouched at the bank for
// every other water, a pond, a river or a lake.
func (f *farm) fillAction(col, row int) action {
	if !f.onFarm(col, row) {
		if t := f.towns.At(col, row); t != nil && atWell(wellRect(f.art, t), col, row) {
			return actFillWell
		}
	}
	return actFillRiver
}
