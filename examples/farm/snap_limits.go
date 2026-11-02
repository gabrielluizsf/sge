package farm

// The most a save may hold: a save over any of these is refused, so no
// number or list in a file can make the game huge or odd.
const (
	// maxSaveBytes is the largest save file read.
	maxSaveBytes = 4 << 20
	maxMoney     = 1_000_000_000
	maxCount     = 1_000_000
	maxDay       = 1_000_000
	maxYears     = 1_000_000
	maxCoord     = 1 << 24 // in pixels
	maxCell      = 1 << 19 // in cells
	maxPlots     = 20_000
	maxStumps    = 10_000
	maxPeople    = 5_000
	maxFamily    = 20_000
	maxTowns     = 500
	maxStacks    = 64
	maxText      = 48 // bytes of a name or id
	maxGene      = 63
	maxDays      = 1_000 // days a crop has been watered
	maxSeconds   = 1_000_000
)
