package main

import "sync"

// Counter counts the number of matches crawled in a concurrency-safe manner
type Counter struct {
	Total int32
	mux   sync.Mutex
}

// NewCounter returns a new GameCounter object and initializes the number of games crawled to 0
func NewCounter() *Counter {
	return &Counter{
		Total: 0,
	}
}

// Reset sets the internal counter back to 0
func (c *Counter) Reset() {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.Total = 0
}

// Inc increments the internal counter concurrency-safely
func (c *Counter) Inc() {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.Total++
}

// GetCount returns the current Count
func (c *Counter) GetCount() int32 {
	c.mux.Lock()
	defer c.mux.Unlock()
	return c.Total
}
