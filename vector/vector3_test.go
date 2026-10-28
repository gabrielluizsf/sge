package vector

import "testing"

func TestArithmetic(t *testing.T) {
	a := Vector3{1, 2, 3}
	b := Vector3{4, 5, 6}
	near(t, "Add", a.Add(b), Vector3{5, 7, 9})
	near(t, "Sub", b.Sub(a), Vector3{3, 3, 3})
	near(t, "Scale", a.Scale(2), Vector3{2, 4, 6})
	near(t, "Mul", a.Mul(b), Vector3{4, 10, 18})
	near(t, "Cross", Vector3{X: 1}.Cross(Vector3{Y: 1}), Vector3{Z: 1})
	if got := a.Dot(b); got != 32 {
		t.Errorf("Dot = %v, want 32", got)
	}
}

func TestLength(t *testing.T) {
	v := Vector3{3, 4, 0}
	if got := v.Length(); got != 5 {
		t.Errorf("Length = %v, want 5", got)
	}
	if got := v.SqrLength(); got != 25 {
		t.Errorf("SqrLength = %v, want 25", got)
	}
	if got := v.Distance(Vector3{}); got != 5 {
		t.Errorf("Distance = %v, want 5", got)
	}
	near(t, "Normalized", v.Normalized(), Vector3{0.6, 0.8, 0})
}

func TestNormalizingNothingIsNothing(t *testing.T) {
	near(t, "Normalized", Vector3{}.Normalized(), Vector3{})
}

func TestLerp(t *testing.T) {
	a := Vector3{0, 0, 0}
	b := Vector3{10, -10, 4}
	near(t, "Lerp 0", a.Lerp(b, 0), a)
	near(t, "Lerp 1", a.Lerp(b, 1), b)
	near(t, "Lerp .5", a.Lerp(b, 0.5), Vector3{5, -5, 2})
}

func TestMethodsDoNotChangeTheReceiver(t *testing.T) {
	a := Vector3{1, 1, 1}
	_ = a.Add(Vector3{1, 1, 1})
	_ = a.Scale(3)
	near(t, "a", a, Vector3{1, 1, 1})
}
