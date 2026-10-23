package collider

// clamp is v held between lo and hi.
func clamp(v, lo, hi float64) float64 { return max(lo, min(v, hi)) }
