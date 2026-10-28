package vector

import "math"

// Length is the euclidean length of v.
func (v Vector3) Length() float64 { return math.Sqrt(v.SqrLength()) }
