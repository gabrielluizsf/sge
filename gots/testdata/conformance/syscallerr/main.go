// os.SyscallError, which net wraps its system calls' errors in.
package main

import (
	"errors"
	"fmt"
	"os"
	"syscall"
)

func main() {
	err := os.NewSyscallError("connect", syscall.ECONNREFUSED)
	var se *os.SyscallError
	fmt.Println(errors.As(err, &se), se.Syscall, errors.Is(err, syscall.ECONNREFUSED))
	fmt.Println(os.NewSyscallError("read", nil) == nil, se.Timeout())
	fmt.Println(errors.Is(os.ErrDeadlineExceeded, os.ErrDeadlineExceeded), os.ErrDeadlineExceeded)
}
