package sound

// Seconds is how many samples, at Rate, last t seconds.
func Seconds(t float64) int { return int(t * Rate) }
