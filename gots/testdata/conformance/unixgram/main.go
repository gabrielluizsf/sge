// unixgram: Unix datagram sockets, a listener and a client with a local
// address of its own, replies to the sender, deadlines, and packet
// sockets, which Deno has none of, refused. The paths are the ones the
// conformance run grants Deno.
package main

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"time"
)

func main() {
	srvPath := filepath.Join(os.TempDir(), "gots-conformance-gram.sock")
	cliPath := filepath.Join(os.TempDir(), "gots-conformance-gram-cli.sock")
	os.Remove(srvPath)
	os.Remove(cliPath)
	defer os.Remove(srvPath)
	defer os.Remove(cliPath)
	srv, err := net.ListenUnixgram("unixgram", &net.UnixAddr{Name: srvPath, Net: "unixgram"})
	if err != nil {
		panic(err)
	}
	defer srv.Close()
	go func() {
		b := make([]byte, 64)
		for {
			n, from, err := srv.ReadFromUnix(b)
			if err != nil {
				return
			}
			srv.WriteToUnix(append([]byte("echo "), b[:n]...), from)
		}
	}()
	c, err := net.DialUnix("unixgram", &net.UnixAddr{Name: cliPath, Net: "unixgram"}, &net.UnixAddr{Name: srvPath, Net: "unixgram"})
	if err != nil {
		panic(err)
	}
	defer c.Close()
	c.SetReadDeadline(time.Now().Add(5 * time.Second))
	for _, m := range []string{"one", "two"} {
		c.Write([]byte(m))
		b := make([]byte, 64)
		n, err := c.Read(b)
		fmt.Println(string(b[:n]), err, c.LocalAddr(), c.RemoteAddr())
	}
	c.SetReadDeadline(time.Now().Add(10 * time.Millisecond))
	_, err = c.Read(make([]byte, 8))
	fmt.Println("timeout:", err.(net.Error).Timeout())
	_, err = net.Dial("unixpacket", srvPath)
	fmt.Println("unixpacket refused:", err != nil)
}
