package farm

import "math"

// Facings are the ways a person may face on screen, for a Body of n
// facings: down, up, left and right, and with 8 the ways between them,
// "downleft", "downright", "upleft" and "upright".
func Facings(n int) []string {
	if n == 8 {
		return eightWays[:]
	}
	return []string{"down", "up", "left", "right"}
}

// eightWays are the eight facings, clockwise on screen from the right.
var eightWays = [8]string{"right", "downright", "down", "downleft", "left", "upleft", "up", "upright"}

// facingOf is the facing, of n, nearest the screen's way x, y: with 4,
// the side it leans to most, or up or down when it leans to neither.
func facingOf(x, y float64, n int) string {
	if n == 8 {
		i := int(math.Round(math.Atan2(y, x) / (math.Pi / 4)))
		return eightWays[(i%8+8)%8]
	}
	if math.Abs(x) > math.Abs(y) {
		return map[bool]string{true: "right", false: "left"}[x > 0]
	}
	return map[bool]string{true: "down", false: "up"}[y > 0]
}
