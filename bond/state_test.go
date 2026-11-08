package bond

import "testing"

func TestStateRoundTrip(t *testing.T) {
	b := Bond{Affection: 7, Courted: Engaged, gifted: 3, talked: 4, wedding: 9}
	if got := b.State().Bond(); got != b {
		t.Errorf("%+v, want %+v", got, b)
	}
}
