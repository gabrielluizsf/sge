//go:build gots_overlay

// Overlay for the size of int, which strconv works out from uint's: under
// gots uint has 32 bits but int has 64, a number.

package strconv

const intSize = 64
