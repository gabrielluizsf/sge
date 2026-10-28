// TCP over loopback: a server that echoes lines, a client that sends a
// few, a read that runs out of time and loses nothing, and a refusal.
package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"syscall"
	"time"
)

func serve(l net.Listener) {
	for {
		c, err := l.Accept()
		if err != nil {
			return
		}
		go func() {
			defer c.Close()
			r := bufio.NewReader(c)
			for {
				line, err := r.ReadString('\n')
				if err != nil {
					return
				}
				fmt.Fprintf(c, "echo %s", strings.ToUpper(line))
			}
		}()
	}
}

func main() {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		fmt.Println("listen:", err)
		return
	}
	go serve(l)
	addr := l.Addr().(*net.TCPAddr)
	fmt.Println("listening on loopback:", addr.IP.IsLoopback(), addr.Port > 0)

	c, err := net.Dial("tcp", l.Addr().String())
	if err != nil {
		fmt.Println("dial:", err)
		return
	}
	r := bufio.NewReader(c)
	for _, s := range []string{"hello", "deno", "and go"} {
		fmt.Fprintln(c, s)
		got, err := r.ReadString('\n')
		fmt.Printf("%q %v\n", got, err)
	}

	// Nothing comes: the read runs out of time, then the next one gets
	// what is sent after.
	c.SetReadDeadline(time.Now().Add(50 * time.Millisecond))
	buf := make([]byte, 64)
	_, err = c.Read(buf)
	fmt.Println("timed out:", errors.Is(err, os.ErrDeadlineExceeded))
	c.SetReadDeadline(time.Time{})
	fmt.Fprintln(c, "later")
	got, err := r.ReadString('\n')
	fmt.Printf("%q %v\n", got, err)

	// Half-closing: the server sees the end and closes, the client reads EOF.
	c.(*net.TCPConn).CloseWrite()
	_, err = r.ReadString('\n')
	fmt.Println("server closed:", err == io.EOF)
	c.Close()

	l.Close()
	// A port listened on and let go, with no accept waiting on it.
	free, _ := net.Listen("tcp", "127.0.0.1:0")
	port := free.Addr().(*net.TCPAddr).Port
	free.Close()
	_, err = net.Dial("tcp", fmt.Sprint("127.0.0.1:", port))
	fmt.Println("refused:", errors.Is(err, syscall.ECONNREFUSED))
}
