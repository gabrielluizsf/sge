//go:build gots_overlay

// Overlay for the system's roots, which under gots are those the runtime
// verifies with by default: Deno hands them to node:tls's
// rootCertificates, once --allow-sys=ca lets it. js/wasm has no files of
// roots for Go to read.

package x509

import (
	"errors"
	"syscall/js"
)

func loadSystemRoots() (*CertPool, error) {
	pems, err := runtimeRoots()
	if err != nil {
		return nil, err
	}
	pool := NewCertPool()
	for _, p := range pems {
		pool.AppendCertsFromPEM([]byte(p))
	}
	return pool, nil
}

// runtimeRoots is the runtime's roots, in PEM.
func runtimeRoots() (pems []string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("x509: the runtime's roots cannot be read: Deno hands them to programs run with --allow-sys=ca")
		}
	}()
	p := js.Global().Get("process")
	if p.IsUndefined() || p.Get("getBuiltinModule").IsUndefined() {
		return nil, errors.New("x509: this runtime has no roots to hand Go")
	}
	roots := p.Call("getBuiltinModule", "node:tls").Get("rootCertificates")
	for i := range roots.Length() {
		pems = append(pems, roots.Index(i).String())
	}
	return pems, nil
}
