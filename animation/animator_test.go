package animation

import (
	"testing"

	"github.com/gabrielluizsf/antui/canvas"
	"github.com/gabrielluizsf/sge/sprite"
)

func frames(t *testing.T, n int) []*sprite.Sprite {
	t.Helper()
	var out []*sprite.Sprite
	for i := 0; i < n; i++ {
		s, err := sprite.Rect(1, 1, canvas.RGB(uint8(i), 0, 0))
		if err != nil {
			t.Fatal(err)
		}
		out = append(out, s)
	}
	return out
}

func TestLoopingAnimationWrapsAround(t *testing.T) {
	a := NewAnimator()
	fr := frames(t, 3)
	a.Add("walk", &Animation{Frames: fr, FPS: 10, Loop: true})
	if !a.Play("walk") {
		t.Fatal("Play of an added animation said it is not there")
	}
	steps := []struct {
		dt   float64
		want int
	}{{0, 0}, {0.1, 1}, {0.1, 2}, {0.1, 0}, {0.25, 2}}
	for i, s := range steps {
		a.Update(s.dt)
		if got := a.FrameIndex(); got != s.want {
			t.Errorf("step %d: frame %d, want %d", i, got, s.want)
		}
	}
	if a.Frame() != fr[2] || a.Done() {
		t.Error("Frame or Done disagree with the index")
	}
}

func TestOneShotStopsOnTheLastFrame(t *testing.T) {
	a := NewAnimator()
	a.Add("die", &Animation{Frames: frames(t, 2), FPS: 4})
	a.Play("die")
	a.Update(0.3)
	if a.FrameIndex() != 1 || a.Done() {
		t.Errorf("after 0.3s: frame %d done %v, want 1 and not done", a.FrameIndex(), a.Done())
	}
	a.Update(5)
	if a.FrameIndex() != 1 || !a.Done() {
		t.Errorf("long after: frame %d done %v, want 1 and done", a.FrameIndex(), a.Done())
	}
}

func TestPlayingTheSameAnimationCarriesOn(t *testing.T) {
	a := NewAnimator()
	a.Add("walk", &Animation{Frames: frames(t, 4), FPS: 10, Loop: true})
	a.Add("jump", &Animation{Frames: frames(t, 2), FPS: 10})
	a.Play("walk")
	a.Update(0.25)
	a.Play("walk")
	if a.FrameIndex() != 2 {
		t.Errorf("replaying the same animation restarted it at frame %d", a.FrameIndex())
	}
	a.Play("jump")
	if a.Current() != "jump" || a.FrameIndex() != 0 {
		t.Error("switching animation does not start the new one from its first frame")
	}
}

func TestStopFreezes(t *testing.T) {
	a := NewAnimator()
	a.Add("spin", &Animation{Frames: frames(t, 4), FPS: 10, Loop: true})
	a.Play("spin")
	a.Update(0.1)
	a.Stop()
	a.Update(1)
	if a.FrameIndex() != 1 {
		t.Errorf("a stopped animation moved on to frame %d", a.FrameIndex())
	}
	a.Play("spin")
	if a.FrameIndex() != 0 {
		t.Error("playing a stopped animation does not start it over")
	}
}

func TestNothingToPlay(t *testing.T) {
	a := NewAnimator()
	if a.Play("missing") {
		t.Error("Play of an unknown animation reported success")
	}
	if a.Frame() != nil || a.Done() || a.FrameIndex() != 0 {
		t.Error("an empty animator shows something")
	}
}
