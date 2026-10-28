// Unix stream sockets: an echo over a socket file, which Close removes.
package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"path/filepath"
)

func main() {
	// The path the conformance runner lets Deno listen on.
	path := filepath.Join(os.TempDir(), "gots-conformance.sock")
	os.Remove(path)
	l, err := net.Listen("unix", path)
	if err != nil {
		fmt.Println("listen:", err)
		return
	}
	go func() {
		c, err := l.Accept()
		if err != nil {
			return
		}
		line, _ := bufio.NewReader(c).ReadString('\n')
		fmt.Fprintf(c, "echo %s", line)
		c.Close()
	}()
	c, err := net.Dial("unix", path)
	if err != nil {
		fmt.Println("dial:", err)
		return
	}
	fmt.Fprintln(c, "over a file")
	got, _ := bufio.NewReader(c).ReadString('\n')
	fmt.Printf("%q %s\n", got, c.RemoteAddr().Network())
	c.Close()
	_, statErr := os.Stat(path)
	fmt.Println("file while listening:", statErr == nil)
	l.Close()
	_, statErr = os.Stat(path)
	fmt.Println("file after close:", os.IsNotExist(statErr))
}
