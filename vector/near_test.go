package vector

import (
	"math"
	"testing"
)

// near fails t when got is not within a rounding error of want.
func near(t *testing.T, what string, got, want Vector3) {
	t.Helper()
	if math.Abs(got.X-want.X) > 1e-9 || math.Abs(got.Y-want.Y) > 1e-9 || math.Abs(got.Z-want.Z) > 1e-9 {
		t.Errorf("%s = %v, want %v", what, got, want)
	}
}
