
package server

// We wrap to hold onto optional items for /connz.
type closedClient struct {
	ConnInfo
	subs []string
	user string
}

// Fixed sized ringbuffer for closed connections.
type closedRingBuffer struct {
	total uint64
	conns []*closedClient
}

// Create a new ring buffer with at most max items.
func newClosedRingBuffer(max int) *closedRingBuffer {
	rb := &closedRingBuffer{}
	rb.conns = make([]*closedClient, max)
	return rb
}

// Adds in a new closed connection. If there is no more room,
// remove the oldest.
func (rb *closedRingBuffer) append(cc *closedClient) {
	rb.conns[rb.next()] = cc
	rb.total++
}

func (rb *closedRingBuffer) next() int {
	return int(rb.total % uint64(cap(rb.conns)))
}

func (rb *closedRingBuffer) len() int {
	if rb.total > uint64(cap(rb.conns)) {
		return cap(rb.conns)
	}
	return int(rb.total)
}

func (rb *closedRingBuffer) totalConns() uint64 {
	return rb.total
}

// This will not be sorted. Will return a copy of the list
// which recipient can modify. If the contents of the client
// itself need to be modified, meaning swapping in any optional items,
// a copy should be made. We could introduce a new lock and hold that
// but since we return this list inside monitor which allows programatic
// access, we do not know when it would be done.
func (rb *closedRingBuffer) closedClients() []*closedClient {
	dup := make([]*closedClient, rb.len())
	if rb.total <= uint64(cap(rb.conns)) {
		copy(dup, rb.conns[:rb.len()])
	} else {
		first := rb.next()
		next := cap(rb.conns) - first
		copy(dup, rb.conns[first:])
		copy(dup[next:], rb.conns[:next])
	}
	return dup
}
