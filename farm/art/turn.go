package art

import "math"

// cosTurn and sinTurn are the cosine and sine of a turns of a circle.
func cosTurn(a float64) float64 { return math.Cos(2 * math.Pi * a) }
func sinTurn(a float64) float64 { return math.Sin(2 * math.Pi * a) }
