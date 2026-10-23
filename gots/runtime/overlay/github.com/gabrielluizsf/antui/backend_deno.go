//go:build gots_overlay

// Overlay for backend_other.go: on Deno a window is the native one the
// gotsWindow bridge opens through SDL2, reached with syscall/js. Without
// the bridge Open fails as it does on any platform without a port, and
// drawing to a Canvas still works.

package antui

import (
	"errors"
	"syscall/js"

	"github.com/gabrielluizsf/antui/backend"
	"github.com/gabrielluizsf/antui/canvas"
)

type noWindows struct {
	w js.Value
}

func newBackend() platform { return &noWindows{} }

var errNoBackend = errors.New("antui: there is no window bridge to open a window with")

func (n *noWindows) open(win *Window, title string, width, height int) error {
	bridge := js.Global().Get("gotsWindow")
	if bridge.IsUndefined() {
		return errNoBackend
	}
	w := bridge.Call("open", title, width, height)
	if w.Type() == js.TypeString {
		return errors.New("antui: " + w.String())
	}
	n.w = w
	return nil
}

func (n *noWindows) close() {
	if !n.w.IsUndefined() {
		n.w.Call("close")
	}
}

func (n *noWindows) pump(win *Window) {
	events := n.w.Call("poll")
	for i := 0; i < events.Length(); i++ {
		n.push(win, events.Index(i))
	}
}

// push hands one bridge event to the window as the backend events it is.
func (n *noWindows) push(win *Window, ev js.Value) {
	switch ev.Get("type").String() {
	case "close":
		win.PushSimple(backend.EventClose)
	case "expose":
		win.PushSimple(backend.EventExpose)
	case "focus":
		win.Push(backend.Event{Type: backend.EventFocus, Focused: ev.Get("focused").Bool()})
	case "resize":
		w, h := ev.Get("width").Int(), ev.Get("height").Int()
		if win.ResizeCanvas(w, h) {
			win.Push(backend.Event{Type: backend.EventResize, Width: w, Height: h})
		}
	case "keydown", "keyup":
		t := backend.EventKeyUp
		if ev.Get("type").String() == "keydown" {
			t = backend.EventKeyDown
		}
		win.Push(backend.Event{
			Type:   t,
			Key:    sdlKey(ev.Get("sym").Int()),
			Mods:   sdlMods(ev.Get("mods").Int()),
			Repeat: ev.Get("repeat").Bool(),
		})
	case "text":
		for _, r := range ev.Get("text").String() {
			win.Push(backend.Event{Type: backend.EventText, Rune: r, Text: string(r)})
		}
	case "move":
		win.Push(backend.Event{
			Type: backend.EventMouseMove,
			X:    ev.Get("x").Int(), Y: ev.Get("y").Int(),
			DX: ev.Get("dx").Int(), DY: ev.Get("dy").Int(),
		})
	case "down", "up":
		var b backend.MouseButton
		switch ev.Get("button").Int() {
		case 1:
			b = backend.MouseLeft
		case 2:
			b = backend.MouseMiddle
		case 3:
			b = backend.MouseRight
		default:
			return
		}
		t := backend.EventMouseUp
		if ev.Get("type").String() == "down" {
			t = backend.EventMouseDown
		}
		win.Push(backend.Event{Type: t, Button: b, X: ev.Get("x").Int(), Y: ev.Get("y").Int()})
	case "wheel":
		win.Push(backend.Event{Type: backend.EventMouseWheel, Wheel: ev.Get("wheel").Int()})
	case "drop":
		files := ev.Get("files")
		var paths []string
		for i := 0; i < files.Length(); i++ {
			paths = append(paths, files.Index(i).String())
		}
		win.Push(backend.Event{Type: backend.EventDropFiles, Files: paths})
	}
}

// sdlKey is the antui key for an SDL keycode.
func sdlKey(sym int) backend.Key {
	switch {
	case sym >= 'a' && sym <= 'z':
		return backend.Key(sym - 'a' + 'A')
	case sym >= '0' && sym <= '9', sym == ' ', sym == '\'', sym == ',', sym == '-',
		sym == '.', sym == '/', sym == ';', sym == '=', sym == '[', sym == '\\',
		sym == ']', sym == '`':
		return backend.Key(sym)
	case sym >= 0x4000003A && sym <= 0x40000045:
		return backend.KeyF1 + backend.Key(sym-0x4000003A)
	}
	switch sym {
	case 13:
		return backend.KeyEnter
	case 27:
		return backend.KeyEscape
	case 8:
		return backend.KeyBackspace
	case 9:
		return backend.KeyTab
	case 127:
		return backend.KeyDelete
	case 0x40000039:
		return backend.KeyCapsLock
	case 0x40000049:
		return backend.KeyInsert
	case 0x4000004A:
		return backend.KeyHome
	case 0x4000004B:
		return backend.KeyPageUp
	case 0x4000004D:
		return backend.KeyEnd
	case 0x4000004E:
		return backend.KeyPageDown
	case 0x4000004F:
		return backend.KeyRight
	case 0x40000050:
		return backend.KeyLeft
	case 0x40000051:
		return backend.KeyDown
	case 0x40000052:
		return backend.KeyUp
	case 0x400000E0:
		return backend.KeyLeftControl
	case 0x400000E1:
		return backend.KeyLeftShift
	case 0x400000E2:
		return backend.KeyLeftAlt
	case 0x400000E3:
		return backend.KeyLeftSuper
	case 0x400000E4:
		return backend.KeyRightControl
	case 0x400000E5:
		return backend.KeyRightShift
	case 0x400000E6:
		return backend.KeyRightAlt
	case 0x400000E7:
		return backend.KeyRightSuper
	}
	return backend.KeyUnknown
}

// sdlMods are the antui modifiers for SDL's.
func sdlMods(m int) backend.Mod {
	var out backend.Mod
	if m&0x0003 != 0 {
		out |= backend.ModShift
	}
	if m&0x00C0 != 0 {
		out |= backend.ModControl
	}
	if m&0x0300 != 0 {
		out |= backend.ModAlt
	}
	if m&0x0C00 != 0 {
		out |= backend.ModSuper
	}
	return out
}

func (n *noWindows) present(win *Window, dirty canvas.Area) bool {
	cv := win.Canvas()
	if cv.Pixels == nil {
		return false
	}
	return n.w.Call("present", cv.Pixels, cv.Stride, cv.Width, cv.Height).Bool()
}

func (n *noWindows) setTitle(title string) { n.w.Call("setTitle", title) }

func (n *noWindows) setFullscreen(on bool) bool { return n.w.Call("setFullscreen", on).Bool() }

func (n *noWindows) setOpacity(a uint8) bool { return n.w.Call("setOpacity", int(a)).Bool() }

func (n *noWindows) displaySize() (int, int, bool) {
	d := n.w.Call("display")
	if d.IsNull() {
		return 0, 0, false
	}
	return d.Index(0).Int(), d.Index(1).Int(), true
}

func (n *noWindows) displayRefresh() int {
	d := n.w.Call("display")
	if d.IsNull() {
		return 0
	}
	return d.Index(2).Int()
}

func (n *noWindows) setLimits(win *Window, l Limits) {
	n.w.Call("setLimits", l.MinWidth, l.MinHeight, l.MaxWidth, l.MaxHeight)
}

func (n *noWindows) setSize(win *Window, width, height int) bool {
	return n.w.Call("setSize", width, height).Bool()
}

func (n *noWindows) displayScale() (float64, bool) { return 0, false }

func (n *noWindows) systemDark() (bool, bool) { return false, false }

func (n *noWindows) clipboard() (string, []string) {
	return n.w.Call("clipboard").String(), nil
}

func (n *noWindows) setIcon([]*canvas.Canvas) bool { return false }

func (n *noWindows) setCursor(CursorIcon) bool { return false }

func (n *noWindows) sendNotification(string, string) bool { return false }

func (n *noWindows) setClipboard(text string) bool {
	return n.w.Call("setClipboard", text).Bool()
}
