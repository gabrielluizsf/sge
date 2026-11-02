// Package textfield is what typing into a one-line field needs, without
// any window: which characters a name may be made of, adding the ones
// typed up to a limit, taking the last one away, and checking a name
// read from elsewhere. A game feeds it the text its window says was
// typed this frame.
package textfield

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// Allowed reports whether r may be in a name: a letter of any script
// (accents included), a space, a hyphen or an apostrophe. Control and
// format characters, digits, marks and other symbols are not.
func Allowed(r rune) bool {
	return unicode.IsLetter(r) || r == ' ' || r == '-' || r == '\''
}

// Insert is cur with the allowed characters of typed added at its end, no
// more than max characters in all; the rest of typed is dropped.
func Insert(cur, typed string, max int) string {
	n := utf8.RuneCountInString(cur)
	var b strings.Builder
	b.WriteString(cur)
	for _, r := range typed {
		if n >= max {
			break
		}
		if r == utf8.RuneError || !Allowed(r) {
			continue
		}
		b.WriteRune(r)
		n++
	}
	return b.String()
}

// Backspace is cur without its last character.
func Backspace(cur string) string {
	_, size := utf8.DecodeLastRuneInString(cur)
	return cur[:len(cur)-size]
}

// Valid reports whether s is a good name: valid UTF-8 of 1 to max
// characters, all Allowed, not beginning or ending with a space.
func Valid(s string, max int) bool {
	n := utf8.RuneCountInString(s)
	if !utf8.ValidString(s) || n < 1 || n > max || s != strings.TrimSpace(s) || strings.HasPrefix(s, " ") {
		return false
	}
	for _, r := range s {
		if !Allowed(r) {
			return false
		}
	}
	return true
}

// Trim is s without the spaces at its ends.
func Trim(s string) string { return strings.TrimSpace(s) }
