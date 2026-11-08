package human

// hairShade is a colour of hair and how common it is among those of
// African, Asian and European ancestry.
type hairShade struct {
	colour                   Colour
	african, asian, european float64
}

// hairShades are black to blond and red.
var hairShades = []hairShade{
	{RGB(0.04, 0.035, 0.035), 0.9, 0.85, 0.12},
	{RGB(0.14, 0.08, 0.05), 0.1, 0.13, 0.25},
	{RGB(0.28, 0.17, 0.09), 0, 0.02, 0.26},
	{RGB(0.42, 0.20, 0.10), 0, 0, 0.08},
	{RGB(0.55, 0.20, 0.08), 0, 0, 0.05},
	{RGB(0.55, 0.40, 0.20), 0, 0, 0.14},
	{RGB(0.80, 0.65, 0.38), 0, 0, 0.14},
}
