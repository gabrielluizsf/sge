// A goroutine panics, and nothing recovers: the program dies of it.
package main

import "fmt"

func main() {
	done := make(chan bool)
	go func() {
		panic(fmt.Sprintf("boom %d", 42))
	}()
	<-done
}
