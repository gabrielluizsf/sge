// tlsconfig: where the runtime has no sockets, a browser, a Transport with
// a TLSClientConfig fails rather than go through fetch, which would leave
// its roots, pins and client certificates aside; one without it fetches.
package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"strings"
)

func main() {
	pinned := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{
		VerifyPeerCertificate: func([][]byte, [][]*x509.Certificate) error { return nil },
	}}}
	_, err := pinned.Get("https://example.test/")
	fmt.Println("pinned fetched:", err == nil || strings.Contains(err.Error(), "fetch"))
	_, err = http.Get("https://example.test/")
	fmt.Println("plain:", err)
}
