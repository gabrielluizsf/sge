// Command farmart writes the farm's sprites as SVG files, and a preview
// of them as a PNG: every sprite enlarged, the animations frame by frame,
// a small map made of the tiles to see that they meet, and the HUD.
//
//	go run ./examples/farm/cmd/farmart -assets examples/farm/assets -preview examples/farm/preview.png
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/gabrielluizsf/sge/examples/farm/art"
)

func main() {
	assets := flag.String("assets", "assets", "directory to write the SVG files to")
	preview := flag.String("preview", "preview.png", "PNG to write the preview to")
	flag.Parse()
	if err := os.MkdirAll(*assets, 0o755); err != nil {
		fail(err)
	}
	sprites := art.Sprites()
	names := make([]string, 0, len(sprites))
	for name := range sprites {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		if err := os.WriteFile(filepath.Join(*assets, name+".svg"), []byte(art.SVG(sprites[name])), 0o644); err != nil {
			fail(err)
		}
	}
	cv, err := previewSheet(*assets)
	if err != nil {
		fail(err)
	}
	if err := writePNG(cv, *preview); err != nil {
		fail(err)
	}
	fmt.Printf("%d sprites, preview %dx%d\n", len(names), cv.Width, cv.Height)
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, "farmart:", err)
	os.Exit(1)
}
