package installer

import "io/fs"

// Setup is one installing: App installed from Files onto Place, with
// Uninstaller (a program, or nil) written as the one that takes it away;
// Steps run after the files are in, the game's own, such as baking its
// pictures, into Cache: a folder of the game's, marked as the app's and
// taken away with it (never its saves). Limits bound the files; Run runs
// the system's own programs (nil runs them for real).
type Setup struct {
	App         App
	Files       fs.FS
	Place       Place
	Uninstaller []byte
	Steps       []Step
	Cache       string
	Limits      Limits
	Run         Runner
}

// Step is one of the game's own steps in installing, named for the
// window's line: Run is given the folder the game's files are in.
type Step struct {
	Name string
	Run  func(files string) error
}

// Runner runs program with args, env added to the environment.
type Runner func(program string, args, env []string) error

func (s Setup) runner() Runner {
	if s.Run != nil {
		return s.Run
	}
	return runSystem
}
