//go:build gots_overlay

// Overlay for nistec's precomputed P-256 table, which Go reads by casting
// bytes to elements with unsafe: here the limbs are set one by one.

package fiat

// GotsSetLimbs sets e to the element whose Montgomery limbs are l, least
// significant first, as the table stores them.
func (e *P256Element) GotsSetLimbs(l [4]uint64) *P256Element {
	e.x = p256MontgomeryDomainFieldElement(l)
	return e
}
