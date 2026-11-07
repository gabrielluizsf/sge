package main

import (
	"fmt"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/label"
)

// The HUD's ink on the dark panels and over the world: warm parchment on
// a panel, and over the world cream, gold, coral and leaf, outlined in
// near black so they read on any ground.
var (
	hudOutline = canvas.RGB(0x14, 0x10, 0x0f)
	hudInks    = map[rune]canvas.Color{
		'k': canvas.RGB(0xf3, 0xe9, 0xd2), // on a panel
		'w': canvas.RGB(0xff, 0xf6, 0xe4),
		'y': canvas.RGB(0xf6, 0xca, 0x62),
		'r': canvas.RGB(0xff, 0x8c, 0x7c),
		'g': canvas.RGB(0xa8, 0xe8, 0x8a),
	}
)

// hudOutlineWidth is how many pixels round text over the world is
// outlined.
const hudOutlineWidth = 2

// Text is s in Source Serif, in the palette's ink: outlined round, but
// for what is on a panel.
func (a *realArt) Text(s string, ink rune) *canvas.Canvas {
	key := fmt.Sprint("t", string(ink), s)
	if cv, ok := a.hud[key]; ok {
		return cv
	}
	cv := serifText(s, ink)
	a.hud[key] = cv
	return cv
}

// serifText is s in Source Serif, in the palette's ink, uncached: outlined
// round, but for what is on a panel.
func serifText(s string, ink rune) *canvas.Canvas {
	st := label.Style{Ink: hudInks[ink]}
	if ink != 'k' {
		st.Outline, st.Width = hudOutline, hudOutlineWidth
	}
	return label.New(hudFace, s, st)
}
