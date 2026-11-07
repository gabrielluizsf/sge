package bakes

import (
	"fmt"
	"image"
	"image/color"

	"github.com/gabrielluizsf/sge/bake"
	"github.com/gabrielluizsf/sge/examples/farm/kit"
	"github.com/gabrielluizsf/sge/gltf"
	"github.com/gabrielluizsf/sge/human"
)

// DoubleBedModel is the asset a couple sleep in: Poly Haven's Gothic Bed
// 01, CC0, by Kirill Sannikov.
const DoubleBedModel = "GothicBed_01"

// DoubleBedTurn is the heading of the double bed's own pictures, in quarter
// turns, that lies as the couple lie in it: head to the back wall, foot to the room.
const DoubleBedTurn = 3

// A couple lie side by side, each a doubleAside metres off the middle of
// the bed, across it, and baked in a picture of their own that keeps only
// their half of the bed, so the two pictures put together are one bed with
// two in it, and neither is inside the other or the frame.
var (
	doubleAt   = [3]float64{0, doubleDown, 0}
	doubleTurn = 0.0
)

// doubleAside is how many metres each of the couple is from the bed's
// middle, and doubleDown how many the bed is let down below their backs.
const (
	doubleAside = 0.37
	doubleDown  = 0.0
)

// layInDoubleBed gives m, a person already animated, the double bed to
// sleep in on each side, if they have the motions of one.
func (b *Baker) layInDoubleBed(m *gltf.Model, motions []human.Motion) error {
	for _, mo := range human.SleepingTogether {
		if !hasMotion(motions, mo) {
			continue
		}
		bed, scale, err := b.doubleBed()
		if err != nil {
			return err
		}
		at := doubleAt
		at[0] = doubleSide(mo) * doubleAside
		if err := human.Lay(m, bed, human.BedPlace{At: at, Yaw: doubleTurn, Scale: scale, During: mo.Name}); err != nil {
			return err
		}
		m.Nodes[m.Roots[len(m.Roots)-1]].Name = bedNode + mo.Name
	}
	return nil
}

// doubleSide is -1 for the left side of the bed and 1 for the right one:
// the way the bed's middle lies off the sleeper, across the bed.
func doubleSide(mo human.Motion) float64 {
	if mo.Name == human.SleepLeft.Name {
		return -1
	}
	return 1
}

// doubleBed is the double bed's model and the scale it is drawn at.
func (b *Baker) doubleBed() (*gltf.Model, float64, error) {
	for _, a := range b.Manifest.Assets {
		if a.ID != DoubleBedModel {
			continue
		}
		bed, err := Load(a)
		if err != nil {
			return nil, 0, fmt.Errorf("%s: %w", a.ID, err)
		}
		return bed, scaleOf(a), nil
	}
	return nil, 0, fmt.Errorf("no asset %s to sleep in as a couple", DoubleBedModel)
}

// SleepShiftSide is how many pixels, across and down, the double bed's
// middle is from the pivot of the picture of the sleeper on side motion mo,
// which the game puts on the bed's middle.
func SleepShiftSide(mo human.Motion) (dx, dy int) {
	return sideShift(mo)
}

// sideShift is the double bed's middle in pixels from a sleeper's pivot:
// across the bed runs up to the right on the screen, two pixels across for
// one up.
func sideShift(mo human.Motion) (int, int) {
	px := -doubleSide(mo) * doubleAside * kit.Cell
	return int(px), -int(px / 2)
}

// withoutBed is sp, a sleeper's picture on a side of the double bed, with
// all that the bed alone has, bare, taken away: what is left of it is the
// sleeper, their shadow on the sheet and the sheet's give, which the game
// lays over the double bed's own picture. Two of them make a couple in it,
// and one alone makes a bed with one in it.
func withoutBed(sp, bare bake.Sprite) bake.Sprite {
	at := func(s bake.Sprite, x, y int) color.NRGBA {
		p := image.Pt(x-sp.Pivot.X+s.Pivot.X, y-sp.Pivot.Y+s.Pivot.Y)
		if !p.In(s.Image.Rect) {
			return color.NRGBA{}
		}
		return s.Image.NRGBAAt(p.X, p.Y)
	}
	out := image.NewNRGBA(sp.Image.Rect)
	for y := out.Rect.Min.Y; y < out.Rect.Max.Y; y++ {
		for x := out.Rect.Min.X; x < out.Rect.Max.X; x++ {
			if c := sp.Image.NRGBAAt(x, y); c != at(bare, x, y) {
				out.SetNRGBA(x, y, c)
			}
		}
	}
	box := image.Rectangle{}
	for y := out.Rect.Min.Y; y < out.Rect.Max.Y; y++ {
		for x := out.Rect.Min.X; x < out.Rect.Max.X; x++ {
			if out.NRGBAAt(x, y).A > 0 {
				box = box.Union(image.Rect(x, y, x+1, y+1))
			}
		}
	}
	if box.Empty() {
		return sp
	}
	crop := image.NewNRGBA(image.Rect(0, 0, box.Dx(), box.Dy()))
	for y := range box.Dy() {
		for x := range box.Dx() {
			crop.SetNRGBA(x, y, out.NRGBAAt(box.Min.X+x, box.Min.Y+y))
		}
	}
	return bake.Sprite{Image: crop, Pivot: sp.Pivot.Sub(box.Min)}
}

// bareBed is the double bed alone of the sleeper on side mo, as it is
// shot, to take away from the picture of the sleeper in it.
func bareBed(m *gltf.Model, mo human.Motion, sh bake.Shot) bake.Sprite {
	bed := *m // the same model, shot with the bed for its only root, where it lies
	bed.Roots = []int{bedRoot(m, mo)}
	return shoot(&bed, frame{shot: sh})
}

// bedRoot is the root of m that is the double bed of the side mo.
func bedRoot(m *gltf.Model, mo human.Motion) int {
	for _, r := range m.Roots {
		if m.Nodes[r].Name == bedNode+mo.Name {
			return r
		}
	}
	return 0
}

// bedNode is the start of the name of the node of a double bed, the
// motion it is for ending it.
const bedNode = "double_bed_"
