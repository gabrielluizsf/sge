package human

// SleepLeft and SleepRight are the same sleep in a double bed, each on
// one side of it: the bed is placed under the sleeper by the motion's name,
// so a couple lie side by side, each in a picture of their own.
var SleepLeft, SleepRight = sleepIn("sleep_l"), sleepIn("sleep_r")

// SleepingTogether are the motions of one asleep in a double bed.
var SleepingTogether = []Motion{SleepLeft, SleepRight}

// sleepIn is Sleep under another name.
func sleepIn(name string) Motion {
	m := Sleep
	m.Name = name
	return m
}
