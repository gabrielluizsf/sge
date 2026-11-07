//go:build linux && amd64

package screenpick

import (
	"os"
	"testing"
)

// TestRealScreen reads the real screen; it runs only with SGE_REAL_SCREEN=1.
func TestRealScreen(t *testing.T) {
	if os.Getenv("SGE_REAL_SCREEN") != "1" {
		t.Skip("set SGE_REAL_SCREEN=1 to read the real screen")
	}
	s, err := Open()
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	r, err := s.Read()
	t.Logf("reading %+v, err %v", r, err)
	if err != nil {
		t.Fatal(err)
	}
}
