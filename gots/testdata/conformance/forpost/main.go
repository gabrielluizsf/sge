// Loops whose post statement takes a call with several results, as net
// reads its configuration files line by line.
package main

import (
	"fmt"
	"strings"
)

type file struct{ lines []string }

func (f *file) readLine() (string, bool) {
	if len(f.lines) == 0 {
		return "", false
	}
	l := f.lines[0]
	f.lines = f.lines[1:]
	return l, true
}

func main() {
	f := &file{lines: strings.Split("nameserver 1.1.1.1\n# comment\nsearch example.com", "\n")}
	var funcs []func() string
	for line, ok := f.readLine(); ok; line, ok = f.readLine() {
		if strings.HasPrefix(line, "#") {
			continue
		}
		funcs = append(funcs, func() string { return line })
		fmt.Println(line)
	}
	// Each iteration has its own line, the one its post statement set.
	for _, g := range funcs {
		fmt.Println("kept:", g())
	}
	xs := []int{3, 1, 2}
	for i, j := 0, len(xs)-1; i < j; i, j = i+1, j-1 {
		xs[i], xs[j] = xs[j], xs[i]
	}
	fmt.Println(xs)
}
