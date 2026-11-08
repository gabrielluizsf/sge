// Command pack writes a game's installer: the game's installer program,
// built for the system it installs on, with the game's files after it.
//
//	go run github.com/gabrielluizsf/sge/installer/pack -stub setup.exe -files build/windows -o Game-Setup.exe
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/gabrielluizsf/sge/installer"
)

func main() {
	stub := flag.String("stub", "", "the game's installer program, built for the system it installs on")
	files := flag.String("files", "", "the folder of the game's files, its program among them")
	out := flag.String("o", "", "the installer to write")
	flag.Parse()
	if *stub == "" || *files == "" || *out == "" {
		flag.Usage()
		os.Exit(2)
	}
	if err := installer.Pack(*stub, *files, *out); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
