package art

// water is a water tile f steps into its ripple: the waves move two
// pixels a step, and repeat every eight, so tiles meet without a seam.
func water(f int) Grid {
	g := fill('v')
	for y := 0; y < Size; y++ {
		for x := 0; x < Size; x++ {
			switch {
			case y%4 == 1 && (x+2*f+y)%8 < 3:
				g = g.set(x, y, 'x')
			case y%4 == 3 && (x-2*f+y+16)%8 < 2:
				g = g.set(x, y, 'u')
			}
		}
	}
	return g
}
