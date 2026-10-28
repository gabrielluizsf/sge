// A fault JavaScript raises, recovered as Go's runtime.Error, prints and
// behaves as an error.
package main

import "fmt"

func write() (r any) {
	defer func() { r = recover() }()
	var m map[string]int
	m["x"] = 1
	return nil
}

func main() {
	r := write()
	_, isErr := r.(error)
	fmt.Println(r != nil, isErr, fmt.Sprint(r) != "", fmt.Sprintf("%v", r) != "")
}
