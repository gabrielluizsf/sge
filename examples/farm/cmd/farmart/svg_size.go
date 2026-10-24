package main

import "fmt"

// svgSize reads the width and height art writes into every drawing.
func svgSize(src string) (w, h int) {
	var a, b int
	fmt.Sscanf(src, `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %d %d" width="%d" height="%d"`, &a, &b, &w, &h)
	return w, h
}
