package human

// The crawl's stroke of one arm, at eight moments evenly round a stroke,
// in radians. The arm is lifted Raise from the body's side, 0 at the hip
// and pi reaching ahead, and turned Plane from the water's surface, down
// into it when positive, up out of it when negative; the Elbow bends.
// So the hand reaches (the entry), catches, pulls under the chest, pushes
// to the thigh, leaves, and comes back through the air by the side, the
// elbow high and the hand hanging from it: a stroke of about 140 degrees
// at the shoulder, as a swimmer's is, not a full turn.
var (
	strokeRaise = []float64{2.70, 2.40, 1.60, 0.70, 0.35, 1.00, 1.70, 2.30}
	strokePlane = []float64{0.15, 0.60, 1.10, 0.70, -0.10, -0.50, -0.45, -0.20}
	strokeElbow = []float64{0.25, 1.00, 1.50, 0.90, 0.50, 1.50, 1.40, 0.70}
)

// rollAmount is how far the body rolls to each side, in radians: some 35
// degrees at the trunk, in the 35 to 60 swimmers roll (Payton et al.,
// Journal of Applied Biomechanics, 1997; Yanai, 2003).
const rollAmount = 0.6

// Swim is the crawl, a stroke of each arm in turn, 1.6 seconds to the
// pair: lying on the water, the body rolling to the arm that pulls, the
// legs kicking six beats to a pair of strokes, the head up.
var Swim = loop{"swim", 1.6, 16}.of(func(t float64) Pose {
	kick := wave(t, 3, 0)
	arm := func(phase float64) (raise, plane, elbow float64) {
		return smooth(strokeRaise, t+phase), smooth(strokePlane, t+phase), smooth(strokeElbow, t+phase)
	}
	rl, pl, el := arm(0)
	rr, pr, er := arm(0.5)
	return Pose{Pitch: 1.45, Nod: -0.9, Twist: rollAmount * wave(t, 1, 0),
		Hip: [2]float64{0.25 * kick, -0.25 * kick}, Knee: [2]float64{0.2 + 0.2*max(0, kick), 0.2 + 0.2*max(0, -kick)},
		Ankle: both(0.85), Abduct: both(0.1),
		Raise: [2]float64{rl, rr}, Plane: [2]float64{pl, pr}, Elbow: [2]float64{el, er},
		Bob: -1.02}
})
