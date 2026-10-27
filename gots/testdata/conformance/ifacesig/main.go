// ifacesig: a value implements an interface only if its methods have the
// interface's types, not only its names: a TCPConn's ReadFrom reads a
// Reader, so it is no PacketConn.
package main

import (
	"fmt"
	"io"
	"net"
)

type reader struct{}

func (reader) ReadFrom(r io.Reader) (int64, error) { return 0, nil }

type packets struct{}

func (packets) ReadFrom(b []byte) (int, net.Addr, error) { return 0, nil, nil }

type readsFrom interface {
	ReadFrom(b []byte) (int, net.Addr, error)
}

type ptrOnly struct{}

func (*ptrOnly) ReadFrom(b []byte) (int, net.Addr, error) { return 0, nil, nil }

func main() {
	for _, v := range []any{reader{}, packets{}, &ptrOnly{}, (*ptrOnly)(nil), &net.TCPConn{}, 3} {
		_, ok := v.(readsFrom)
		_, rf := v.(io.ReaderFrom)
		fmt.Printf("%T %v %v\n", v, ok, rf)
	}
	var e error = &net.OpError{Op: "dial"}
	_, isTimeout := e.(interface{ Timeout() bool })
	_, isWrong := e.(interface{ Timeout() int })
	fmt.Println(isTimeout, isWrong)
	defer func() { fmt.Println("recovered:", recover() != nil) }()
	_ = any(reader{}).(readsFrom)
}
