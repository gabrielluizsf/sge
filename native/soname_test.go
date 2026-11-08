package native

import "testing"

// Library planting: anything dlopen would read as a path, or that smuggles
// a NUL, must be refused before it reaches the loader.
func TestValidSonameRefusesPaths(t *testing.T) {
	for _, s := range []string{
		"", "./libEGL.so.1", "/tmp/libEGL.so.1", "../lib/libGL.so", "libGL.so.1/../../x",
		"$ORIGIN/libGL.so.1", "libEGL.so.1\x00/tmp/x", "libGL.so.1 ", "libGL.so.", "libGL.so.1.",
		"libGL.so.x", "libGL", "lib.so.1", "EGL.so.1", "libGL.dll", "libGL.so.1;rm",
	} {
		if validSoname(s) {
			t.Errorf("validSoname(%q) = true", s)
		}
	}
}

func TestValidSonameAcceptsSystemNames(t *testing.T) {
	for _, s := range []string{"libEGL.so.1", "libGL.so.1", "libGLX.so.0", "libX11.so.6", "libc.so.6", "libm.so.6", "libvulkan.so.1", "libstdc++.so.6"} {
		if !validSoname(s) {
			t.Errorf("validSoname(%q) = false", s)
		}
	}
}

func TestOpenRefusesPathAndCallRefusesNil(t *testing.T) {
	if _, err := Open("./libc.so.6"); err == nil {
		t.Fatal("Open accepted a relative path")
	}
	defer func() {
		if recover() == nil {
			t.Fatal("Call through a nil Proc did not panic")
		}
	}()
	var f Frame
	f.Call(0)
}

func TestCallThroughLibc(t *testing.T) {
	l, err := Open("libc.so.6")
	if err != nil {
		t.Fatal(err)
	}
	labs, err := l.Symbol("labs")
	if err != nil {
		t.Fatal(err)
	}
	var f Frame
	if got := int64(f.Int(uintptr(1<<64 - 42)).Call(labs)); got != 42 {
		t.Fatalf("labs(-42) = %d", got)
	}
}
