// Dialer and ListenConfig Control hooks, as servers fetching URLs from
// outside refuse private addresses with: one that refuses stops the dial
// before it connects, one that allows is asked about the address dialed.
package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"runtime"
	"sync/atomic"
	"syscall"
	"time"
)

func main() {
	l, _ := net.Listen("tcp", "127.0.0.1:0")
	defer l.Close()
	var accepted atomic.Int32
	go func() {
		for {
			c, err := l.Accept()
			if err != nil {
				return
			}
			accepted.Add(1)
			c.Close()
		}
	}()
	deny := errors.New("private address")
	d := &net.Dialer{Control: func(network, address string, c syscall.RawConn) error { return deny }}
	_, err := d.Dial("tcp", l.Addr().String())
	fmt.Println("refused:", errors.Is(err, deny))
	dc := &net.Dialer{ControlContext: func(ctx context.Context, network, address string, c syscall.RawConn) error { return deny }}
	_, err = dc.DialContext(context.Background(), "tcp", l.Addr().String())
	fmt.Println("refused by ControlContext:", errors.Is(err, deny))
	time.Sleep(20 * time.Millisecond)
	fmt.Println("connections refused dials made:", accepted.Load())

	var asked []string
	allow := &net.Dialer{Control: func(network, address string, c syscall.RawConn) error {
		asked = append(asked, network+" "+address)
		return nil
	}}
	c, err := allow.Dial("tcp", l.Addr().String())
	fmt.Println("allowed:", err == nil, len(asked) == 1 && asked[0] == "tcp4 "+l.Addr().String())
	if err == nil {
		c.Close()
	}
	// A hook that wants the descriptor: natively it gets one, on Deno it
	// is told there is none, and the dial fails rather than skip it.
	fd := &net.Dialer{Control: func(network, address string, c syscall.RawConn) error {
		return c.Control(func(uintptr) {})
	}}
	_, err = fd.Dial("tcp", l.Addr().String())
	fmt.Println("descriptor:", (err == nil) == (runtime.GOOS != "js"))

	lc := net.ListenConfig{Control: func(network, address string, c syscall.RawConn) error { return deny }}
	_, err = lc.Listen(context.Background(), "tcp", "127.0.0.1:0")
	fmt.Println("listen refused:", errors.Is(err, deny))
}
