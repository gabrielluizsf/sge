// reflect calling functions and methods that block.
package main

import (
	"fmt"
	"reflect"
	"time"
)

type worker struct{ ch chan string }

func (w *worker) Take() string { return <-w.ch }

func slowAdd(a, b int) int {
	time.Sleep(time.Millisecond)
	return a + b
}

func main() {
	out := reflect.ValueOf(slowAdd).Call([]reflect.Value{reflect.ValueOf(2), reflect.ValueOf(3)})
	fmt.Println(out[0].Int())

	w := &worker{ch: make(chan string, 1)}
	w.ch <- "job"
	got := reflect.ValueOf(w).MethodByName("Take").Call(nil)
	fmt.Println(got[0].String())

	recv := reflect.ValueOf(func(c chan int) int { return <-c })
	c := make(chan int, 1)
	c <- 42
	fmt.Println(recv.Call([]reflect.Value{reflect.ValueOf(c)})[0].Int())
	// One that does not block, as before.
	fmt.Println(reflect.ValueOf(fmt.Sprint).Call([]reflect.Value{reflect.ValueOf("ok")})[0])
}
