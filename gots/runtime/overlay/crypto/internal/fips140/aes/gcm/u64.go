//go:build gots_overlay

// A 64-bit word as two 32-bit ones: see ghash.go.

package gcm

type u64 struct{ h, l uint32 }

func (a u64) xor(b u64) u64 { return u64{a.h ^ b.h, a.l ^ b.l} }

// shl is a << s, for s from 1 to 63.
func (a u64) shl(s uint) u64 {
	if s >= 32 {
		return u64{a.l << (s - 32), 0}
	}
	return u64{a.h<<s | a.l>>(32-s), a.l << s}
}

// shr is a >> s, for s from 1 to 63.
func (a u64) shr(s uint) u64 {
	if s >= 32 {
		return u64{0, a.h >> (s - 32)}
	}
	return u64{a.h >> s, a.l>>s | a.h<<(32-s)}
}
