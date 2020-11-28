package semaphore

import (
	"fmt"
	"sync"
	"testing"
)

// TestScheduling tests waiting and signaling.
func TestScheduling(t *testing.T) {
	// This tests scheduling using a semaphore. There are two functions that
	// execute on two different threads. The idea is that the first thread
	// acquire a lock and is then blocked until the other thread is called
	// to execute. After which the first thread is signalled to continue execution.
	//
	// If the semaphore works correctly, the first thread executes first
	// despite being hypothetically blocked until the second is executed.
	// This happens because the first thread is using the semaphore as a
	// mutex. So, as long as the first is blocked, the second is also
	// blocked.
	wg := new(sync.WaitGroup)
	wg.Add(2)

	semaphore := New(1)
	got := ""

	waiter := new(bool)
	*waiter = true

	signaller := new(bool)
	*signaller = false

	go func(signaller, waiter *bool) {
		semaphore.Wait()
		*signaller = true
		for *waiter {
		}

		got += "1"
		semaphore.Signal()
		wg.Done()
	}(signaller, waiter)

	for !*signaller {
	}

	go func() {
		semaphore.Wait()
		got += "2"
		semaphore.Signal()
		wg.Done()
	}()
	*waiter = false
	wg.Wait()

	want := "12"
	if got != want {
		t.Errorf("scheduling failed: want %s, got %s", want, got)
	}
}

// TestSchedulingMany runs a test similar to TestScheduling on a larger
// number of threads.
func TestSchedulingMany(t *testing.T) {
	threads := 100
	wg := new(sync.WaitGroup)
	wg.Add(threads)

	got := ""
	semaphore := New(0)

	for i := 0; i < threads; i++ {
		go func(i int) {
			semaphore.Wait()
			got = fmt.Sprintf("%s%d", got, i)
			semaphore.Signal()

			wg.Done()
		}(i)

		for semaphore.counter != -1*(i+1) {
		}
	}

	semaphore.Signal()
	wg.Wait()

	want := ""
	for i := 0; i < threads; i++ {
		want = fmt.Sprintf("%s%d", want, i)
	}

	if want != got {
		t.Errorf("unexpected gotult: want %s, got %s", want, got)
	}
}
