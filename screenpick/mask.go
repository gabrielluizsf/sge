package screenpick

import "math/bits"

// channelOf is the part of pixel that mask picks out, from 0 to 255.
func channelOf(pixel, mask uint64) uint8 {
	if mask == 0 {
		return 0
	}
	shift := bits.TrailingZeros64(mask)
	width := bits.Len64(mask >> shift)
	v := (pixel & mask) >> shift
	switch {
	case width == 8:
		return uint8(v)
	case width > 8:
		return uint8(v >> (width - 8))
	}
	return uint8(v * 255 / (1<<width - 1))
}
