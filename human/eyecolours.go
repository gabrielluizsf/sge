package human

// eyeColour is a material of the eyes folder and how common it is among
// those of African, Asian and European ancestry.
type eyeColour struct {
	item                     string
	african, asian, european float64
}

// eyeColours are the pack's irises: brown above all but among Europeans,
// who have the blues, greens and greys.
var eyeColours = []eyeColour{
	{"eyes/materials/brown.mhmat", 1, 1, 0.45},
	{"eyes/materials/brownlight.mhmat", 0.3, 0.3, 0.15},
	{"eyes/materials/bluegreen.mhmat", 0, 0, 0.06},
	{"eyes/materials/blue.mhmat", 0, 0, 0.15},
	{"eyes/materials/deepblue.mhmat", 0, 0, 0.08},
	{"eyes/materials/lightblue.mhmat", 0, 0, 0.10},
	{"eyes/materials/green.mhmat", 0, 0, 0.10},
	{"eyes/materials/grey.mhmat", 0, 0, 0.08},
	{"eyes/materials/ice.mhmat", 0, 0, 0.03},
}
