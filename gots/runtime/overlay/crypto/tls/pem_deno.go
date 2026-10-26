//go:build gots_overlay

// Certificates and keys as node:tls takes them: PEM.

package tls

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

// pemPair is cert's chain and private key in PEM.
func pemPair(cert *Certificate) (chain, key string, err error) {
	var b []byte
	for _, der := range cert.Certificate {
		b = append(b, pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})...)
	}
	k, err := x509.MarshalPKCS8PrivateKey(cert.PrivateKey)
	if err != nil {
		return "", "", fmt.Errorf("tls: the certificate's private key cannot be handed to node:tls: %w", err)
	}
	return string(b), string(pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: k})), nil
}

// anys is s as JavaScript takes an array.
func anys(s []string) []any {
	a := make([]any, len(s))
	for i, v := range s {
		a[i] = v
	}
	return a
}
