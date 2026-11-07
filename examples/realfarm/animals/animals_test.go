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
