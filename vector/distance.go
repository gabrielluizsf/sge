package vector

// Distance is how far v is from o.
func (v Vector3) Distance(o Vector3) float64 { return v.Sub(o).Length() }
