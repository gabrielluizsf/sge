package main

import (
	_ "embed"

	"github.com/gabrielluizsf/antui/canvas"
)

// fontData is the realistic farm's typeface: Source Serif 4 SemiBold, a
// static TrueType font from Adobe's official release
// (github.com/adobe-fonts/source-serif) under the SIL Open Font License
// 1.1 (fonts/OFL.txt), made for legibility in interfaces. It is the only
// font the game reads, never one from the player or the network.
//
//go:embed fonts/SourceSerif4-Semibold.ttf
var fontData []byte

// fontPixels is the size the HUD's text is drawn at, in screen pixels.
const fontPixels = 18

// hudFace is the HUD's face, read once.
var hudFace = func() *canvas.Face {
	f, err := canvas.ParseFace(fontData, fontPixels)
	if err != nil {
		panic(err)
	}
	return f
}()
