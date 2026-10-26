//go:build gots_overlay

// A client's failed check of the server's certificate, as Go reports it:
// a CertificateVerificationError around the x509 error Go would have
// found, so code that looks for one finds it.

package tls

import (
	"crypto/x509"
	"errors"
	"fmt"
)

func verifyError(err error) error {
	f, ok := err.(*jsFailure)
	if !ok {
		return err
	}
	var inner error
	switch f.code {
	case "UNABLE_TO_VERIFY_LEAF_SIGNATURE", "UNABLE_TO_GET_ISSUER_CERT", "UNABLE_TO_GET_ISSUER_CERT_LOCALLY",
		"DEPTH_ZERO_SELF_SIGNED_CERT", "SELF_SIGNED_CERT_IN_CHAIN", "CERT_UNTRUSTED", "CERT_SIGNATURE_FAILURE":
		inner = x509.UnknownAuthorityError{}
	case "CERT_HAS_EXPIRED", "CERT_NOT_YET_VALID":
		inner = x509.CertificateInvalidError{Reason: x509.Expired, Detail: fmt.Sprintf("%q", f.message)}
	case "ERR_TLS_CERT_ALTNAME_INVALID":
		inner = errors.New("x509: certificate is not valid for the server's name: " + fmt.Sprintf("%q", f.message))
	default:
		return err
	}
	return &CertificateVerificationError{Err: inner}
}
