package farm

import (
	"fmt"
	"os"
	"strconv"
)

// Main plays the farm drawn by a as a command does, in a window titled
// title, 1280 by 720 pixels. FARM_DEMO=n in the environment plays n
// frames by itself instead, without a window, at a fixed step, and prints
// how the farm did, so a run in Go and one on Deno can be compared;
// FARM_OPEN=inventory, pause, credits, licence or talk opens that screen for the
// last frames; FARM_SHOT=file.png also saves the last frame, of a view
// FARM_VIEW=widthxheight in size, and FARM_AT=col,row stands the player
// on a cell first, to look at a place. Main exits on an error.
func Main(a Art, title string) {
	frames, _ := strconv.Atoi(os.Getenv("FARM_DEMO"))
	w, h := 1280, 720
	if frames <= 0 {
		exitOn(New(a, Keyboard()).Run(title, w, h))
		return
	}
	s := New(a, Bot(10))
	var col, row int
	if n, _ := fmt.Sscanf(os.Getenv("FARM_AT"), "%d,%d", &col, &row); n == 2 {
		s.Stand(col, row)
	}
	fmt.Sscanf(os.Getenv("FARM_VIEW"), "%dx%d", &w, &h)
	cv, err := s.Play(Demo{Frames: frames, Width: w, Height: h})
	exitOn(err)
	if open := os.Getenv("FARM_OPEN"); open != "" {
		cv, err = s.Open(open, Demo{Frames: 200, Width: w, Height: h})
		exitOn(err)
	}
	fmt.Println(s.Report())
	if shot := os.Getenv("FARM_SHOT"); shot != "" {
		exitOn(s.Shot(cv, shot))
	}
}

// exitOn prints err and exits, when there is one.
func exitOn(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
