//go:build gots_overlay

// Overlay for the system's roots, which under gots are the runtime's: the
// store Deno or Node verifies with, which crypto/tls hands TLS when a
// pool says so. Go itself verifies nothing against them.

package x509

import (
	"encoding/pem"
	"errors"
)

func loadSystemRoots() (*CertPool, error) { return &CertPool{systemPool: true}, nil }

// GotsPEMs is what crypto/tls hands the runtime's TLS of s: its
// certificates in PEM, and whether the system's roots are among them.
// A certificate added with a constraint is refused, as the runtime
// cannot check it.
func GotsPEMs(s *CertPool) (pems []string, system bool, err error) {
	for i := range s.len() {
		c, constraint, err := s.cert(i)
		if err != nil {
			return nil, false, err
		}
		if constraint != nil {
			return nil, false, errors.New("x509: a certificate added with a constraint cannot be checked by the runtime's TLS")
		}
		pems = append(pems, string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: c.Raw})))
	}
	return pems, s.systemPool, nil
}
