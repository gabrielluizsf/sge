package human

import (
	"io/fs"
	"strings"
	"testing"
)

func TestCatalogueIsBigEnough(t *testing.T) {
	c := Catalogue()
	if len(hairStyles) < 10 || len(eyeColours) < 8 || len(c) < 80 {
		t.Errorf("%d hair styles, %d eye colours, %d entries", len(hairStyles), len(eyeColours), len(c))
	}
	seen := map[string]bool{}
	for _, name := range c {
		if seen[name] {
			t.Errorf("%s twice", name)
		}
		seen[name] = true
	}
}

func TestEveryCatalogueEntryIsInThePacksAndCC0(t *testing.T) {
	k := realKit(t)
	for _, name := range Catalogue() {
		data, err := fs.ReadFile(k.fsys, name)
		if err != nil {
			t.Errorf("%s: %v", name, err)
			continue
		}
		if strings.HasSuffix(name, ".mhclo") {
			p, err := parseProxy(data, len(k.base.pos))
			if err != nil {
				t.Errorf("%s: %v", name, err)
			} else if _, err := fs.Stat(k.fsys, beside0(name, p.obj)); err != nil {
				t.Errorf("%s: %v", name, err)
			}
		}
		if lic, err := k.Licence(name); err != nil || lic != "CC0" {
			t.Errorf("%s is %q under %v", name, lic, err)
		}
	}
}

// beside0 is file name of the folder of proxy.
func beside0(proxy, file string) string {
	n, _ := beside(proxy, file)
	return n
}

func TestRandomPeopleAreMadeFromCC0Only(t *testing.T) {
	k := realKit(t)
	for seed := uint32(1); seed <= 50; seed++ {
		b, l := Random(Spec{Seed: seed})
		lic, err := k.Licences(b, l)
		if err != nil {
			t.Fatalf("seed %d: %v", seed, err)
		}
		for name, v := range lic {
			if v != "CC0" {
				t.Errorf("seed %d: %s is %q", seed, name, v)
			}
		}
	}
	b, l := Random(Spec{Seed: 3})
	if _, err := k.Make(b, l); err != nil {
		t.Fatal(err)
	}
}
