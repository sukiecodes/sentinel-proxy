// Essentially a state manager, defines various structs, and handles concurrency
package lb

import (
	"net/http/httputil"
	"net/url"
	"sync"
	"sync/atomic"
)

// Backend holds the data about a server we are proxying to 
type Backend struct {
	URL *url.URL 
	Alive bool 
	mux sync.RWMutex // RWMutex so multiple servers can read simultaneously but not write
	ReverseProxy *httputil.ReverseProxy 
}

// SetAlive updates the status of the backend safely using a Lock 
func (b *Backend) SetAlive(alive bool) {
	b.mux.Lock()
	b.Alive = alive 
	b.mux.Unlock()
}

// Returns whether or not backend server is alive
func (b *Backend) IsAlive() bool {
	b.mux.Lock()
	alive := b.Alive
	b.mux.Unlock() 
	return alive
}

// Holds information about reachable backends
type ServerPool struct {
	Backends []*Backend 
	current uint64 // used for atomic operations
}

// Atomically increases the counter and returns an index
func (s *ServerPool) NextIndex() int {
	// used to handle multiple people hitting the proxy at once 
	// counter increments correctly without skipping or duplicating numbers
	return int(atomic.AddUint64(&s.current, uint64(1)) % uint64(len(s.Backends)))
}

// Returns the next alive backendn server to send a request to
func (s *ServerPool) GetNextPeer() *Backend {
	next := s.NextIndex() // getting the next index via Round Robin 

	// start from "next" index and loop through all backends until we find an alive one
	l := len(s.Backends) + next 
	for i:=next; i<l; i++ {
		idx := i % len(s.Backends) // % to make sure index in range of number of backends
		if (s.Backends[idx].IsAlive()) {
			// if one is alive, keep track of it as "current"
			if (i != next) {
				atomic.StoreUint64(&s.current, uint64(idx))
			}
			return s.Backends[idx]
		}
	}
	return nil
}