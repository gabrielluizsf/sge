//go:build gots_overlay

// The syscall.RawConn a Control hook is given: Deno's sockets have no
// file descriptor to hand it, so one that asks for it is told so, and the
// dial or listen fails rather than go on without what it wanted set.

package net

import "syscall"

type denoRawConn struct{}

func (denoRawConn) Control(f func(fd uintptr)) error { return syscall.ENOTSUP }

func (denoRawConn) Read(f func(fd uintptr) bool) error { return syscall.ENOTSUP }

func (denoRawConn) Write(f func(fd uintptr) bool) error { return syscall.ENOTSUP }
