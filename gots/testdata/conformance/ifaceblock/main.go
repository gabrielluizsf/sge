// Methods that block, reached through interfaces: a reader and a writer
// that wait on channels, read and written by bufio, io.Copy and fmt from
// goroutines, while the same code over strings stays synchronous.
package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"sync"
)

// pipeReader hands out what is sent on its channel, a chunk a read, and
// ends when the channel is closed.
type pipeReader struct {
	ch   chan []byte
	left []byte
}

func (p *pipeReader) Read(b []byte) (int, error) {
	if len(p.left) == 0 {
		chunk, ok := <-p.ch
		if !ok {
			return 0, io.EOF
		}
		p.left = chunk
	}
	n := copy(b, p.left)
	p.left = p.left[n:]
	return n, nil
}

// slowWriter passes what it is given to a channel.
type slowWriter struct{ ch chan string }

func (w slowWriter) Write(b []byte) (int, error) {
	w.ch <- string(b)
	return len(b), nil
}

// lines reads the lines of r through bufio, as a server reads a request.
func lines(r io.Reader) []string {
	var out []string
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		out = append(out, sc.Text())
	}
	return out
}

func main() {
	// Synchronous: the same code over a strings.Reader.
	fmt.Println(lines(strings.NewReader("one\ntwo\nthree")))

	p := &pipeReader{ch: make(chan []byte)}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("read:", lines(p))
	}()
	for _, s := range []string{"al", "pha\nbe", "ta\n", "gamma"} {
		p.ch <- []byte(s)
	}
	close(p.ch)
	wg.Wait()

	// io.Copy and fmt.Fprintf write through io.Writer to a writer that
	// waits, from a goroutine; a method value and a closure carry it too.
	w := slowWriter{ch: make(chan string)}
	done := make(chan bool)
	go func() {
		var got []string
		for s := range w.ch {
			got = append(got, s)
		}
		fmt.Printf("written: %q\n", got)
		done <- true
	}()
	go func() {
		io.Copy(w, strings.NewReader("copied"))
		fmt.Fprintf(w, "%s-%d", "printed", 7)
		write := w.Write
		write([]byte("by method value"))
		say := func(s string) { io.WriteString(w, s) }
		say("by closure")
		close(w.ch)
	}()
	<-done
	fmt.Println(lines(strings.NewReader("sync\nagain")))
}
