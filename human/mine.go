package human

// Mine is breaking rock with a pickaxe held in both hands: the same
// overhead swing as the felling cut, the point driven into the rock where
// the axe's edge bites the trunk, so the beats are the chop's.
var Mine = loop{"mine", 1.4, 40}.of(func(t float64) Pose { return poseAt(chopBeats, t) })

// Mining are the motions of those who break rock, besides Motions.
var Mining = []Motion{Mine}

// MineVersion changes whenever the swing does, or the pickaxe in the
// hands that swing it, so what was baked of it is baked again, and no
// other motion is.
const MineVersion = 1
