package farm

// runSpeed is how fast the farmer runs on land, holding Shift: nearly
// twice as fast as they walk.
const runSpeed = walkSpeed * 1.8

// speedOf is how fast a farmer moves: running on land if they mean to,
// walking otherwise, wading in shallow water and swimming in deep, where
// there is no running.
func speedOf(run, wet, swims bool) float64 {
	switch {
	case swims:
		return swimSpeed
	case wet:
		return wadeSpeed
	case run:
		return runSpeed
	}
	return walkSpeed
}

// moveClip is the suffix of the clip a farmer is shown by as they move:
// swimming, running or walking.
func moveClip(running, swims bool) string {
	switch {
	case swims:
		return "_swim"
	case running:
		return "_run"
	}
	return "_walk"
}

// running is whether the farmer, told to as in says, runs: moving, on
// land, with their will to.
func running(in intent, wet, swims bool) bool {
	return in.run && in.dir.Length() > 0 && !wet && !swims
}
