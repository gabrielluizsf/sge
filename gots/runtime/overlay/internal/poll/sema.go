//go:build gots_overlay

// Overlay for the semaphore internal/poll's FD lock sleeps on. With one
// thread and the net overlay doing its own I/O, the lock is never
// contended: acquiring it when it is held is a bug, not a wait.

package poll

func runtime_Semacquire(sema *uint32) {
	if *sema == 0 {
		panic("gots: internal/poll: FD lock contended")
	}
	*sema--
}

func runtime_Semrelease(sema *uint32) { *sema++ }
