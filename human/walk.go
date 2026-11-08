package human

import "math"

// The left leg's joints through a stride, in degrees, sixteen to a stride
// from its heel's strike: the normal adult gait of Perry and Burnfield's
// Gait Analysis and Winter's Biomechanics and Motor Control of Human Gait.
// The hip is flexed at the strike, goes back to extend as the body passes
// over the foot, and swings forward in the air; the knee gives a little
// as the weight comes on, straightens, then folds to clear the ground, at
// most 60 in the swing; the ankle drops flat, rolls over the foot, pushes
// off, and lifts the toes through the swing. The right leg is the same
// half a stride on.
var (
	strideHip   = [16]float64{30, 28, 22, 14, 5, -2, -8, -12, -14, -10, -2, 10, 22, 30, 34, 33}
	strideKnee  = [16]float64{5, 12, 16, 14, 8, 5, 8, 20, 18, 28, 42, 58, 62, 48, 25, 10}
	strideAnkle = [16]float64{0, 5, 3, -2, -6, -9, -11, -10, -6, 6, 16, 6, -2, -3, -2, -1}
)

// Walking's other numbers. The arms swing against the legs, 20 degrees
// to a side at the shoulder, the elbows a little bent, more as the arm
// comes forward. The feet land close to the body's middle line: each leg
// is brought in 2.5 degrees, which puts the ankles about a hand's width
// closer than the hips, a step's width of some 15 centimetres.
const (
	armSwing   = 20 * math.Pi / 180
	elbowBend  = 25 * math.Pi / 180
	legIn      = -0.045
	thighLen   = 0.45 // metres, a thigh and a shank, for the body's rise
	shankLen   = 0.44
	bobEase    = 0.6 // how much of the hips' fall the foot's roll gives back
	pelvisTurn = 0.07
)

// at is v's value k keys on, round a stride.
func at(v [16]float64, k int) float64 { return v[((k%16)+16)%16] * math.Pi / 180 }

// legReach is how far below the hip a leg's ankle hangs, in metres, its
// thigh swung hip and its knee bent knee radians.
func legReach(hip, knee float64) float64 {
	return thighLen*math.Cos(hip) + shankLen*math.Cos(hip-knee)
}

// Walk is walking, a stride a second: each leg swings through a normal
// gait's hip, knee and ankle, heel first, off the toes last; the arms
// swing against the legs; the body rises over each foot and drops as the
// next lands.
var Walk = loop{"walk", 1, 16}.of(func(t float64) Pose {
	k := int(math.Round(t * 16))
	var p Pose
	for s, key := range [2]int{k, k + 8} {
		p.Hip[s], p.Knee[s], p.Ankle[s] = at(strideHip, key), at(strideKnee, key), at(strideAnkle, key)
		p.Abduct[s] = legIn
	}
	// Each hip's swing, in about -1 to 1: an arm is forward as its own
	// leg is back.
	mid, span := 10*math.Pi/180, 24*math.Pi/180
	nL, nR := (p.Hip[0]-mid)/span, (p.Hip[1]-mid)/span
	p.Shoulder = [2]float64{-armSwing * nL, -armSwing * nR}
	p.Raise = both(0.08)
	p.Elbow = [2]float64{elbowBend + 0.2*max(0, -nL), elbowBend + 0.2*max(0, -nR)}
	p.Twist = pelvisTurn * nL
	p.Lean, p.Nod = 0.05, 0.03
	// The body is as high as the hips are above the lower foot; some of
	// the fall is given back by the foot's roll.
	reach := max(legReach(p.Hip[0], p.Knee[0]), legReach(p.Hip[1], p.Knee[1]))
	p.Bob = bobEase * (reach - (thighLen + shankLen))
	return p
})
