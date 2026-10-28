package vector

import (
	"math"
	"testing"
)

func TestIdentityRotatesNothing(t *testing.T) {
	v := Vector3{1, 2, 3}
	near(t, "Rotate", Identity().Rotate(v), v)
}

func TestAxisAngleRotatesAQuarterTurn(t *testing.T) {
	q := AxisAngle(Vector3{Z: 2}, math.Pi/2)
	near(t, "Rotate", q.Rotate(Vector3{X: 1}), Vector3{Y: 1})
}

func TestMulComposesRotations(t *testing.T) {
	quarter := AxisAngle(Vector3{Z: 1}, math.Pi/2)
	half := quarter.Mul(quarter)
	near(t, "Rotate", half.Rotate(Vector3{X: 1}), Vector3{X: -1})
	first := AxisAngle(Vector3{X: 1}, math.Pi/2)
	near(t, "order", quarter.Mul(first).Rotate(Vector3{Y: 1}), quarter.Rotate(first.Rotate(Vector3{Y: 1})))
}
