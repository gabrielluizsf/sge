package vitals

import (
	"math"
	"testing"
)

func TestFoodFallsWithTime(t *testing.T) {
	v := Full()
	v.Tick(100)
	if v.Food >= Max {
		t.Fatal("food did not fall")
	}
}

func TestHealthFallsWhenStarving(t *testing.T) {
	v := Vitals{Health: 50}
	v.Tick(10)
	if v.Health >= 50 {
		t.Fatal("health did not fall with no food")
	}
}

func TestHealthMendsWhenFed(t *testing.T) {
	v := Vitals{Health: 50, Food: 90}
	v.Tick(10)
	if v.Health <= 50 {
		t.Fatal("health did not mend")
	}
}

func TestHurtKillsAtZeroAndNeverGoesBelow(t *testing.T) {
	v := Full()
	if v.Hurt(30) || !v.Hurt(500) || v.Health != 0 {
		t.Fatalf("%+v", v)
	}
}

func TestEatAndHealStopAtMax(t *testing.T) {
	v := Vitals{Health: 90, Food: 90}
	v.Eat(50)
	v.Heal(50)
	if v != Full() {
		t.Fatalf("%+v", v)
	}
}

func TestValidRefusesWhatIsOutOfRange(t *testing.T) {
	for _, v := range []Vitals{{-1, 5}, {5, -1}, {101, 5}, {5, 101}, {math.NaN(), 5}, {5, math.NaN()}} {
		if v.Valid() {
			t.Fatalf("%+v valid", v)
		}
	}
	if !Full().Valid() {
		t.Fatal("full is not valid")
	}
}
