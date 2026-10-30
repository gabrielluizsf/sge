package farm

import "github.com/gabrielluizsf/sge/life"

// daysPerYear is how many of the game's days make a year: 28, four weeks,
// as farm games make a season. A child is born a year, lineage's
// Gestation, after their parents decide on them, and is grown 18 years,
// 504 days, on; one of life's spans, 68 to 94 years, is a few thousand
// days, and the player, a young adult of founderYears when the game
// starts, lives 48 years at least, so their children grow up and take
// over the farm. F9 passes a year at once, to see it happen.
const daysPerYear = 28

// calendar is how the farm's days make years.
var calendar = life.Calendar{DaysPerYear: daysPerYear}

// founderYears is how old the farmer who starts the game is.
const founderYears = 20

// fertileUntil is the age past which a woman has no more children.
const fertileUntil = 45

// yearsOf is how old one of l is on day.
func yearsOf(l life.Life, day int) float64 { return calendar.Age(l, day) }

// grownOn reports whether one of l is grown up on day.
func grownOn(l life.Life, day int) bool { return life.Grown(yearsOf(l, day)) }
