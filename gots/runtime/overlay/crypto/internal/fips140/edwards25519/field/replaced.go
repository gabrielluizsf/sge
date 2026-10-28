//go:build gots_overlay

// Go's declarations on its limbs that element.go's Element has no use
// for, put out of the way: nothing calls them.

package field

func (v *Element) reduce() *Element { return v }

func (v *Element) carryPropagate() *Element { return v }

func feMulGeneric(v, a, b *Element) { feMul(v, a, b) }

func feSquare(v, a *Element) { feMul(v, a, a) }

func feSquareN(v, a *Element, n int) { v.SquareN(a, n) }

func mul51(a uint64, b uint32) (lo uint64, hi uint64) { return }
