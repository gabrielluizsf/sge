//go:build gots_overlay

// Constant-time choices between Elements: see element.go.

package field

// Select sets v to a if cond == 1, and to b if cond == 0.
func (v *Element) Select(a, b *Element, cond int) *Element {
	m := -int32(cond)
	v.l0 = b.l0 ^ (a.l0^b.l0)&m
	v.l1 = b.l1 ^ (a.l1^b.l1)&m
	v.l2 = b.l2 ^ (a.l2^b.l2)&m
	v.l3 = b.l3 ^ (a.l3^b.l3)&m
	v.l4 = b.l4 ^ (a.l4^b.l4)&m
	v.l5 = b.l5 ^ (a.l5^b.l5)&m
	v.l6 = b.l6 ^ (a.l6^b.l6)&m
	v.l7 = b.l7 ^ (a.l7^b.l7)&m
	v.l8 = b.l8 ^ (a.l8^b.l8)&m
	v.l9 = b.l9 ^ (a.l9^b.l9)&m
	v.l10 = b.l10 ^ (a.l10^b.l10)&m
	v.l11 = b.l11 ^ (a.l11^b.l11)&m
	return v
}

// Swap swaps v and u if cond == 1 or leaves them unchanged if cond == 0, and returns v.
func (v *Element) Swap(u *Element, cond int) {
	var t Element
	t.Select(u, v, cond)
	u.Select(v, u, cond)
	*v = t
}
