package farm

// stageOf is the stage, from 0 to ripe, that crop c shows after days days
// of watering: the seed under its mound before the first, the stages
// between spread over the days it takes, ripe once it has had them all.
func stageOf(c crop, days int) int {
	switch {
	case days <= 0:
		return 0
	case days >= c.daysToRipe:
		return ripe
	}
	return 1 + (days-1)*(ripe-1)/(c.daysToRipe-1)
}
