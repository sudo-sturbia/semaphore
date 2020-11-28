package semaphore

// queue is a simple queue used by the Semaphore to keep track of
// blocked objects.
type queue []*bool

// push adds an element to the end of the queue.
func (q *queue) push(b *bool) {
	*q = append(*q, b)
}

// pop returns the first element of the queue and removes it.
func (q *queue) pop() *bool {
	b := (*q)[0]
	*q = (*q)[1:]
	return b
}

// empty returns true if the queue is empty, false otherwise.
func (q *queue) empty() bool {
	return len(*q) == 0
}
