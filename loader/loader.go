// Package loader loads a game's heavy assets without stopping it: jobs
// queue by priority and run a few milliseconds a frame, within a budget,
// so the frame rate holds while pictures decode; what they made is kept
// up to a memory budget, the least lately used let go past it and made
// again when asked for; and jobs belong to groups whose progress a
// loading screen shows.
package loader

import "time"

// Job makes an asset: its name, the group it belongs to, how much memory
// it holds once made, roughly, in bytes, and how to make it.
type Job struct {
	Name  string
	Group string
	Bytes int64
	Make  func() (any, error)
}

// Loader runs jobs as a game's script: Budget is the time a frame may
// spend on them (4ms when 0), Memory how much what they made may hold
// before the least lately used is let go (no limit when 0).
type Loader struct {
	Budget time.Duration
	Memory int64

	queue  []queued
	assets map[string]*asset
	groups map[string]*progress
	used   int64
	clock  int64 // ticks of use, for the least lately used
	seq    int64
}

// queued is a job waiting, and how urgent it is: lower first, then in the
// order asked.
type queued struct {
	job      Job
	priority int
	seq      int64
}

// asset is a job's result, and when it was last asked for.
type asset struct {
	value any
	err   error
	bytes int64
	used  int64
	group string
}

// progress is a group's jobs asked for, and how many are done.
type progress struct {
	asked, done int
}
