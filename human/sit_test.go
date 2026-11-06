package human

import (
	"math"
	"testing"
)

// Seated, the thighs are level, the shanks plumb, and the body sits as low
// as the legs leave it, with the feet on the floor and every joint healthy.
func TestSittingHasLevelThighsBentKneesAndFeetOnTheFloor(t *testing.T) {
	for k, p := range Sit.Poses {
		if b := p.Beyond(); len(b) > 0 {
			t.Errorf("pose %d: %v", k, b)
		}
		for s := range 2 {
			if h := p.Hip[s]; math.Abs(h-math.Pi/2) > 0.15 {
				t.Errorf("pose %d: the thigh is %.2f rad from level", k, h-math.Pi/2)
			}
			if p.Knee[s] < 1.2 {
				t.Errorf("pose %d: the knee is hardly bent: %.2f", k, p.Knee[s])
			}
		}
		// The hips are as high above the floor as the legs reach: about a
		// chair's seat, a shank's length.
		hips := thighLen + shankLen + p.Bob
		if math.Abs(hips-legReach(p.Hip[0], p.Knee[0])) > 0.01 || math.Abs(hips-shankLen) > 0.06 {
			t.Errorf("pose %d: the hips are %.3f m up", k, hips)
		}
	}
	if len(Sit.Poses) == 0 || Sit.Length <= 0 {
		t.Error("an empty motion")
	}
}
