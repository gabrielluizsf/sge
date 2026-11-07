// Package colourcode reads a colour a person typed: #RGB, #RRGGBB,
// #RRGGBBAA, rgb(r, g, b) and rgba(r, g, b, a), and writes one back. It
// reads text that may be anything, so it is strict: a length limit, digits
// only, every number in range, and no regular expression.
package colourcode

import (
	"errors"
	"strings"
)

// MaxLen is the longest text Parse looks at, in bytes.
const MaxLen = 40

// ErrInvalid is what Parse gives for text that is not a colour.
var ErrInvalid = errors.New("colourcode: not a colour")

// RGBA is a colour of 8 bits a channel, A 255 being opaque.
type RGBA struct{ R, G, B, A uint8 }

// Parse is the colour s spells. Spaces around it, and inside rgb( ), are
// free; letters may be either case. Alpha is 0 to 1 in rgba, as in CSS,
// or a percent; in hex it is the last pair.
func Parse(s string) (RGBA, error) {
	s = strings.TrimSpace(s)
	if len(s) == 0 || len(s) > MaxLen {
		return RGBA{}, ErrInvalid
	}
	if s[0] == '#' {
		return parseHex(s[1:])
	}
	return parseFunc(strings.ToLower(s))
}

// Allowed reports whether r may be typed into a colour code.
func Allowed(r rune) bool {
	switch {
	case r >= '0' && r <= '9', r >= 'a' && r <= 'f', r >= 'A' && r <= 'F':
		return true
	}
	return strings.ContainsRune("#(),. %rgbaRGBA", r)
}

// Insert is cur with the allowed characters of typed added, up to MaxLen.
func Insert(cur, typed string) string {
	var b strings.Builder
	b.WriteString(cur)
	for _, r := range typed {
		if b.Len() >= MaxLen {
			break
		}
		if Allowed(r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// Format is c as #RRGGBB, or #RRGGBBAA when it is not opaque.
func Format(c RGBA) string {
	const digits = "0123456789ABCDEF"
	out := []byte("#")
	vals := []uint8{c.R, c.G, c.B}
	if c.A != 255 {
		vals = append(vals, c.A)
	}
	for _, v := range vals {
		out = append(out, digits[v>>4], digits[v&15])
	}
	return string(out)
}
