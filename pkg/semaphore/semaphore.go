// package semaphore contains an implementation of a counting semaphore.
// The semaphore uses FIFO for scheduling. While the implementation
// provided here works perfectly well, it's very inefficient as it employs
// busy waiting for blocking.
package semaphore

import (
	"sync/atomic"
)

// Semaphore is a simple implementation of a counting semaphore.
type Semaphore struct {
	counter int
	q       queue

	sync *int32
}

// New creates and returns a new semaphore.
func New(count int) *Semaphore {
	i := int32(0)
	return &Semaphore{
		counter: count,
		q:       make([]*bool, 0),
		sync:    &i,
	}
}

// Wait signals to the semaphore that a thread is entering a critical
// section. Depending on the state of the semaphore, the thread is either
// blocked or continued. A thread remains in a blocked state until another
// thread calls semaphore.Signal.
func (s *Semaphore) Wait() {
	for atomic.SwapInt32(s.sync, 1) == 1 {
	}
	s.counter--
	if s.counter < 0 {
		waiter := true
		s.q.push(&waiter)
		*s.sync = 0

		for waiter {
		}
		return
	}
	*s.sync = 0
}

// Signal signals the semaphore to add to it's counter. This most likely means
// that the thread that acquired the semaphore has finished execution of its
// critical section.
//
// When Signal is called, the semaphore checks if there're any blocked threads.
// If so, the oldest blocked thread is signaled to unblock and enter its critical
// section. Otherwise, the semaphore's counter is updated and the first thread
// to call Signal afterwards doesn't block.
func (s *Semaphore) Signal() {
	for atomic.SwapInt32(s.sync, 1) == 1 {
	}
	s.counter++
	if s.counter <= 0 {
		*s.q.pop() = false
	}
	*s.sync = 0
}
