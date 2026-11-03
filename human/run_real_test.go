package human

import (
	"math"
	"testing"
)

// In a run the arms stay close and bent: in every pose of the real kit's
// person, for both arms, the wrist is within 36 cm of the shoulder and the
// elbow, measured from the bones' places, is bent between 70 and 110
// degrees; the hands swing from near the hip to near the chest.
func TestRunArmsStayBentAndClose(t *testing.T) {
	m, err := realKit(t).Make(baselines[0].body, baselines[0].look)
	if err != nil {
		t.Fatal(err)
	}
	if err := Animate(m, Run); err != nil {
		t.Fatal(err)
	}
	a := &m.Animations[0]
	low, high := 9.0, -9.0
	for k := range Run.Poses {
		w := worldAt(m, a, k)
		for _, side := range []string{".L", ".R"} {
			s, e, r := w[boneOf(m, "upperarm01"+side)].at, w[boneOf(m, "lowerarm01"+side)].at, w[boneOf(m, "wrist"+side)].at
			bend := math.Acos(dot3(normalize(sub(s, e)), normalize(sub(r, e)))) * 180 / math.Pi
			d := sub(r, s)
			if bend < 70 || bend > 110 {
				t.Errorf("pose %d%s: the elbow is bent %.0f degrees", k, side, bend)
			}
			if dot3(d, d) > 0.36*0.36 {
				t.Errorf("pose %d%s: the wrist is %.2f m from the shoulder", k, side, math.Sqrt(dot3(d, d)))
			}
			low, high = min(low, r[1]-s[1]), max(high, r[1]-s[1])
		}
	}
	if low > -0.3 || high < -0.12 {
		t.Errorf("the hands swing from %.2f to %.2f m below the shoulder", -low, -high)
	}
}
