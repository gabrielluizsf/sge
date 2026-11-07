package animals

import "testing"

// Every frame animals.json promises is embedded and placed, in all eight
// facings, so the game never finds a hole in a walk.
func TestEveryPromisedFrameIsThere(t *testing.T) {
	infos, err := Infos()
	if err != nil || len(infos) != 4 {
		t.Fatalf("animals.json: %v, %d animals", err, len(infos))
	}
	sprites, err := Sprites()
	if err != nil {
		t.Fatal(err)
	}
	for name, in := range infos {
		for _, coat := range in.Coats {
			for motion, m := range in.Motions {
				for d := range 8 {
					for f := range m.Frames {
						n := Frame(name, coat, motion, d, f)
						if sprites[n].W == 0 {
							t.Fatalf("%s is not placed", n)
						}
						if b, err := PNG(n); err != nil || len(b) == 0 {
							t.Fatalf("%s is not embedded: %v", n, err)
						}
					}
				}
			}
		}
	}
}

// Every coat of every animal, in every motion and facing, has a picture of
// the size of an animal: a coat with a blank one is an animal no one can see.
func TestEveryFrameHasAPicture(t *testing.T) {
	infos, err := Infos()
	if err != nil {
		t.Fatal(err)
	}
	places, err := Sprites()
	if err != nil {
		t.Fatal(err)
	}
	for kind, in := range infos {
		for _, coat := range in.Coats {
			for name, m := range in.Motions {
				for d := range 8 {
					for f := range m.Frames {
						at, ok := places[Frame(kind, coat, name, d, f)]
						if !ok || at.W < 16 || at.H < 16 {
							t.Errorf("%s: %+v", Frame(kind, coat, name, d, f), at)
						}
					}
				}
			}
		}
	}
}
