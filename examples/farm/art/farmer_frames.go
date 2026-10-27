package art

// stepFront is a front or back frame with one leg, left or right, lifted
// a pixel mid-stride, its boot off the ground.
func stepFront(g Grid, left bool) Grid {
	out := g.copyOf()
	x0, x1 := 4, 9
	if !left {
		x0, x1 = 10, 15
	}
	for y := 21; y < FarmerH; y++ {
		r := []rune(out[y])
		for x := x0; x <= x1; x++ {
			below := '.'
			if y+1 < FarmerH {
				below = []rune(g[y+1])[x]
			}
			r[x] = below
		}
		out[y] = string(r)
	}
	return out
}

// stepSide is a side frame with its legs mid-stride.
func stepSide(g Grid) Grid {
	out := g.copyOf()
	copy(out[FarmerH-len(sideStride):], sideStride)
	return out
}

// can is the watering can, pouring from its spout on the right.
var can = grid(`
.kkk..
kMMMk.
kMmMkx
.kkk.x
.....x`)

// holding is g with the can held at x, y, mirrored to pour leftwards.
func holding(g Grid, x, y int, leftwards bool) Grid {
	c := can
	if leftwards {
		c = can.mirror()
	}
	out := g.copyOf()
	for dy, row := range c {
		for dx, ch := range []rune(row) {
			if ch != '.' {
				out = out.set(x+dx, y+dy, ch)
			}
		}
	}
	return out
}

// Farmer is every frame of the farmer, by direction and what it shows.
var Farmer = map[string]Grid{
	"down":    farmerDown,
	"down1":   stepFront(farmerDown, true),
	"down2":   stepFront(farmerDown, false),
	"downAct": holding(farmerDown, 14, 18, false),
	"up":      farmerUp,
	"up1":     stepFront(farmerUp, true),
	"up2":     stepFront(farmerUp, false),
	"upAct":   holding(farmerUp, 0, 18, true),
	"right":   farmerRight,
	// From the side both strides look the same: the walk passes through
	// standing between them.
	"right1":   stepSide(farmerRight),
	"right2":   stepSide(farmerRight),
	"rightAct": holding(farmerRight, 13, 16, false),
}

func init() {
	for _, f := range []string{"", "1", "2", "Act"} {
		Farmer["left"+f] = Farmer["right"+f].mirror()
	}
}
