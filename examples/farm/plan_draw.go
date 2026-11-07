package farm

import "github.com/gabrielluizsf/sge/random"

// planCall is what a day's plan is drawn from: who, of what trade, which
// day, and how many hours it takes to walk to the next town, none when
// there is no road.
type planCall struct {
	id, role string
	day      int
	trip     float64
}

// weights are how often each trade turns to each chore, in the order of
// the chores.
var weights = map[string][choreCount]int{
	roleFarmer:     {2, 1, 6, 0, 2, 0},
	roleLumberjack: {6, 2, 1, 0, 2, 0},
	roleSmith:      {1, 3, 1, 6, 2, 0},
	roleWanderer:   {3, 3, 3, 0, 4, 0},
	roleDealer:     {0, 0, 0, 0, 2, 0, 8},
}

// visitChance is how often each trade spends a day's stretch in the next
// town.
var visitChance = map[string]float64{roleFarmer: 0.3, roleLumberjack: 0.25, roleSmith: 0.15, roleWanderer: 0.5, roleDealer: 0.1}

// A visit lasts at least minStay hours there, and is over by backBy, so
// they are home before bedBy.
const (
	minStay = 2.0
	backBy  = bedBy - 1
)

// mix is the seed of a villager's day: the world's, theirs and the day's.
func mix(id string, day int) uint32 {
	h := uint32(2166136261) ^ uint32(seed)
	for i := 0; i < len(id); i++ {
		h ^= uint32(id[i])
		h *= 16777619
	}
	h ^= uint32(day+1) * 2654435761
	return h*2246822519 + 3266489917
}

// workHours are the hours a trade works, as its schedule says.
func workHours(role string) (from, to float64) {
	for _, e := range schedules[role] {
		if e.Do == doWork {
			return e.From, e.To
		}
	}
	return 0, 0
}

// planFor is the day c asks for: the hours of work cut into stretches,
// each on a chore their trade favours, and now and then a trip to the next
// town in the middle. The same call gives the same plan, in any order.
func planFor(c planCall) dayPlan {
	src := random.New(mix(c.id, c.day))
	from, to := workHours(c.role)
	if to <= from {
		return nil
	}
	if c.trip > 0 && src.Float() < visitChance[c.role] {
		length := 2*c.trip + minStay + 3*src.Float()
		if latest := backBy - length; latest >= from+1 {
			s := from + 1 + (latest-from-1)*src.Float()
			out := fill(src, c.role, spell{from: from, to: s, doing: choreCount})
			out = append(out, spell{from: s, to: s + length, doing: choreVisit})
			return append(out, fill(src, c.role, spell{from: s + length, to: to})...)
		}
	}
	return fill(src, c.role, spell{from: from, to: to})
}

// maxStretch is the longest a villager keeps to one chore, in hours.
const maxStretch = 3.0

// fill cuts the hours of in into stretches of a chore each, from the first
// hour to the last, none the same as the one before; none if they come
// to less than an hour.
func fill(src *random.Source, role string, in spell) dayPlan {
	span := in.to - in.from
	if span < 1 {
		return nil
	}
	n := int(span/maxStretch) + 1
	var out dayPlan
	prev := choreCount
	for i := 0; i < n; i++ {
		c := pickChore(src, role, prev)
		from, to := in.from+span*float64(i)/float64(n), in.from+span*float64(i+1)/float64(n)
		if i == 0 {
			from = in.from
		}
		if i == n-1 {
			to = in.to
		}
		if i > 0 {
			from = out[i-1].to
		}
		out = append(out, spell{from: from, to: to, doing: c})
		prev = c
	}
	return out
}

// pickChore is a chore of role by their weights, other than not.
func pickChore(src *random.Source, role string, not chore) chore {
	w, ok := weights[role]
	if !ok {
		w = weights[roleWanderer]
	}
	total := 0
	for c := chore(0); c < choreCount; c++ {
		if c != not {
			total += w[c]
		}
	}
	n := src.Intn(total)
	for c := chore(0); c < choreCount; c++ {
		if c == not {
			continue
		}
		if n < w[c] {
			return c
		}
		n -= w[c]
	}
	return choreTrade
}
