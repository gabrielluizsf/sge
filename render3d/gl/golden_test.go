package gl

import (
	"image"
	"testing"

	"github.com/gabrielluizsf/sge/render3d"
	"github.com/gabrielluizsf/sge/render3d/soft"
)

// draw is s's scene drawn on dev, read back.
func draw(t *testing.T, dev render3d.Device, s scene) *image.NRGBA {
	t.Helper()
	r := render3d.NewRenderer(dev)
	m, err := r.Upload(s.model)
	if err != nil {
		t.Fatal(err)
	}
	if err := r.Draw(s.show(m)); err != nil {
		t.Fatal(err)
	}
	img, err := dev.Read(render3d.Screen)
	if err != nil {
		t.Fatal(err)
	}
	return img
}

// unlike is how many pixels of a and b differ by more than 8 in any
// channel.
func unlike(a, b *image.NRGBA) int {
	bad := 0
	for i := 0; i < len(a.Pix); i += 4 {
		for k := range 4 {
			if d := int(a.Pix[i+k]) - int(b.Pix[i+k]); d > 8 || d < -8 {
				bad++
				break
			}
		}
	}
	return bad
}

// Every scene soft's tests draw, GL draws the same: but for a pixel here
// and there on the edges of triangles, which the two may round apart.
func TestDrawsWhatSoftDraws(t *testing.T) {
	const size = 64
	dev := offscreen(t, size, size)
	for _, s := range scenes() {
		t.Run(s.name, func(t *testing.T) {
			want := draw(t, soft.New(size, size), s)
			got := draw(t, dev, s)
			if bad := unlike(got, want); bad > size*size/100 {
				t.Errorf("%d pixels of %d unlike soft's", bad, size*size)
			}
		})
	}
}
