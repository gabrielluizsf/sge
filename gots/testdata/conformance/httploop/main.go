// HTTP over loopback: Go's own server, and a client asking it for text,
// posting a body, following a redirect and refusing one.
package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /hello/{name}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Greeting", "yes")
		fmt.Fprintf(w, "hello, %s", r.PathValue("name"))
	})
	mux.HandleFunc("POST /echo", func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, "%s %d %s", r.Header.Get("Content-Type"), len(b), strings.ToUpper(string(b)))
	})
	mux.HandleFunc("/old", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/hello/redirected", http.StatusFound)
	})
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		fmt.Println("listen:", err)
		return
	}
	srv := &http.Server{Handler: mux}
	go srv.Serve(l)
	base := "http://" + l.Addr().String()

	resp, err := http.Get(base + "/hello/deno")
	if err != nil {
		fmt.Println("get:", err)
		return
	}
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	fmt.Println(resp.StatusCode, resp.Header.Get("X-Greeting"), string(body))

	resp, err = http.Post(base+"/echo", "text/plain", strings.NewReader("posted body"))
	if err != nil {
		fmt.Println("post:", err)
		return
	}
	body, _ = io.ReadAll(resp.Body)
	resp.Body.Close()
	fmt.Println(resp.StatusCode, string(body))

	resp, err = http.Get(base + "/missing")
	if err == nil {
		resp.Body.Close()
		fmt.Println("missing:", resp.StatusCode)
	}

	resp, err = http.Get(base + "/old")
	if err == nil {
		body, _ = io.ReadAll(resp.Body)
		resp.Body.Close()
		fmt.Println("followed:", resp.StatusCode, string(body))
	}

	// The client's own redirect policy decides, not fetch.
	stop := errors.New("not following")
	c := &http.Client{CheckRedirect: func(req *http.Request, via []*http.Request) error { return stop }}
	_, err = c.Get(base + "/old")
	fmt.Println("refused redirect:", errors.Is(err, stop))
	srv.Close()
}
