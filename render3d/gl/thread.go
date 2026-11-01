//go:build linux && amd64

package gl

import "runtime"

// thread is a goroutine locked to its OS thread, the one the device's
// context is current on, running there what it is given, in order.
type thread struct {
	jobs chan func()
	done chan struct{}
}

// startThread is a thread, running.
func startThread() *thread {
	t := &thread{jobs: make(chan func()), done: make(chan struct{})}
	go t.loop()
	return t
}

// loop runs the thread's jobs. It never unlocks: when it ends, its OS
// thread ends too, and the context current there with it.
func (t *thread) loop() {
	runtime.LockOSThread()
	for job := range t.jobs {
		job()
		t.done <- struct{}{}
	}
}

// do runs job on the thread and waits for it to be done.
func (t *thread) do(job func()) {
	t.jobs <- job
	<-t.done
}

// stop ends the thread, once its last job is done.
func (t *thread) stop() { close(t.jobs) }
