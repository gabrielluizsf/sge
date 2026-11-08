package colourcode

// parseHex reads the digits after a #: 3, 6 or 8 of them.
func parseHex(h string) (RGBA, error) {
	var v [8]uint8
	for i := 0; i < len(h); i++ {
		d, ok := hexDigit(h[i])
		if !ok || i >= len(v) {
			return RGBA{}, ErrInvalid
		}
		v[i] = d
	}
	switch len(h) {
	case 3:
		return RGBA{v[0] * 17, v[1] * 17, v[2] * 17, 255}, nil
	case 6:
		return RGBA{v[0]<<4 | v[1], v[2]<<4 | v[3], v[4]<<4 | v[5], 255}, nil
	case 8:
		return RGBA{v[0]<<4 | v[1], v[2]<<4 | v[3], v[4]<<4 | v[5], v[6]<<4 | v[7]}, nil
	}
	return RGBA{}, ErrInvalid
}

// hexDigit is the value of the hex digit c.
func hexDigit(c byte) (uint8, bool) {
	switch {
	case c >= '0' && c <= '9':
		return c - '0', true
	case c >= 'a' && c <= 'f':
		return c - 'a' + 10, true
	case c >= 'A' && c <= 'F':
		return c - 'A' + 10, true
	}
	return 0, false
}
