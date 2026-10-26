// UDP over loopback: an echo, a connected socket, a datagram cut to the
// buffer, and a read that runs out of time.
package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	srv, err := net.ListenPacket("udp", "127.0.0.1:0")
	if err != nil {
		fmt.Println("listen:", err)
		return
	}
	defer srv.Close()
	go func() {
		buf := make([]byte, 64)
		for {
			n, from, err := srv.ReadFrom(buf)
			if err != nil {
				return
			}
			srv.WriteTo(append([]byte("echo "), buf[:n]...), from)
		}
	}()
	c, err := net.Dial("udp", srv.LocalAddr().String())
	if err != nil {
		fmt.Println("dial:", err)
		return
	}
	defer c.Close()
	buf := make([]byte, 64)
	for _, s := range []string{"one", "two"} {
		c.Write([]byte(s))
		c.SetReadDeadline(time.Now().Add(2 * time.Second))
		n, err := c.Read(buf)
		fmt.Printf("%q %v\n", buf[:n], err)
	}
	c.Write([]byte("a datagram longer than its buffer"))
	small := make([]byte, 8)
	n, _ := c.Read(small)
	fmt.Printf("cut: %q\n", small[:n])

	quiet, _ := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1)})
	defer quiet.Close()
	quiet.SetReadDeadline(time.Now().Add(50 * time.Millisecond))
	_, _, err = quiet.ReadFromUDP(buf)
	fmt.Println("timed out:", errors.Is(err, os.ErrDeadlineExceeded))
}
