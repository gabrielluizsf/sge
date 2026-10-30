package building

// Props are the things of a farm and a town that are not houses.
var Props = []Building{Stall, Smithy, Anvil, Well, FenceStraight, FenceCorner, FenceEnd, Bench, Wheelbarrow}

// All are every house and every prop.
func All() []Building {
	return append(append([]Building{}, Houses...), Props...)
}
