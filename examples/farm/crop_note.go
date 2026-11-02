package farm

import "fmt"

// cropNote is the short line in Portuguese that tells how p's crop is
// doing: how long until its next stage, or that it is ripe or withered,
// and whether it is watered today.
func cropNote(p *plot) string {
	c := crops[p.crop]
	name := upperFirst(c.label)
	switch {
	case p.dead:
		return name + " - murcho"
	case p.stage >= ripe:
		return name + " - pronto para colher"
	}
	verb := "cresce"
	if p.stage == 0 {
		verb = "brota"
	}
	return fmt.Sprint(name, " - ", verb, " em ", daysLabel(daysToNextStage(c, p)), " - ", wateredLabel(p))
}

// daysToNextStage is how many more days of watering p's crop needs to
// show its next stage.
func daysToNextStage(c crop, p *plot) int {
	d := p.days + 1
	for stageOf(c, d) <= p.stage && d < c.daysToRipe {
		d++
	}
	return d - p.days
}

// daysLabel is n days in words.
func daysLabel(n int) string {
	if n == 1 {
		return "1 dia"
	}
	return fmt.Sprint(n, " dias")
}

// wateredLabel says whether p has been watered today.
func wateredLabel(p *plot) string {
	if p.watered {
		return "regado"
	}
	return "precisa de água"
}

// upperFirst is s with its first letter capitalised, for the plain
// ASCII and accented lower-case initials crop names have.
func upperFirst(s string) string {
	r := []rune(s)
	if len(r) > 0 && r[0] >= 'a' && r[0] <= 'z' {
		r[0] -= 'a' - 'A'
	}
	return string(r)
}
