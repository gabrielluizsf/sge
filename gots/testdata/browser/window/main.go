// window: a window in a browser, which browser_test drives: it draws red,
// and tells the keys, clicks and text the page sends it.
package main

import (
	"fmt"

	"github.com/gabrielluizsf/antui"
	"github.com/gabrielluizsf/antui/canvas"
)

func main() {
	win, err := antui.Open("gots browser test", 64, 48)
	if err != nil {
		fmt.Println("open:", err)
		return
	}
	var key, click, text bool
	for frame := 0; frame < 1000 && !(key && click && text) && win.Begin(); frame++ {
		win.Clear(canvas.Red)
		if win.KeyPressed(antui.KeyRight) {
			key = true
			fmt.Println("key right")
		}
		if win.MousePressed(antui.MouseLeft) {
			click = true
			fmt.Println("click left")
		}
		if t := win.TextInput(); t != "" {
			text = true
			fmt.Println("text", t)
		}
		win.End()
	}
	fmt.Println("done")
}
