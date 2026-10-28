// resolverdial: a Dialer with a Resolver of the program's own, Go's DNS
// client dialing a DNS server of the program's own, which tells it where
// game.test is: the dial goes where that server said.
package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"strconv"
)

func main() {
	app, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	go func() {
		c, err := app.Accept()
		if err == nil {
			c.Write([]byte("reached"))
			c.Close()
		}
	}()
	dns, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	go serveDNS(dns)

	r := &net.Resolver{PreferGo: true, Dial: func(ctx context.Context, _, _ string) (net.Conn, error) {
		return (&net.Dialer{}).DialContext(ctx, "tcp", dns.Addr().String())
	}}
	ips, err := r.LookupHost(context.Background(), "game.test")
	fmt.Println(ips, err)
	_, port, _ := net.SplitHostPort(app.Addr().String())
	c, err := (&net.Dialer{Resolver: r}).Dial("tcp4", net.JoinHostPort("game.test", port))
	if err != nil {
		fmt.Println("dial:", err)
		return
	}
	b, _ := io.ReadAll(c)
	fmt.Println(string(b))
	_, err = (&net.Dialer{Resolver: r}).Dial("tcp4", "nowhere.test:"+strconv.Itoa(1))
	fmt.Println(err != nil)
}

// serveDNS answers each query, over TCP, for game.test with 127.0.0.1, and
// for anything else with no answer.
func serveDNS(l net.Listener) {
	for {
		c, err := l.Accept()
		if err != nil {
			return
		}
		go func() {
			defer c.Close()
			for {
				var n uint16
				if binary.Read(c, binary.BigEndian, &n) != nil {
					return
				}
				q := make([]byte, n)
				if _, err := io.ReadFull(c, q); err != nil {
					return
				}
				a := answer(q)
				binary.Write(c, binary.BigEndian, uint16(len(a)))
				c.Write(a)
			}
		}()
	}
}

func answer(q []byte) []byte {
	// The question: names from byte 12, then its type and class.
	end := 12
	name := ""
	for q[end] != 0 {
		name += string(q[end+1:end+1+int(q[end])]) + "."
		end += 1 + int(q[end])
	}
	end++
	qtype := binary.BigEndian.Uint16(q[end:])
	end += 4
	found := name == "game.test." && qtype == 1
	rcode := uint16(0)
	if name != "game.test." {
		rcode = 3 // no such name
	}
	a := append([]byte{}, q[:2]...)
	ancount := uint16(0)
	if found {
		ancount = 1
	}
	a = binary.BigEndian.AppendUint16(a, 0x8580|rcode)
	a = binary.BigEndian.AppendUint16(a, 1)
	a = binary.BigEndian.AppendUint16(a, ancount)
	a = append(a, 0, 0, 0, 0)
	a = append(a, q[12:end]...)
	if found {
		a = append(a, 0xc0, 12, 0, 1, 0, 1, 0, 0, 0, 60, 0, 4, 127, 0, 0, 1)
	}
	return a
}
