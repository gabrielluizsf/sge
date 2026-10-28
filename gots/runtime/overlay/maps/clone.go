//go:build gots_overlay

// Overlay for maps.Clone, whose copy the Go runtime makes: here it is a
// loop, keeping nil nil.

package maps

func clone(m any) any { panic("gots: maps.clone is not called") }

func Clone[M ~map[K]V, K comparable, V any](m M) M {
	if m == nil {
		return nil
	}
	r := make(M, len(m))
	for k, v := range m {
		r[k] = v
	}
	return r
}
