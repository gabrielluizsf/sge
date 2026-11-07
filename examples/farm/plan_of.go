package farm

import "github.com/gabrielluizsf/sge/life"

// planOf is v's plan for today, drawn once a day.
func (f *farm) planOf(v *villager) dayPlan {
	if v.planDay != f.clock.Day+1 {
		v.planDay = f.clock.Day + 1
		v.plan = f.planOn(v, f.clock.Day)
	}
	return v.plan
}

// planOn is v's plan for any day, from who they are and the day alone.
// Only a grown villager of a town has one.
func (f *farm) planOn(v *villager, day int) dayPlan {
	if v.ours || v.town == nil || v.stage() < life.Adult {
		return nil
	}
	trip := 0.0
	if r := f.roadFrom(v.town); r != nil {
		trip = r.hours
	}
	return planFor(planCall{id: v.id, role: v.role, day: day, trip: trip})
}
