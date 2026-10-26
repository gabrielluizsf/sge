// Closing a connection while a goroutine waits to read from it, as a
// server shutting down does: the read ends with net.ErrClosed.
package main

import (
	"errors"
	"fmt"
	"net"
	"time"
)

func main() {
	l, _ := net.Listen("tcp", "127.0.0.1:0")
	accepted := make(chan net.Conn)
	go func() {
		c, _ := l.Accept()
		accepted <- c
	}()
	client, _ := net.Dial("tcp", l.Addr().String())
	server := <-accepted
	done := make(chan error)
	go func() {
		_, err := server.Read(make([]byte, 8))
		done <- err
	}()
	time.Sleep(20 * time.Millisecond)
	server.Close()
	select {
	case err := <-done:
		fmt.Println("read ended:", errors.Is(err, net.ErrClosed))
	case <-time.After(2 * time.Second):
		fmt.Println("read still waiting")
	}
	client.Close()
	l.Close()
}
