// Write deadlines: a write to a peer that never reads runs out of time
// having written part, and a deadline already passed fails at once.
package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	l, _ := net.Listen("tcp", "127.0.0.1:0")
	defer l.Close()
	peer := make(chan net.Conn, 1)
	go func() {
		c, _ := l.Accept()
		peer <- c // and never reads
	}()
	c, err := net.Dial("tcp", l.Addr().String())
	if err != nil {
		fmt.Println("dial:", err)
		return
	}
	defer c.Close()
	p := <-peer // kept, unread: dropped, Go's finalizer would close it
	defer p.Close()
	big := make([]byte, 64<<20)
	c.SetWriteDeadline(time.Now().Add(100 * time.Millisecond))
	start := time.Now()
	n, err := c.Write(big)
	fmt.Println("timed out:", errors.Is(err, os.ErrDeadlineExceeded), "partial:", n < len(big), "soon:", time.Since(start) < 5*time.Second)

	c.SetWriteDeadline(time.Now().Add(-time.Second))
	n, err = c.Write([]byte("x"))
	fmt.Println("passed:", n, errors.Is(err, os.ErrDeadlineExceeded))
}
