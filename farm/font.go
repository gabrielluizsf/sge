package farm

import (
	_ "embed"

	"github.com/gabrielluizsf/antui/canvas"
)

// fontData is the HUD's typeface: Sniglet ExtraBold, round and friendly,
// a static TrueType font from google/fonts under the SIL Open Font License
// (fonts/OFL.txt). It is the only font the game reads, never one from the
// player or the network.
//
//go:embed fonts/Sniglet-ExtraBold.ttf
var fontData []byte

// fontPixels is the size the HUD's text is drawn at, in screen pixels:
// the font is drawn at the screen's own resolution, never scaled up.
const fontPixels = 18

// hudFace is the HUD's face, read once.
var hudFace = func() *canvas.Face {
	f, err := canvas.ParseFace(fontData, fontPixels)
	if err != nil {
		panic(err)
	}
	return f
}()
