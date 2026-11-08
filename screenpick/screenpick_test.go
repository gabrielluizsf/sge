package screenpick

import "testing"

func TestChannelOf(t *testing.T) {
	for _, tc := range []struct {
		pixel, mask uint64
		want        uint8
	}{
		{0xFF8040, 0xFF0000, 0xFF}, {0xFF8040, 0x00FF00, 0x80}, {0xFF8040, 0x0000FF, 0x40},
		{0xF800, 0xF800, 255}, {0x07E0, 0x07E0, 255}, {0x0000, 0xF800, 0},
		{0x3FF00000, 0x3FF00000, 255}, {1, 0, 0}, {^uint64(0), ^uint64(0), 255},
	} {
		if got := channelOf(tc.pixel, tc.mask); got != tc.want {
			t.Errorf("channelOf(%#x, %#x) = %d, want %d", tc.pixel, tc.mask, got, tc.want)
		}
	}
}

func TestUsable(t *testing.T) {
	env := func(m map[string]string) func(string) string { return func(k string) string { return m[k] } }
	for name, tc := range map[string]struct {
		env map[string]string
		ok  bool
	}{
		"x11":            {map[string]string{"DISPLAY": ":0", "XDG_SESSION_TYPE": "x11"}, true},
		"no session var": {map[string]string{"DISPLAY": ":0"}, true},
		"wayland":        {map[string]string{"DISPLAY": ":0", "XDG_SESSION_TYPE": "wayland"}, false},
		"wayland only":   {map[string]string{"WAYLAND_DISPLAY": "wayland-0"}, false},
		"no display":     {map[string]string{}, false},
	} {
		if err := usable(env(tc.env)); (err == nil) != tc.ok {
			t.Errorf("%s: usable = %v, want ok=%v", name, err, tc.ok)
		}
	}
}
