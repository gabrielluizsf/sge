// cipheralias: ciphers work in place, as their docs allow, and refuse a
// destination that overlaps the source anywhere else.
package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
)

func main() {
	key := make([]byte, 16)
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	gcm, _ := cipher.NewGCM(block)
	nonce := make([]byte, gcm.NonceSize())
	msg := []byte("sixteen byte msg and some more")
	sealed := gcm.Seal(msg[:0], nonce, msg, nil) // in place
	fmt.Println(hex.EncodeToString(sealed))
	opened, err := gcm.Open(sealed[:0], nonce, sealed, nil)
	fmt.Println(string(opened), err)

	buf := make([]byte, 32)
	fmt.Println("xor in place", subtle.XORBytes(buf, buf, buf))
	fmt.Println("shifted", refused(func() { subtle.XORBytes(buf[1:17], buf[:16], buf[:16]) }))
	other := make([]byte, 32)
	fmt.Println("apart", subtle.XORBytes(other, buf, buf))
	fmt.Println("tail", refused(func() { subtle.XORBytes(buf[:16], buf[8:24], other[:16]) }))
}

func refused(f func()) (msg any) {
	defer func() { msg = recover() }()
	f()
	return nil
}
