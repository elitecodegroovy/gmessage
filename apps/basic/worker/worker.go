// Package work manages a pool of goroutines to perform work.
package worker

import "sync"

// Worker must be implemented by types that want to use
// the worker pool.
type Worker interface {
	Task()
}

// Pool provides a pool of goroutines that can execute any Worker
// tasks that are submitted.
type Pool struct {
	worker chan Worker
	wg     sync.WaitGroup
}

// New creates a new worker pool.
func New(maxGoroutines int) *Pool {
	p := Pool{
		worker: make(chan Worker),
	}

	p.wg.Add(maxGoroutines)
	for i := 0; i < maxGoroutines; i++ {
		go func() {
			for w := range p.worker {
				w.Task()
			}
			p.wg.Done()
		}()
	}

	return &p
}

// Run submits worker to the pool.
func (p *Pool) Run(w Worker) {
	p.worker <- w
}

// Shutdown waits for all the goroutines to shutdown.
func (p *Pool) Shutdown() {
	close(p.worker)
	p.wg.Wait()
}
