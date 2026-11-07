package kit

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gabrielluizsf/sge/human"
	"github.com/gabrielluizsf/sge/save"
)

func someone(seed uint32) Character {
	b, l := human.Random(human.Spec{Seed: seed})
	return Character{Body: b, Look: l}
}

func TestRandomCharactersAreValid(t *testing.T) {
	for seed := range uint32(60) {
		if err := someone(seed).Check(); err != nil {
			t.Fatalf("seed %d: %v", seed, err)
		}
	}
}

func TestTheCharacterRoundTripsThroughTheSave(t *testing.T) {
	s := save.In(t.TempDir(), characterVersion)
	if _, ok := LoadCharacter(s); ok {
		t.Fatal("a character from nothing")
	}
	c := someone(7)
	if err := SaveCharacter(s, c); err != nil {
		t.Fatal(err)
	}
	got, ok := LoadCharacter(s)
	if !ok || got.Tag() != c.Tag() || got.Body != c.Body || len(got.Look.Wear) != len(c.Look.Wear) {
		t.Errorf("read back %+v ok=%v, saved %+v", got, ok, c)
	}
}

func TestDifferentCharactersHaveDifferentTags(t *testing.T) {
	if someone(1).Tag() == someone(2).Tag() || someone(1).Tag() != someone(1).Tag() {
		t.Error("tags do not tell characters apart")
	}
}

// hostile writes data as the look file and reports whether it was read.
func hostile(t *testing.T, data string) bool {
	t.Helper()
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, characterFile+".json"), []byte(data), 0o600); err != nil {
		t.Fatal(err)
	}
	_, ok := LoadCharacter(save.In(dir, characterVersion))
	return ok
}

func TestHostileLooksAreRefused(t *testing.T) {
	good := someone(3)
	s := save.In(t.TempDir(), characterVersion)
	if err := SaveCharacter(s, good); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(filepath.Join(s.Dir(), characterFile+".json"))
	if err != nil {
		t.Fatal(err)
	}
	text := string(raw)
	if !hostile(t, text) {
		t.Fatal("the good file is refused")
	}
	wear := good.Look.Wear[len(good.Look.Wear)-1]
	cases := map[string]string{
		"unknown proxy":  strings.Replace(text, wear, "clothes/../../etc/passwd", 1),
		"unknown field":  strings.Replace(text, `"Body"`, `"Evil":1,"Body"`, 1),
		"body too big":   strings.Replace(text, `"Muscle":`, `"Muscle":7e3,"Pad":`, 1),
		"negative":       strings.Replace(text, `"Weight":`, `"Weight":-1,"Pad":`, 1),
		"tone too big":   strings.Replace(text, `"SkinTone":`, `"SkinTone":2,"Pad":`, 1),
		"trailing":       text + `{}`,
		"oversized":      strings.Replace(text, `"Wear":[`, `"Wear":["`+strings.Repeat("a", characterLimit)+`",`, 1),
		"not json":       "\x00\x01",
		"wrong version":  strings.Replace(text, `"version":1`, `"version":9`, 1),
		"a skin file":    strings.Replace(text, `"Skin":""`, `"Skin":"skins/x.mhmat"`, 1),
		"bad colour":     strings.Replace(text, `"R":`, `"R":3,"Pad":`, 1),
		"eyes elsewhere": strings.Replace(text, `"Eyes":"eyes/materials/`, `"Eyes":"hair/`, 1),
	}
	for name, data := range cases {
		if name != "unknown field" && name != "wrong version" && data == text {
			t.Errorf("%s: the case changed nothing", name)
		}
		if hostile(t, data) {
			t.Errorf("%s: read", name)
		}
	}
}

func TestCheckRefusesNaNAndTooMuchWear(t *testing.T) {
	c := someone(4)
	c.Body.Height = nan()
	if c.Check() == nil {
		t.Error("NaN height")
	}
	c = someone(4)
	c.Look.Wear = append(c.Look.Wear, c.Look.Wear...)
	if c.Check() == nil {
		t.Error("repeated wear")
	}
	c = someone(4)
	c.Look.ClothesColours = make([]human.Colour, len(c.Look.Wear)+1)
	if c.Check() == nil {
		t.Error("more colours than clothes")
	}
}

func TestSetPlayerTagsTheirPictures(t *testing.T) {
	defer ResetPlayer()
	c := someone(5)
	SetPlayer(c)
	var v Villager
	for _, x := range Villagers {
		if x.ID == "player" {
			v = x
		}
	}
	if v.Body != c.Body || !strings.Contains(v.LookTag(), c.Tag()) && strings.Contains(SpriteTag("player"), c.Tag()) {
		t.Errorf("player %+v, tag %q", v.Body, v.LookTag())
	}
	for _, x := range Villagers {
		if x.ID != "player" && PlayerTag(x) != "" {
			t.Errorf("%s tagged", x.ID)
		}
	}
}

func TestResetPlayerBringsTheDefaultBack(t *testing.T) {
	before := PlayerCharacter()
	SetPlayer(someone(6))
	ResetPlayer()
	if PlayerCharacter().Tag() != before.Tag() || SpriteTag("player") != "" {
		t.Error("the default player did not come back")
	}
}

func TestTheNameRoundTripsAndDoesNotChangeThePictures(t *testing.T) {
	s := save.In(t.TempDir(), characterVersion)
	c := someone(8)
	plain := c.Tag()
	c.Name = "João d'Ávila-Souza"
	if c.Tag() != plain {
		t.Error("a name changed the tag, so the pictures would be baked again")
	}
	if err := SaveCharacter(s, c); err != nil {
		t.Fatal(err)
	}
	got, ok := LoadCharacter(s)
	if !ok || got.Name != c.Name {
		t.Errorf("read back %q, ok=%v", got.Name, ok)
	}
}

func TestCheckRefusesBadNames(t *testing.T) {
	for _, ok := range []string{"", "Ana", "Tupã", "Mary-Jane"} {
		c := someone(9)
		c.Name = ok
		if err := c.Check(); err != nil {
			t.Errorf("%q: %v", ok, err)
		}
	}
	for _, bad := range []string{" Ana", "Ana\n", "A\x00", "<b>", "../x", strings.Repeat("a", 21), "\xff", "A‮B"} {
		c := someone(9)
		c.Name = bad
		if c.Check() == nil {
			t.Errorf("%q passed", bad)
		}
	}
}

// A hostile name in the file is not taken: the character is, with the
// default name.
func TestAHostileNameInTheFileFallsBackToTheDefault(t *testing.T) {
	good := someone(3)
	good.Name = "Ana"
	s := save.In(t.TempDir(), characterVersion)
	if err := SaveCharacter(s, good); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(filepath.Join(s.Dir(), characterFile+".json"))
	if err != nil {
		t.Fatal(err)
	}
	text := string(raw)
	if !strings.Contains(text, `"Name":"Ana"`) {
		t.Fatalf("no name in the file: %s", text)
	}
	for name, bad := range map[string]string{
		"control":  `"Name":"A\u0000na"`,
		"newline":  `"Name":"Ana\nroot"`,
		"long":     `"Name":"` + strings.Repeat("a", 500) + `"`,
		"escape":   `"Name":"\u001b[31mAna"`,
		"path":     `"Name":"../../etc"`,
		"bidi":     `"Name":"‮ana"`,
		"spaces":   `"Name":"  "`,
		"not text": `"Name":7`,
	} {
		dir := t.TempDir()
		data := strings.Replace(text, `"Name":"Ana"`, bad, 1)
		if err := os.WriteFile(filepath.Join(dir, characterFile+".json"), []byte(data), 0o600); err != nil {
			t.Fatal(err)
		}
		got, ok := LoadCharacter(save.In(dir, characterVersion))
		if name == "not text" {
			if ok {
				t.Errorf("%s: read", name)
			}
			continue
		}
		if !ok || got.Name != "" {
			t.Errorf("%s: ok=%v name %q, want the default (empty)", name, ok, got.Name)
		}
	}
}
