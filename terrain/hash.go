package terrain

// Hash is a number that looks random but is the same for the same seed and
// cell every time. It works in 32 bits, which JavaScript does fast.
func Hash(seed uint32, col, row int) uint32 {
	h := seed ^ uint32(col)*0x27d4eb2d ^ uint32(row)*0x165667b1
	h ^= h >> 15
	h *= 0x85ebca6b
	h ^= h >> 13
	h *= 0xc2b2ae35
	h ^= h >> 16
	return h
}
