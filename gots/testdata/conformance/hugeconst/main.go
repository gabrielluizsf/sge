// int64 arithmetic whose constants are beyond 2^53, which a JavaScript
// number cannot hold exactly, though the result is small: time does this
// on every date it works out.
package main

import (
	"fmt"
	"time"
)

const (
	far    int64 = -(292277022400*365.2425 + 306) * 86400
	unix   int64 = (1969*365 + 1969/4 - 1969/100 + 1969/400) * 86400
	toUnix       = far - unix
)

type abs uint64

func main() {
	sec := int64(1792843200)
	a := abs(sec + unix - far)
	fmt.Println(a, int64(a)+toUnix, int64(uint64(a))+far-unix == sec)
	x := int64(7)
	fmt.Println(x*86400+toUnix-toUnix, -(x - far) + far, (x<<60)>>58)

	t := time.Date(2026, 10, 24, 12, 0, 0, 0, time.UTC)
	fmt.Println(t.Format(time.RFC3339), t.Unix(), t.Weekday(), t.YearDay())
	p, err := time.Parse(time.RFC3339, "2020-01-02T03:04:05Z")
	fmt.Println(p, err, p.Add(36*time.Hour).Format(time.RFC1123))
	b, _ := t.MarshalJSON()
	fmt.Println(string(b), time.Unix(0, 0).UTC(), time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC).Unix())
}
