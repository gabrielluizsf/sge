// What HTTPS and its friends lean on: hashes, HMAC, constant-time
// comparison and the system's randomness.
package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
)

func main() {
	sum := sha256.Sum256([]byte("sge"))
	fmt.Println(hex.EncodeToString(sum[:]))
	mac := hmac.New(sha256.New, []byte("key"))
	mac.Write([]byte("message"))
	m := mac.Sum(nil)
	fmt.Println(hex.EncodeToString(m))
	fmt.Println(hmac.Equal(m, m), hmac.Equal(m, sum[:]))
	fmt.Println(subtle.ConstantTimeCompare([]byte("abc"), []byte("abc")), subtle.ConstantTimeCompare([]byte("abc"), []byte("abd")))
	a, b := make([]byte, 32), make([]byte, 32)
	rand.Read(a)
	rand.Read(b)
	fmt.Println("random differs:", !hmac.Equal(a, b), len(rand.Text()))
}
