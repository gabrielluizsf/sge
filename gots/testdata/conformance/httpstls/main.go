// httpstls: an HTTPS server of httptest's, reached by the client that
// trusts its certificate, and refused by clients that do not, or that pin
// another one.
package main

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
)

func main() {
	srv := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello over %s, tls %v", r.Proto, r.TLS != nil)
	}))
	// The handshakes refused below are logged, with the time.
	srv.Config.ErrorLog = log.New(io.Discard, "", 0)
	srv.StartTLS()
	defer srv.Close()

	res, err := srv.Client().Get(srv.URL)
	if err != nil {
		panic(err)
	}
	b, _ := io.ReadAll(res.Body)
	res.Body.Close()
	fmt.Println(res.StatusCode, string(b))

	// Roots that do not hold the server's certificate.
	strict := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{RootCAs: x509.NewCertPool()}}}
	_, err = strict.Get(srv.URL)
	_, verify := errors.AsType[*tls.CertificateVerificationError](err)
	fmt.Println("empty roots refused:", err != nil, verify)

	// The right roots, and a pin that refuses.
	tr := srv.Client().Transport.(*http.Transport).Clone()
	tr.TLSClientConfig.VerifyPeerCertificate = func([][]byte, [][]*x509.Certificate) error { return errors.New("pinned elsewhere") }
	_, err = (&http.Client{Transport: tr}).Get(srv.URL)
	fmt.Println("pin refused:", err != nil)
}
