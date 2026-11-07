package colourcode

import "strings"

// parseFunc reads rgb(r, g, b) and rgba(r, g, b, a), already lower case.
func parseFunc(s string) (RGBA, error) {
	var want int
	switch {
	case strings.HasPrefix(s, "rgba("):
		s, want = s[5:], 4
	case strings.HasPrefix(s, "rgb("):
		s, want = s[4:], 3
	default:
		return RGBA{}, ErrInvalid
	}
	if !strings.HasSuffix(s, ")") {
		return RGBA{}, ErrInvalid
	}
	parts := strings.Split(strings.TrimSuffix(s, ")"), ",")
	if len(parts) != want {
		return RGBA{}, ErrInvalid
	}
	var v [3]uint8
	for i := 0; i < 3; i++ {
		n, ok := channel(strings.TrimSpace(parts[i]))
		if !ok {
			return RGBA{}, ErrInvalid
		}
		v[i] = n
	}
	c := RGBA{v[0], v[1], v[2], 255}
	if want == 4 {
		a, ok := alpha(strings.TrimSpace(parts[3]))
		if !ok {
			return RGBA{}, ErrInvalid
		}
		c.A = a
	}
	return c, nil
}

// channel is a whole number 0 to 255 of at most 3 digits.
func channel(p string) (uint8, bool) {
	if len(p) == 0 || len(p) > 3 {
		return 0, false
	}
	n := 0
	for i := 0; i < len(p); i++ {
		if p[i] < '0' || p[i] > '9' {
			return 0, false
		}
		n = n*10 + int(p[i]-'0')
	}
	return uint8(n), n <= 255
}

// alpha is an alpha as 0 to 1 ("0.5", ".5", "1") or a percent ("50%"),
// in 0 to 255. At most 4 digits are read after the point.
func alpha(p string) (uint8, bool) {
	scale := 1000
	if strings.HasSuffix(p, "%") {
		p, scale = strings.TrimSuffix(p, "%"), 10
	}
	if len(p) == 0 || len(p) > 8 {
		return 0, false
	}
	whole, frac, digits, dot, seen := 0, 0, 0, false, false
	for i := 0; i < len(p); i++ {
		switch c := p[i]; {
		case c == '.' && !dot && scale == 1000:
			dot = true
		case c >= '0' && c <= '9':
			seen = true
			if !dot {
				whole = whole*10 + int(c-'0')
				if whole > 1000 {
					return 0, false
				}
			} else if digits < 3 {
				frac, digits = frac*10+int(c-'0'), digits+1
			}
		default:
			return 0, false
		}
	}
	for ; digits < 3; digits++ {
		frac *= 10
	}
	v := whole*scale + frac // thousandths of 1
	if v > 1000 || !seen {
		return 0, false
	}
	return uint8((v*255 + 500) / 1000), true
}
