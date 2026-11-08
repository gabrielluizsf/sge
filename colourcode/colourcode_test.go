package colourcode

import (
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	for _, tc := range []struct {
		in   string
		want RGBA
		ok   bool
	}{
		{"#f00", RGBA{255, 0, 0, 255}, true},
		{"#FfA", RGBA{255, 255, 170, 255}, true},
		{"#336699", RGBA{0x33, 0x66, 0x99, 255}, true},
		{"  #336699  ", RGBA{0x33, 0x66, 0x99, 255}, true},
		{"#33669980", RGBA{0x33, 0x66, 0x99, 0x80}, true},
		{"rgba(10, 20, 30, 0.5)", RGBA{10, 20, 30, 128}, true},
		{"RGBA(10,20,30,1)", RGBA{10, 20, 30, 255}, true},
		{"rgba(10, 20, 30, .25)", RGBA{10, 20, 30, 64}, true},
		{"rgba(10, 20, 30, 50%)", RGBA{10, 20, 30, 128}, true},
		{"rgba(10, 20, 30, 0)", RGBA{10, 20, 30, 0}, true},
		{"rgb(255, 255, 255)", RGBA{255, 255, 255, 255}, true},
		{"", RGBA{}, false},
		{"   ", RGBA{}, false},
		{"#", RGBA{}, false},
		{"#12", RGBA{}, false},
		{"#1234", RGBA{}, false},
		{"#12345", RGBA{}, false},
		{"#1234567", RGBA{}, false},
		{"#123456789", RGBA{}, false},
		{"#gggggg", RGBA{}, false},
		{"#+12345", RGBA{}, false},
		{"ff0000", RGBA{}, false},
		{"rgba(256, 0, 0, 1)", RGBA{}, false},
		{"rgba(-1, 0, 0, 1)", RGBA{}, false},
		{"rgba(99999999999999999999, 0, 0, 1)", RGBA{}, false},
		{"rgba(1, 2, 3, 1.5)", RGBA{}, false},
		{"rgba(1, 2, 3, 101%)", RGBA{}, false},
		{"rgba(1, 2, 3, 1e9)", RGBA{}, false},
		{"rgba(1, 2, 3, 1.2.3)", RGBA{}, false},
		{"rgba(1, 2, 3, .)", RGBA{}, false},
		{"rgba(1, 2, 3)", RGBA{}, false},
		{"rgb(1, 2, 3, 4)", RGBA{}, false},
		{"rgba(1, 2, 3, 1", RGBA{}, false},
		{"rgba(1, 2, , 1)", RGBA{}, false},
		{"rgba(1.5, 2, 3, 1)", RGBA{}, false},
		{"rgba(0x10, 2, 3, 1)", RGBA{}, false},
		{"rgba(1, 2, 3, 1)x", RGBA{}, false},
		{"rgba()", RGBA{}, false},
		{"#fff\x00", RGBA{}, false},
		{"#fféff", RGBA{}, false},
		{"javascript:alert(1)", RGBA{}, false},
		{"#" + strings.Repeat("f", 100000), RGBA{}, false},
		{"rgba(" + strings.Repeat("1,", 100000) + ")", RGBA{}, false},
		{strings.Repeat("(", 100000), RGBA{}, false},
	} {
		got, err := Parse(tc.in)
		if (err == nil) != tc.ok || got != tc.want {
			t.Errorf("Parse(%.30q) = %v, %v; want %v ok=%v", tc.in, got, err, tc.want, tc.ok)
		}
	}
}

func TestFormatRoundTrips(t *testing.T) {
	for _, c := range []RGBA{{0, 0, 0, 255}, {255, 128, 1, 255}, {1, 2, 3, 4}} {
		got, err := Parse(Format(c))
		if err != nil || got != c {
			t.Errorf("%v -> %q -> %v, %v", c, Format(c), got, err)
		}
	}
}

func TestInsertFiltersAndLimits(t *testing.T) {
	if got := Insert("#", "a;<>/\\12"); got != "#a12" {
		t.Errorf("Insert kept %q", got)
	}
	got := Insert("", strings.Repeat("f", 1000))
	if len(got) != MaxLen {
		t.Errorf("Insert grew to %d, want %d", len(got), MaxLen)
	}
	if Insert("#", "\x00‮\n") != "#" {
		t.Error("Insert kept control characters")
	}
}

func FuzzParse(f *testing.F) {
	f.Add("#fff")
	f.Add("rgba(1,2,3,.5)")
	f.Fuzz(func(t *testing.T, s string) {
		if c, err := Parse(s); err == nil {
			if again, err := Parse(Format(c)); err != nil || again != c {
				t.Fatalf("%q parsed to %v but %q did not round trip", s, c, Format(c))
			}
		}
	})
}
