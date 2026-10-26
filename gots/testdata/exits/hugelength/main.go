// A Content-Length past what a JavaScript number holds exactly, which
// would be read as another length: the server answers it 400 and the
// client refuses it, on Deno, where Go itself would take it.
package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
)

const huge = "9007199254740993" // 2^53 + 1

func main() {
	srv := &http.Server{Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.Copy(io.Discard, r.Body)
		fmt.Fprint(w, "read the body")
	})}
	l, _ := net.Listen("tcp", "127.0.0.1:0")
	go srv.Serve(l)
	c, err := net.Dial("tcp", l.Addr().String())
	if err != nil {
		fmt.Println("dial:", err)
		return
	}
	fmt.Fprintf(c, "POST / HTTP/1.1\r\nHost: x\r\nContent-Length: %s\r\n\r\nabc", huge)
	status, _ := bufio.NewReader(c).ReadString('\n')
	fmt.Printf("server: %q\n", status)
	c.Close()

	// A server that claims such a body.
	raw, _ := net.Listen("tcp", "127.0.0.1:0")
	go func() {
		c, err := raw.Accept()
		if err != nil {
			return
		}
		bufio.NewReader(c).ReadString('\n')
		fmt.Fprintf(c, "HTTP/1.1 200 OK\r\nContent-Length: %s\r\n\r\nabc", huge)
	}()
	_, err = http.Get("http://" + raw.Addr().String() + "/")
	fmt.Println("client refused:", err != nil)
}
