package main

import "github.com/gabrielluizsf/sge/daytime"

// A villager's trade, which decides what they do by day.
const (
	roleFarmer     = "farmer"
	roleLumberjack = "lumberjack"
	roleSmith      = "smith"
	roleWanderer   = "wanderer"
)

// What a schedule says to do.
const (
	doWork   = "work"
	doStroll = "stroll"
	doSleep  = "sleep"
)

// schedules are each trade's day: farmers rise first, the smith keeps the
// forge late, and the wanderers stroll all day; everyone sleeps at home.
var schedules = map[string]daytime.Schedule{
	roleFarmer:     {{From: 6, To: 18, Do: doWork}, {From: 18, To: 21, Do: doStroll}, {From: 21, To: 6, Do: doSleep}},
	roleLumberjack: {{From: 7, To: 17, Do: doWork}, {From: 17, To: 21, Do: doStroll}, {From: 21, To: 7, Do: doSleep}},
	roleSmith:      {{From: 7, To: 19, Do: doWork}, {From: 19, To: 21, Do: doStroll}, {From: 21, To: 7, Do: doSleep}},
	roleWanderer:   {{From: 8, To: 20, Do: doStroll}, {From: 20, To: 8, Do: doSleep}},
}

// roleLabels are the trades' names on screen.
var roleLabels = map[string]string{roleFarmer: "FAZENDEIRA", roleLumberjack: "MADEIREIRA", roleSmith: "FERREIRA", roleWanderer: "MORADORA"}

// roleLabel is v's trade on screen, as a man's or a woman's.
func roleLabel(v *villager) string {
	l := roleLabels[v.role]
	if !v.woman {
		l = l[:len(l)-1] + "O"
	}
	return l
}
