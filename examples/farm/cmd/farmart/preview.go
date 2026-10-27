package main

import (
	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/examples/farm/art"
)

const scale = 5

// rows are the preview's lines of sprites, each an animation or a group.
var rows = [][]string{
	{"grass0", "grass1", "grass2", "grass3", "flowers0", "flowers1", "flowers2", "flowers3"},
	{"path", "soil_dry", "soil_wet", "fence", "tree0", "tree1", "tree2", "tree3"},
	{"carrot0", "carrot1", "carrot2", "carrot3", "carrot4", "tomato0", "tomato1", "tomato2", "tomato3", "tomato4"},
	{"farmer_down", "farmer_down1", "farmer_down", "farmer_down2", "farmer_downAct", "farmer_up", "farmer_up1", "farmer_up", "farmer_up2", "farmer_upAct"},
	{"farmer_right", "farmer_right1", "farmer_right", "farmer_right2", "farmer_rightAct", "farmer_left", "farmer_left1", "farmer_left", "farmer_left2", "farmer_leftAct"},
	{"farmer_downSwim", "farmer_downSwim1", "farmer_upSwim", "farmer_upSwim1", "farmer_rightSwim", "farmer_rightSwim1", "farmer_leftSwim", "farmer_leftSwim1"},
	{"icon_coin", "icon_carrotSeeds", "icon_tomatoSeeds", "icon_can", "icon_carrot", "icon_tomato", "cursor", "stall"},
}

// testMap is a small farm laid out in tiles, to see them meet.
var testMap = []string{
	"ffffffffffff",
	"gGgggggggTgg",
	"gggggCcggTww",
	"ppppppppppww",
	"gdsdsdsgggww",
	"gdsdsdsgggGg",
}

func previewSheet(dir string) (*canvas.Canvas, error) {
	cell := art.IsoW*scale + 4
	rowH := art.IsoH*scale + 4
	mapH := (len(testMap)+12)*art.Size/2*3 + art.IsoH*3
	cv, err := canvas.NewCanvas(10*cell+8, len(rows)*rowH+2*mapH+8)
	if err != nil {
		return nil, err
	}
	cv.Clear(canvas.RGB(40, 36, 48))
	for r, row := range rows {
		x := 4
		for _, name := range row {
			img, err := load(dir, name, scale)
			if err != nil {
				return nil, err
			}
			cv.Blit(x, 4+r*rowH, img)
			x += img.Width + 4
		}
	}
	top := 8 + len(rows)*rowH
	if err := drawTestMap(cv, dir, 4+len(testMap)*art.Size*3, top); err != nil {
		return nil, err
	}
	if err := drawLake(cv, dir, 4+len(lake)*art.Size*3, top+mapH); err != nil {
		return nil, err
	}
	return cv, drawHUD(cv, 4+36*art.Size*3, top)
}

// isoBlit draws the tile img of the cell at, column and row, as the game
// does, at three times the size: its bottom middle on the near corner of
// the cell's diamond, the diamond of cell 0, 0 topped at origin.
func isoBlit(cv *canvas.Canvas, img *canvas.Canvas, at [2]int, origin [2]int) {
	half := art.Size * 3
	bx := origin[0] + (at[0]-at[1])*half
	by := origin[1] + (at[0]+at[1]+2)*half/2
	cv.Blit(bx-img.Width/2, by-img.Height, img)
}

// drawTestMap lays the test map out at three times the size, seen from
// the corner, crops and the farmer on it, as the game will.
func drawTestMap(cv *canvas.Canvas, dir string, x0, y0 int) error {
	const s = 3
	ground := map[rune]string{'G': "flowers1", 'p': "path", 'd': "soil_dry", 's': "soil_wet"}
	things := map[rune]string{'f': "fence", 'T': "tree0"}
	extra := map[[2]int]string{{2, 4}: "carrot1", {4, 4}: "carrot4", {2, 5}: "tomato2", {6, 5}: "tomato4", {6, 4}: "carrot2"}
	o := [2]int{x0, y0 + art.IsoH*s}
	for d := 0; d < len(testMap)+len(testMap[0]); d++ {
		for r, line := range testMap {
			c := d - r
			if c < 0 || c >= len(line) {
				continue
			}
			names := []string{"grass0"}
			if g, ok := ground[rune(line[c])]; ok {
				names[0] = g
			}
			if line[c] == 'w' {
				names = append(names, "shore15_1")
			}
			if t, ok := things[rune(line[c])]; ok {
				names = append(names, t)
			}
			if e, ok := extra[[2]int{c, r}]; ok {
				names = append(names, e)
			}
			for _, name := range names {
				img, err := load(dir, name, s)
				if err != nil {
					return err
				}
				isoBlit(cv, img, [2]int{c, r}, o)
			}
		}
	}
	stall, err := load(dir, "stall", s)
	if err != nil {
		return err
	}
	// The stall stands on cells 5 and 6 of row 2, its pivot on their middle.
	half := art.Size * 3
	cx, cy := o[0]+(6-2)*half, o[1]+(6+2+1)*half/2
	cv.Blit(cx-int(art.StallPivot[0]*float64(stall.Width)), cy-int(art.StallPivot[1]*float64(stall.Height)), stall)
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
	coin, err := render(art.Icons()["coin"], s)
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
		ic, err := render(art.Icons()[icon], s)
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
