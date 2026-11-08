// Package vitals is a person's health and food: food falls with time,
// health falls with blows and with hunger, and eating, rest and a full
// stomach mend it.
package vitals

// Max is what health and food are when full.
const Max = 100.0

// Rates, per second: how fast food falls, how fast health falls while
// the food is out, how fast it mends while the food is well above Fed.
const (
	Hunger = 0.055
	Starve = 0.15
	Mend   = 0.12
	// Fed is the food above which health mends of itself.
	Fed = 60.0
)

// Vitals is a person's health and food, each from 0 to Max.
type Vitals struct {
	Health, Food float64
}

// Full is a person at their best.
func Full() Vitals { return Vitals{Health: Max, Food: Max} }

// Dead reports whether health has run out.
func (v *Vitals) Dead() bool { return v.Health <= 0 }

// Tick passes dt seconds: food falls, and, with none, health falls with
// it; with plenty, health mends.
func (v *Vitals) Tick(dt float64) {
	v.Food = clamp(v.Food - Hunger*dt)
	switch {
	case v.Food <= 0:
		v.Health = clamp(v.Health - Starve*dt)
	case v.Food > Fed:
		v.Health = clamp(v.Health + Mend*dt)
	}
}

// Hurt takes n from health, and reports whether it ran out.
func (v *Vitals) Hurt(n float64) bool {
	v.Health = clamp(v.Health - n)
	return v.Dead()
}

// Heal gives n to health.
func (v *Vitals) Heal(n float64) { v.Health = clamp(v.Health + n) }

// Eat gives n to food.
func (v *Vitals) Eat(n float64) { v.Food = clamp(v.Food + n) }

// Valid reports whether both are numbers a person could have: from 0 to Max.
func (v Vitals) Valid() bool {
	return v.Health >= 0 && v.Health <= Max && v.Food >= 0 && v.Food <= Max
}

func clamp(x float64) float64 {
	if !(x > 0) {
		return 0
	}
	return min(x, Max)
}
