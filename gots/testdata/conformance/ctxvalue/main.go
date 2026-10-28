// context.WithValue checks its key is comparable, through reflectlite.
package main

import (
	"context"
	"fmt"
)

type key struct{ name string }

func main() {
	ctx := context.WithValue(context.Background(), key{"user"}, "ana")
	fmt.Println(ctx.Value(key{"user"}), ctx.Value(key{"other"}))
	defer func() { fmt.Println("recovered:", recover()) }()
	context.WithValue(ctx, []int{1}, "x")
}
