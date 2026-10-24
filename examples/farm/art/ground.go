package art

// path is a dirt path with pebbles.
var path = grid(`
cccccccccccccccc
ccccbccccccccccc
cccccccccccbcccc
ccbccccccccccccc
ccccccccbbcccccc
cccccccccccccbcc
cbcccccccccccccc
ccccccbccccccccc
cccccccccccccccc
cccbcccccccbcccc
cccccccccccccccc
cccccccacccccccc
ccccccccccccbccc
cbcccccccccccccc
cccccbcccccccccc
cccccccccccccccc`)

// soil is tilled earth, dry: ridges catching the light, furrows between.
func soil(base, ridge, furrow rune) Grid {
	g := fill(base)
	for y := 0; y < Size; y++ {
		for x := 0; x < Size; x++ {
			switch y % 4 {
			case 0:
				g = g.set(x, y, ridge)
			case 2:
				if (x+y/4)%3 != 0 {
					g = g.set(x, y, furrow)
				}
			}
		}
	}
	return g
}

var (
	dry = soil('b', 'c', 'a')
	wet = soil('e', 'a', 'k')
)
