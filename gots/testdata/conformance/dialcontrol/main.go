// A Dialer whose Control refuses an address, as servers fetching URLs
// from outside refuse private ones: the dial fails, never connects.
package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"syscall"
)

func main() {
	l, _ := net.Listen("tcp", "127.0.0.1:0")
	defer l.Close()
	go func() {
		for {
			c, err := l.Accept()
			if err != nil {
				return
			}
			fmt.Println("accepted a connection Control refused")
			c.Close()
		}
	}()
	deny := errors.New("private address")
	d := &net.Dialer{Control: func(network, address string, c syscall.RawConn) error { return deny }}
	_, err := d.Dial("tcp", l.Addr().String())
	fmt.Println("refused:", err != nil)
	dc := &net.Dialer{ControlContext: func(ctx context.Context, network, address string, c syscall.RawConn) error { return deny }}
	_, err = dc.DialContext(context.Background(), "tcp", l.Addr().String())
	fmt.Println("refused by ControlContext:", err != nil)
	lc := net.ListenConfig{Control: func(network, address string, c syscall.RawConn) error { return deny }}
	_, err = lc.Listen(context.Background(), "tcp", "127.0.0.1:0")
	fmt.Println("listen refused:", err != nil)
}
