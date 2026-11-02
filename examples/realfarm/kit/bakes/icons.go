package bakes

import (
	"image"
	"image/draw"

	"github.com/gabrielluizsf/sge/bake"
	"github.com/gabrielluizsf/sge/building"
	"github.com/gabrielluizsf/sge/flora"
	"github.com/gabrielluizsf/sge/gltf"
)

// IconSize is the side, in pixels, of an icon as it is baked: twice what
// the HUD shows it at, so it is shrunk, not stretched.
const IconSize = 64

// iconMargin is the room left round a thing in its icon.
const iconMargin = 4

// IconName is the sprite of the icon called name, as the farm asks for it
// ("icon_can", "icon_carrot", "icon_carrotSeeds"): the same.
func IconName(name string) string { return name }

// iconTasks bake the pictures of what the player carries: each item, and
// each crop as picked and as seeds, rendered from their models.
func (b *Baker) iconTasks() []Task {
	var out []Task
	for _, k := range building.Items {
		out = append(out, Task{Name: "icon " + k.Name, Kind: "icon", ID: k.Name, Core: true, Run: func() error { return b.itemIcon(k) }})
	}
	for _, c := range flora.Crops {
		if c.Name == "withered" {
			continue // a plant's state, not an item
		}
		out = append(out, Task{Name: "icon " + c.Name, Kind: "icon", ID: c.Name, Core: true, Run: func() error { return b.cropIcons(c) }})
	}
	return out
}

// itemIcon bakes item k as icon_<name>.
func (b *Baker) itemIcon(k building.Building) error {
	palette, err := b.Palette()
	if err != nil {
		return err
	}
	return b.Shelf.Put("icon_"+k.Name, Icon(k.Model(palette)))
}

// cropIcons bakes crop c as picked, its ripe plant, and as seeds, a
// sprout: icon_<name> and icon_<name>Seeds.
func (b *Baker) cropIcons(c flora.Crop) error {
	if err := b.Shelf.Put("icon_"+c.Name, Icon(c.Grow(c.Stages-1, 1))); err != nil {
		return err
	}
	return b.Shelf.Put("icon_"+c.Name+"Seeds", Icon(c.Grow(1, 2)))
}

// Icon is model m seen from the game's corner without a shadow, as large
// as fits an icon and in the middle of it.
func Icon(m *gltf.Model) bake.Sprite {
	o := Options()
	o.Shadow = 0
	shot := bake.Shot{Yaw: 45}
	first := bake.Render(m, shot, o)
	b := first.Image.Bounds()
	o.Unit *= float64(IconSize-2*iconMargin) / float64(max(b.Dx(), b.Dy(), 1))
	s := bake.Render(m, shot, o)
	out := image.NewNRGBA(image.Rect(0, 0, IconSize, IconSize))
	sb := s.Image.Bounds()
	at := image.Pt((IconSize-sb.Dx())/2, (IconSize-sb.Dy())/2)
	draw.Draw(out, image.Rectangle{Min: at, Max: at.Add(sb.Size())}, s.Image, sb.Min, draw.Src)
	return bake.Sprite{Image: out, Pivot: image.Pt(IconSize/2, IconSize/2)}
}
