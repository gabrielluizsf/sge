// net/http/httptest's server, as tests start one.
package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
)

func main() {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		w.Write([]byte("not really a png"))
	}))
	defer srv.Close()
	resp, err := http.Get(srv.URL + "/pic.png")
	if err != nil {
		fmt.Println("get:", err)
		return
	}
	b, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	fmt.Println(resp.StatusCode, resp.Header.Get("Content-Type"), string(b))
	nf := httptest.NewServer(http.NotFoundHandler())
	resp, err = http.Get(nf.URL)
	fmt.Println(resp.StatusCode, err)
	resp.Body.Close()
	nf.Close()
}
