package human

import "strings"

// Body is the shape of a person, each from 0 to 1 as MakeHuman's sliders:
// Gender from woman to man; Age from a baby of 1 through 25 at 0.5 to 90;
// African, Asian and European ancestry, shares of one; Muscle, Weight,
// Height and Proportions, from the least to the most, 0.5 the average.
type Body struct {
	Gender, Age              float64
	African, Asian, European float64
	Muscle, Weight, Height   float64
	Proportions              float64
}

// Average is MakeHuman's own average person: between woman and man, 25,
// of equal ancestry, average in all.
var Average = Body{Gender: 0.5, Age: 0.5, African: 1.0 / 3, Asian: 1.0 / 3, European: 1.0 / 3,
	Muscle: 0.5, Weight: 0.5, Height: 0.5, Proportions: 0.5}

// Years is the body's age in years.
func (b Body) Years() float64 {
	if b.Age < 0.5 {
		return 1 + 48*b.Age
	}
	return 25 + 130*(b.Age-0.5)
}

// values are how much of each word of a macro target's name the body has,
// as MakeHuman works them out: a target is as strong as all its words.
func (b Body) values() map[string]float64 {
	c := func(v float64) float64 { return min(max(v, 0), 1) }
	v := map[string]float64{"universal": 1, "female": 1 - c(b.Gender), "male": c(b.Gender)}
	age := c(b.Age)
	if age < 0.5 {
		v["young"] = max(0, (age-0.1875)*3.2)
		v["baby"] = max(0, 1-age*5.333)
		v["child"] = max(0, min(1, 5.333*age)-v["young"])
	} else {
		v["old"] = age*2 - 1
		v["young"] = 1 - v["old"]
	}
	three(v, "muscle", c(b.Muscle))
	three(v, "weight", c(b.Weight))
	three(v, "height", c(b.Height))
	v["idealproportions"] = max(0, c(b.Proportions)*2-1)
	v["uncommonproportions"] = max(0, 1-c(b.Proportions)*2)
	sum := c(b.African) + c(b.Asian) + c(b.European)
	if sum == 0 {
		sum, b.African, b.Asian, b.European = 1, 1.0/3, 1.0/3, 1.0/3
	}
	v["african"], v["asian"], v["caucasian"] = c(b.African)/sum, c(b.Asian)/sum, c(b.European)/sum
	return v
}

// three sets the min, average and max values of a slider at s.
func three(v map[string]float64, name string, s float64) {
	v["max"+name] = max(0, s*2-1)
	v["min"+name] = max(0, 1-s*2)
	v["average"+name] = 1 - v["max"+name] - v["min"+name]
}

// weight is how strongly a target of this name, its words joined by "-",
// moves the body.
func weight(name string, v map[string]float64) float64 {
	w := 1.0
	for _, word := range strings.Split(name, "-") {
		w *= v[word]
	}
	return w
}
