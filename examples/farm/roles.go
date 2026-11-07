package farm

import (
	"github.com/gabrielluizsf/sge/daytime"
	"github.com/gabrielluizsf/sge/examples/i18n"
)

// A villager's trade, which decides what they do by day.
const (
	roleFarmer     = "farmer"
	roleLumberjack = "lumberjack"
	roleSmith      = "smith"
	roleWanderer   = "wanderer"
	roleDealer     = "dealer"
)

// What a schedule says to do.
const (
	doWork   = "work"
	doStroll = "stroll"
	doSleep  = "sleep"
)

// bedBy is the hour the villagers are in bed by: they start home half an
// hour before, to be in by ten.
const bedBy = 21.5

// schedules are each trade's day: farmers rise first, the smith keeps the
// forge late, and the wanderers stroll all day; everyone sleeps at home.
var schedules = map[string]daytime.Schedule{
	roleFarmer:     {{From: 6, To: 18, Do: doWork}, {From: 18, To: bedBy, Do: doStroll}, {From: bedBy, To: 6, Do: doSleep}},
	roleLumberjack: {{From: 7, To: 17, Do: doWork}, {From: 17, To: bedBy, Do: doStroll}, {From: bedBy, To: 7, Do: doSleep}},
	roleSmith:      {{From: 7, To: 19, Do: doWork}, {From: 19, To: bedBy, Do: doStroll}, {From: bedBy, To: 7, Do: doSleep}},
	roleWanderer:   {{From: 8, To: bedBy, Do: doWork}, {From: bedBy, To: 8, Do: doSleep}},
	roleDealer:     {{From: 7, To: 18, Do: doWork}, {From: 18, To: bedBy, Do: doStroll}, {From: bedBy, To: 7, Do: doSleep}},
}

// roleLabels are the trades' names on screen, as English keys, a woman's
// and a man's.
var roleLabels = map[string][2]string{
	roleFarmer:     {"farmer", "farmer"},
	roleLumberjack: {"lumberjack", "lumberjack"},
	roleSmith:      {"smith", "smith"},
	roleWanderer:   {"wanderer", "wanderer"},
	roleDealer:     {"livestock dealer", "livestock dealer"},
}

// roleLabel is v's trade on screen, as a woman's or a man's.
func roleLabel(v *villager) string {
	l := roleLabels[v.role]
	return i18n.TF(v.woman, l[0])
}
