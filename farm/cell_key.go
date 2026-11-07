package farm

// cellOfKey is the cell cellKey made key from.
func cellOfKey(key int) (col, row int) {
	const reach = 1 << 20
	return key/(2*reach) - reach, key%(2*reach) - reach
}
