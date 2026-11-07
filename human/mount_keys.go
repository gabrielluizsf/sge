package human

import "math"

// How long getting on takes, how long getting down, which is done the moment
// it is asked for, and how many frames of either are made.
const (
	MountLength    = 1.8
	DismountLength = 0.7
	mountFrames    = 16
)

// The spots on a horse's frame and by it where the feet go: each stirrup, where
// the seat puts the ankle; those of standing beside it, flat, the ones of being
// turned to it; and the swinging foot's, at the top of its way over the croup.
var (
	nearStirrup = [3]float64{-0.279, 0.814, -0.065}
	farStirrup  = [3]float64{0.279, 0.814, -0.065}
	standSwing  = [3]float64{-0.558, 0.091, -0.004}
	standNear   = [3]float64{-0.842, 0.091, -0.004}
	backOnFloor = [3]float64{-0.67, 0.091, -0.32}
	headOnFloor = [3]float64{-0.73, 0.091, -0.04}
	overCroup   = [3]float64{0.03, 1.56, -0.85}
)

// Toes level on the floor, a little down in the air.
var (
	flat      = [2]float64{0, 0}
	inAir     = [2]float64{0.35, -0.12}
	bothInAir = [2]float64{0.35, 0.35}
	seated    = [2]float64{rideDroop, rideDroop}
)

// rideDroop is how the seat's foot stands from level.
const rideDroop = 0.08

var mountKeys = mountPath{
	{at: 0, root: [3]float64{-MountAside, 0, 0}, feet: [2][3]float64{standSwing, standNear}, droop: flat, upper: Stand.Poses[0]},
	{at: 0.08, root: [3]float64{-0.72, 0, -0.06}, turn: 0.7, lean: 0.02, feet: [2][3]float64{{-0.62, 0.15, -0.18}, {-0.80, 0.091, -0.02}}, droop: inAir, upper: reachUp},
	{at: 0.16, root: [3]float64{-0.76, -0.01, -0.125}, turn: 1.35, lean: 0.05, feet: [2][3]float64{backOnFloor, headOnFloor}, droop: flat, upper: hold[0]},
	{at: 0.28, root: [3]float64{-0.78, -0.02, -0.125}, turn: 1.35, lean: 0.14, feet: [2][3]float64{backOnFloor, nearStirrup}, droop: [2]float64{0, -0.1}, upper: hold[1]},
	{at: 0.34, root: [3]float64{-0.74, -0.06, -0.125}, turn: 1.35, lean: 0.2, feet: [2][3]float64{backOnFloor, nearStirrup}, droop: [2]float64{0, -0.1}, upper: hold[2]},
	{at: 0.40, root: [3]float64{-0.66, 0.12, -0.12}, turn: 1.32, lean: 0.28, feet: [2][3]float64{{-0.66, 0.26, -0.34}, nearStirrup}, droop: [2]float64{0.45, -0.1}, upper: hold[3]},
	{at: 0.64, root: [3]float64{-0.37, 0.68, -0.12}, turn: 1.25, lean: 0.58, feet: [2][3]float64{{-0.46, 0.82, -0.33}, nearStirrup}, droop: inAir, upper: hold[4]},
	{at: 0.76, root: [3]float64{-0.22, 0.60, -0.14}, turn: 0.45, lean: 0.5, feet: [2][3]float64{overCroup, nearStirrup}, droop: inAir, upper: hold[5]},
	{at: 0.87, root: [3]float64{-0.08, 0.50, -0.14}, turn: 0.15, lean: 0.35, feet: [2][3]float64{{0.37, 1.05, -0.45}, nearStirrup}, droop: inAir, upper: hold[6]},
	{at: 1, root: [3]float64{0, rideBob, -RideBack}, lean: 0.08, feet: [2][3]float64{farStirrup, nearStirrup}, droop: seated, upper: ridePose(0)},
}

// Getting down: the far leg up and over the withers, the body leaning back and
// turning with it; sitting sideways on the near side, both legs hanging there,
// the hands on the saddle; a push off it, out and down through the air; landing
// on bent knees, and standing, facing away from the horse.
var dismountKeys = mountPath{
	{at: 0, root: [3]float64{0, rideBob, -RideBack}, lean: 0.08, feet: [2][3]float64{farStirrup, nearStirrup}, droop: seated, upper: ridePose(0)},
	{at: 0.16, root: [3]float64{0, 0.47, -0.16}, turn: -0.2, lean: -0.1, feet: [2][3]float64{{0.40, 1.12, 0.34}, {-0.30, 0.80, -0.08}}, droop: inAir, upper: hop[0]},
	{at: 0.34, root: [3]float64{-0.03, 0.50, -0.14}, turn: -0.35, lean: -0.25, feet: [2][3]float64{{-0.02, 1.60, 0.68}, {-0.32, 0.82, -0.08}}, droop: inAir, upper: hop[1]},
	{at: 0.52, root: [3]float64{-0.12, 0.42, -0.14}, turn: -1.35, lean: 0.05, feet: [2][3]float64{{-0.52, 0.92, 0.10}, {-0.52, 0.92, -0.25}}, droop: bothInAir, upper: hop[2]},
	{at: 0.62, root: [3]float64{-0.16, 0.40, -0.14}, turn: -1.4, lean: 0.22, feet: [2][3]float64{{-0.56, 0.90, 0.10}, {-0.56, 0.90, -0.25}}, droop: bothInAir, upper: hop[3]},
	{at: 0.74, root: [3]float64{-0.50, 0.34, -0.12}, turn: -1.45, lean: 0.12, feet: [2][3]float64{{-0.62, 0.46, 0.10}, {-0.62, 0.46, -0.24}}, droop: bothInAir, upper: hop[4]},
	{at: 0.86, root: [3]float64{-0.66, -0.18, 0}, turn: -1.5, lean: 0.3, feet: [2][3]float64{{-0.74, 0.091, 0.14}, {-0.74, 0.091, -0.14}}, droop: flat, upper: hop[5]},
	{at: 1, root: [3]float64{-MountAside + rootBehind, 0, rootBehind}, turn: -math.Pi / 2, feet: [2][3]float64{{-MountAside, 0.091, 0.142}, {-MountAside, 0.091, -0.142}}, droop: flat, upper: Stand.Poses[0]},
}
