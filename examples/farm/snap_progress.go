package farm

import "github.com/gabrielluizsf/sge/progress"

// progressOf is the level and XP s holds; a save from before there were
// any holds none, and is level 1.
func progressOf(s *snapshot) progress.Progress {
	if s.Level == 0 && s.XP == 0 {
		return progress.New()
	}
	return progress.Progress{Level: s.Level, XP: s.XP}
}

// checkProgress is nil if s's level and XP are ones the game could have
// reached: a level from 1 up to the last, and XP within that level.
func checkProgress(s *snapshot) error {
	return progressOf(s).Valid()
}
