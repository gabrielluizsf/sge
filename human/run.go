package human

import "math"

// The left leg's joints through a stride of running at about 3.5 to 4
// metres a second, in degrees, sixteen to a stride from the foot's
// strike, after Novacheck's The Biomechanics of Running: the foot lands
// under the hips on the midfoot, the knee gives to 40 as the weight comes
// on and the ankle rolls over it, the hip extends to push off at about
// two fifths of the stride, and the thigh is then driven forward to about
// 60 degrees while the knee folds to over 100, the heel coming up to the
// buttock, and unfolds to land again.
var (
	runHip   = [16]float64{28, 22, 14, 6, -2, -9, -15, 5, 25, 42, 55, 60, 55, 45, 36, 30}
	runKnee  = [16]float64{18, 30, 40, 42, 36, 28, 30, 55, 85, 105, 100, 85, 65, 45, 28, 20}
	runAnkle = [16]float64{5, -10, -18, -20, -12, 5, 25, 25, 10, 0, -5, -8, -5, 0, 5, 5}
)

// Running's other numbers: the arms swing 40 degrees to a side of a mean a little forward,
// from the hip to the chest, with the elbows bent to 100 degrees or more, the pelvis and the shoulders turn
// against each other, the trunk leans from the ankles, and in flight the
// body rises 4 centimetres above where the stance leg holds it.
const (
	runArmMean = 0.2  // the arm's swing from the vertical, at the middle of it
	runArmFore = 0.95 // and how far more it goes forward, and
	runArmBack = 0.55 // back, as its leg's hip swings
	runElbow   = 1.75
	runLean    = 0.17
	runFlight  = 0.05
	runStrikes = 8 // keys from one leg's strike to the other's
	runOffKey  = 6 // the key a leg leaves the ground at
)

// footClearance is how high above the ground a pose has the foot of leg s,
// in metres: the body's standing height, lifted by the pose's bob, less
// how far below the hip the foot hangs.
func footClearance(p Pose, s int) float64 {
	return thighLen + shankLen + p.Bob - legReach(p.Hip[s], p.Knee[s])
}

// flightLift is how far the body is thrown above where its stance leg
// holds it at key k of a stride: not at all while a foot is down, and in
// the two flights, from each foot's leaving to the other's landing, a
// rise and fall.
func flightLift(k int) float64 {
	phase := (k - runOffKey + 2*runStrikes) % runStrikes
	if phase > runStrikes-runOffKey || phase > 2 {
		return 0
	}
	return runFlight * math.Pow(math.Sin(math.Pi*float64(phase)/2), 2)
}

// runArmAt is how far forward an arm is, in radians, as its leg's hip is
// at n, from -1 forward to 1 back: forward further than it goes back, as
// the hand rises to the chest and falls to the hip.
func runArmAt(n float64) float64 {
	if n < 0 {
		return runArmMean - runArmFore*n
	}
	return runArmMean - runArmBack*n
}

// Run is running, a pair of strides in 0.7 seconds: each leg goes through
// a running stride, the right half a stride after the left, the arms
// swing against the legs, and between the strides the body is in flight.
var Run = loop{"run", 0.7, 16}.of(func(t float64) Pose {
	k := int(math.Round(t * 16))
	var p Pose
	for s, key := range [2]int{k, k + runStrikes} {
		p.Hip[s], p.Knee[s], p.Ankle[s] = at(runHip, key), at(runKnee, key), at(runAnkle, key)
		p.Abduct[s] = legIn
	}
	mid, span := 15*math.Pi/180, 38*math.Pi/180
	nL, nR := (p.Hip[0]-mid)/span, (p.Hip[1]-mid)/span
	p.Shoulder = [2]float64{runArmAt(nL), runArmAt(nR)}
	p.Raise = both(0.1)
	p.Elbow = [2]float64{runElbow + 0.1*max(0, -nL), runElbow + 0.1*max(0, -nR)}
	p.Twist = 0.12 * nL
	p.Lean, p.Nod = runLean, -0.05
	reach := max(legReach(p.Hip[0], p.Knee[0]), legReach(p.Hip[1], p.Knee[1]))
	p.Bob = reach - (thighLen + shankLen) + flightLift(k)
	return p
})
