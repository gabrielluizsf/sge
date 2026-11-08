package human

import "github.com/gabrielluizsf/sge/random"

// AgeFor is the Body Age of years, the inverse of Body.Years.
func AgeFor(years float64) float64 {
	years = min(max(years, 1), 90)
	if years < 25 {
		return (years - 1) / 48
	}
	return 0.5 + (years-25)/130
}

// yearsOf is the age asked in years, or a grown person's, 18 to 73,
// likelier young.
func yearsOf(years float64, r *random.Source) float64 {
	if years > 0 {
		return min(max(years, 1), 90)
	}
	return 18 + (r.Float()+r.Float())/2*55
}

// bodyOf is a natural body of sex, years old and of ancestry a: muscle,
// weight, height and proportions about what such a person has.
func bodyOf(sex Sex, years float64, a Ancestry, r *random.Source) Body {
	b := Body{Gender: 0, Age: AgeFor(years), African: a.African, Asian: a.Asian, European: a.European}
	muscle, height := 0.38, 0.45
	if sex == Man {
		b.Gender, muscle, height = 1, 0.55, 0.55
	}
	if years < 16 {
		muscle = 0.3
	}
	b.Muscle = spreadOf(r, muscle, 0.3)
	b.Weight = spreadOf(r, 0.5, 0.35)
	b.Height = spreadOf(r, height, 0.4)
	b.Proportions = spreadOf(r, 0.5, 0.4)
	return b
}
