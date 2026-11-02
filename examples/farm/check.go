package farm

import (
	"fmt"
	"math"
	"unicode"
	"unicode/utf8"
)

// inRange is nil if v is from lo to hi, both in, or says what was not.
func inRange(what string, v, lo, hi int) error {
	if v < lo || v > hi {
		return fmt.Errorf("save: %s is %d, not from %d to %d", what, v, lo, hi)
	}
	return nil
}

// finite is nil if v is a number from lo to hi.
func finite(what string, v, lo, hi float64) error {
	if math.IsNaN(v) || v < lo || v > hi {
		return fmt.Errorf("save: %s is %v, not from %v to %v", what, v, lo, hi)
	}
	return nil
}

// plain is nil if text is a short text with nothing odd in it, or empty
// when empty is allowed.
func plain(what, text string) error {
	if text == "" || len(text) > maxText || !utf8.ValidString(text) {
		return fmt.Errorf("save: %s is not a short text", what)
	}
	for _, r := range text {
		if unicode.IsControl(r) {
			return fmt.Errorf("save: %s has a control character", what)
		}
	}
	return nil
}

// firstError is the first of errs that is not nil.
func firstError(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}

// checkPoint is nil if p is a place in the world.
func checkPoint(what string, p point) error {
	return firstError(finite(what+" x", p.X, -maxCoord, maxCoord), finite(what+" y", p.Y, -maxCoord, maxCoord))
}

// checkCell is nil if c is a cell the world has.
func checkCell(what string, col, row int) error {
	return firstError(inRange(what+" col", col, -maxCell, maxCell), inRange(what+" row", row, -maxCell, maxCell))
}
