package human

import (
	"archive/zip"
	"crypto/sha256"
	"fmt"
	"image"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/gabrielluizsf/sge/gltf"
)

// realKit is the kit of MakeHuman's real data and CC0 packs, as the cache
// holds them, loaded once for every test that asks: the test is skipped
// where the cache is empty.
func realKit(t *testing.T) *Kit {
	t.Helper()
	if runtime.GOOS == "js" {
		// Making whole people in JavaScript takes minutes a test: the
		// real kit is checked on Go, and on Deno the small kits check the
		// same code.
		t.Skip("the real kit is checked on Go")
	}
	if shared.kit == nil && shared.skip == "" {
		shared.kit, shared.skip = loadRealKit()
	}
	if shared.kit == nil {
		t.Skip(shared.skip)
	}
	return shared.kit
}

// shared is the real kit the tests share, or why there is none; its files
// are kept open for as long as the tests run.
var shared struct {
	kit  *Kit
	skip string
}

// loadRealKit opens the real kit from the cache, or says why it cannot.
func loadRealKit() (*Kit, string) {
	home, err := os.UserCacheDir()
	if err != nil {
		return nil, "no cache"
	}
	base := filepath.Join(home, "sge-assets", "makehuman")
	root, err := os.OpenRoot(filepath.Join(base, "data"))
	if err != nil {
		return nil, "no MakeHuman data in the cache"
	}
	fss := []fs.FS{root.FS()}
	packs, _ := filepath.Glob(filepath.Join(base, "packs", "*.zip"))
	if len(packs) == 0 {
		return nil, "no MakeHuman packs in the cache"
	}
	for _, p := range packs {
		z, err := zip.OpenReader(p)
		if err != nil {
			return nil, err.Error()
		}
		fss = append(fss, z)
	}
	k, err := Load(Layers(fss...), DefaultLimits)
	if err != nil {
		return nil, err.Error()
	}
	return k, ""
}

// modelHash is a fingerprint of a model's materials and images.
func modelHash(m *gltf.Model) string {
	h := sha256.New()
	for _, mat := range m.Materials {
		fmt.Fprintf(h, "%+v\n", mat)
	}
	for _, img := range m.Images {
		if n, ok := img.(*image.NRGBA); ok {
			h.Write(n.Pix)
		}
	}
	return fmt.Sprintf("%x", h.Sum(nil))[:16]
}
