package farm

import (
	"fmt"
	"slices"
	"testing"
)

// However many pictures are drawn, the cache holds a bounded few, and a
// picture dropped is drawn again the same.
func TestThePictureCacheIsBounded(t *testing.T) {
	svg := func(i int) string {
		return fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 4 4"><rect width="%d" height="4" fill="#c33"/></svg>`, 1+i%4)
	}
	first := picture(svg(1), 4, 4)
	want := slices.Clone(first.Pixels)
	for i := range 3 * pictureKeep {
		picture(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 4 4"><rect width="4" height="%d" fill="#3%03x"/></svg>`, 1+i%4, i), 4, 4)
	}
	if n := pictures.size(); n > 2*pictureKeep {
		t.Fatalf("%d pictures held", n)
	}
	again := picture(svg(1), 4, 4)
	if !slices.Equal(again.Pixels, want) {
		t.Fatal("a picture drawn again differs")
	}
}
