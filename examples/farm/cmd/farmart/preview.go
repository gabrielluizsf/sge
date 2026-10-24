package main

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/examples/farm/art"
)

const scale = 5

// rows are the preview's lines of sprites, each an animation or a group.
var rows = [][]string{
	{"grass0", "grass1", "grass2", "grass3", "flowers0", "flowers1", "flowers2", "flowers3"},
	{"water0", "water1", "water2", "water3", "path", "soil_dry", "soil_wet", "fence"},
	{"stall_roof_l", "stall_roof_r", "tree_crown0", "tree_crown1", "tree_crown2", "tree_crown3", "tree_trunk"},
	{"stall_counter_l", "stall_counter_r"},
	{"carrot0", "carrot1", "carrot2", "carrot3", "carrot4", "tomato0", "tomato1", "tomato2", "tomato3", "tomato4"},
	{"farmer_down", "farmer_down1", "farmer_down", "farmer_down2", "farmer_downAct", "farmer_up", "farmer_up1", "farmer_up", "farmer_up2", "farmer_upAct"},
	{"farmer_right", "farmer_right1", "farmer_right", "farmer_right2", "farmer_rightAct", "farmer_left", "farmer_left1", "farmer_left", "farmer_left2", "farmer_leftAct"},
	{"icon_coin", "icon_carrotSeeds", "icon_tomatoSeeds", "icon_can", "icon_carrot", "icon_tomato"},
}

// testMap is a small farm laid out in tiles, to see them meet.
var testMap = []string{
	"ffffffffffff",
	"gGgggRrggTgg",
	"gggggCcggtwW",
	"ppppppppppwW",
	"gdsdsdsgggwW",
	"gdsdsdsgFggg",
}

func previewSheet(dir string) (*canvas.Canvas, error) {
	cell := art.Size*scale + 4
	cv, err := canvas.NewCanvas(10*cell+8, (len(rows)+len(testMap))*cell+40+len(lake)*art.Size*3+8)
	if err != nil {
		return nil, err
	}
	cv.Clear(canvas.RGB(40, 36, 48))
	for r, row := range rows {
		for c, name := range row {
			img, err := load(dir, name, scale)
			if err != nil {
				return nil, err
			}
			cv.Blit(4+c*cell, 4+r*cell, img)
		}
	}
	top := 8 + len(rows)*cell
	if err := drawTestMap(cv, dir, 4, top); err != nil {
		return nil, err
	}
	if err := drawLake(cv, dir, 4, top+len(testMap)*art.Size*3+8); err != nil {
		return nil, err
	}
	return cv, drawHUD(cv, 4+12*art.Size*3+12, top)
}

// drawTestMap lays the test map out at three times the size, crops and the
// farmer on it, as the game will.
func drawTestMap(cv *canvas.Canvas, dir string, x0, y0 int) error {
	const s = 3
	tiles := map[rune]string{
		'g': "grass0", 'G': "flowers1", 'f': "fence", 'p': "path", 'd': "soil_dry", 's': "soil_wet",
		'w': "water1", 'W': "water2", 'R': "stall_roof_l", 'r': "stall_roof_r", 'C': "stall_counter_l", 'c': "stall_counter_r",
		'T': "grass0", 't': "grass0", 'F': "grass0",
	}
	extra := map[[2]int]string{{2, 4}: "carrot1", {4, 4}: "carrot4", {2, 5}: "tomato2", {6, 5}: "tomato4", {6, 4}: "carrot2", {8, 5}: "farmer_rightAct", {9, 1}: "tree_crown0", {9, 2}: "tree_trunk"}
	for r, line := range testMap {
		for c, ch := range line {
			// Objects are transparent around them: the grass goes under.
			for _, name := range []string{"grass0", tiles[ch]} {
				img, err := load(dir, name, s)
				if err != nil {
					return err
				}
				cv.Blit(x0+c*art.Size*s, y0+r*art.Size*s, img)
			}
			if e, ok := extra[[2]int{c, r}]; ok {
				top, err := load(dir, e, s)
				if err != nil {
					return err
				}
				cv.Blit(x0+c*art.Size*s, y0+r*art.Size*s, top)
			}
		}
	}
	return nil
}

// drawHUD draws the HUD's pieces: the money box, the slots with the one in
// use, and a line of text, at three times the size.
func drawHUD(cv *canvas.Canvas, x0, y0 int) error {
	const s = 3
	money, err := render(art.Panel(44, 14, false), s)
	if err != nil {
		return err
	}
	cv.Blit(x0, y0, money)
	text, err := render(art.Text("$125", 'w'), s)
	if err != nil {
		return err
	}
	cv.Blit(x0+16*s, y0+3*s, text)
	coin, err := render(art.SVG(art.Icons["coin"]), s)
	if err != nil {
		return err
	}
	cv.Blit(x0, y0-1*s, coin)
	for i, icon := range []string{"carrotSeeds", "tomatoSeeds", "can", "carrot"} {
		slot, err := render(art.Panel(20, 20, i == 2), s)
		if err != nil {
			return err
		}
		cv.Blit(x0+i*22*s, y0+20*s, slot)
		ic, err := render(art.SVG(art.Icons[icon]), s)
		if err != nil {
			return err
		}
		cv.Blit(x0+i*22*s+2*s, y0+22*s, ic)
	}
	line, err := render(art.Text("E: LOJA  ESPACO: USAR", 'y'), s/3*2)
	if err != nil {
		return err
	}
	cv.Blit(x0, y0+46*s, line)
	return nil
}
