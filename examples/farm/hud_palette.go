package main

import "github.com/gabrielluizsf/antui/canvas"

// The HUD's palette, warm as the farm's: dark cocoa text on the cream
// panels, and over the world cream, gold, coral and leaf text, each
// outlined in the farm's outline colour so it reads on grass and water.
var (
	inkOutline = canvas.RGB(0x3a, 0x2a, 0x2e)
	inks       = map[rune]canvas.Color{
		'k': canvas.RGB(0x4b, 0x2f, 0x27), // on a panel: cocoa
		'w': canvas.RGB(0xff, 0xf6, 0xe4), // cream
		'y': canvas.RGB(0xff, 0xd1, 0x66), // gold
		'r': canvas.RGB(0xff, 0x7a, 0x6e), // coral
		'g': canvas.RGB(0x9b, 0xe0, 0x7a), // leaf
	}
)

// outlined reports whether text in ink is drawn outlined: all but the
// panels' cocoa.
func outlined(ink rune) bool { return ink != 'k' }
