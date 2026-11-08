package human

// reachUp is the arms on their way up to the horse.
var reachUp = Pose{Nod: 0.1, Shoulder: both(0.6), Raise: both(0.12), Elbow: both(0.8)}

// hold are the arms at each key of getting on in which the hands are on the
// horse.
var hold = [7]Pose{
	{Nod: 0.15, Shoulder: [2]float64{0.89, 1.63}, Raise: [2]float64{0.70, 0.26}, Plane: [2]float64{2.12, 0.00}, Elbow: [2]float64{0.40, 0.07}},
	{Nod: 0.15, Shoulder: [2]float64{0.97, 1.73}, Raise: [2]float64{0.70, 0.26}, Plane: [2]float64{1.94, 0.00}, Elbow: [2]float64{0.62, 0.11}},
	{Nod: 0.15, Shoulder: [2]float64{1.08, 1.64}, Raise: [2]float64{0.71, 0.31}, Plane: [2]float64{1.76, -0.02}, Elbow: [2]float64{0.97, 0.53}},
	{Nod: 0.15, Shoulder: [2]float64{1.43, 0.91}, Raise: [2]float64{-0.10, 0.49}, Plane: [2]float64{1.60, 0.00}, Elbow: [2]float64{1.33, 1.19}},
	{Nod: 0.15, Shoulder: [2]float64{-0.07, 0.76}, Raise: [2]float64{0.67, 0.31}, Plane: [2]float64{2.30, 0.00}, Elbow: [2]float64{0.08, 0.05}},
	{Nod: 0.15, Shoulder: [2]float64{0.06, 1.02}, Raise: [2]float64{0.52, -0.09}, Plane: [2]float64{2.30, 0.00}, Elbow: [2]float64{0.14, 0.06}},
	{Nod: 0.15, Shoulder: [2]float64{-0.13, 0.90}, Raise: [2]float64{0.67, -0.10}, Plane: [2]float64{2.30, 0.33}, Elbow: [2]float64{0.25, 0.28}},
}

// hop are the arms of getting down: on the reins, thrown back for balance as the
// leg comes over, on the saddle either side of the hips, pushing off it, and out
// to land.
var hop = [6]Pose{
	{Nod: 0.1, Shoulder: both(0.2), Raise: both(0.15), Plane: both(0.25), Elbow: both(1.1)},
	{Nod: 0.05, Shoulder: both(-0.1), Raise: both(0.4), Plane: both(0.2), Elbow: both(0.5)},
	{Nod: 0.15, Shoulder: [2]float64{-0.10, -0.30}, Raise: [2]float64{0.22, 0.06}, Plane: [2]float64{-0.11, 0.00}, Elbow: [2]float64{0.34, 0.22}},
	{Nod: 0.15, Shoulder: [2]float64{-0.31, -0.40}, Raise: [2]float64{0.30, 0.00}, Plane: [2]float64{0.09, -0.24}, Elbow: [2]float64{0.49, 0.50}},
	{Nod: 0.2, Shoulder: both(0.4), Raise: both(0.35), Elbow: both(0.5)},
	{Nod: 0.3, Shoulder: both(0.8), Raise: both(0.25), Elbow: both(0.6)},
}
