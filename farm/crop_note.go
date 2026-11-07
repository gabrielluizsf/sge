package farm

import "github.com/gabrielluizsf/sge/i18n"

// cropNote is the short line that tells how p's crop is doing: how long
// until its next stage, or that it is ripe or withered, and whether it is
// watered today.
func cropNote(p *plot) string {
	c := crops[p.crop]
	name := upperFirst(i18n.T(c.label)) + ownerLabel(p)
	switch {
	case p.dead:
		return i18n.T("%s - withered", name)
	case p.stage >= ripe:
		return i18n.T("%s - ready to harvest", name)
	}
	verb := "grows"
	if p.stage == 0 {
		verb = "sprouts"
	}
	return i18n.T("%s - %s in %s - %s", name, i18n.T(verb), daysLabel(daysToNextStage(c, p)), wateredLabel(p))
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
		return i18n.T("1 day")
	}
	return i18n.T("%d days", n)
}

// wateredLabel says whether p has been watered today.
func wateredLabel(p *plot) string {
	if p.watered {
		return i18n.T("watered")
	}
	return i18n.T("needs water")
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
