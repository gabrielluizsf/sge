// main waits on a channel nothing will ever send on.
package main

import "fmt"

func main() {
	c := make(chan int)
	go func() { fmt.Println("worker ran") }()
	<-c
}
