// udpcontrol: a Dialer's Control hook sees the address family of the UDP
// peer it is about to reach, udp4 or udp6, as on every other platform.
package main

import (
	"fmt"
	"net"
	"syscall"
)

func main() {
	srv, err := net.ListenPacket("udp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	defer srv.Close()
	d := net.Dialer{Control: func(network, address string, _ syscall.RawConn) error {
		fmt.Println("control", network)
		return nil
	}}
	c, err := d.Dial("udp", srv.LocalAddr().String())
	if err != nil {
		panic(err)
	}
	c.Close()
}
