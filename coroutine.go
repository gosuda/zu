package zu

import (
	"iter"
	"sync"
)

var _coroutine_pool sync.Pool = sync.Pool{
	New: func() interface{} {
		return &Coroutine{}
	},
}

func _acquire_coroutine() *Coroutine {
	return _coroutine_pool.Get().(*Coroutine)
}

func _release_coroutine(c *Coroutine) {
	c.next = nil
	c.yield = nil
	c.Release()
}

type Coroutine struct {
	next  func() (struct{}, bool)
	yield func(struct{}) bool
}

// Switch switches the control between the coroutine and the caller.
// If the coroutine is paused, it will resume execution.
func (c *Coroutine) Switch() {
	if c.next != nil {
		c.next()
	} else {
		c.yield(struct{}{})
	}
}

// Release releases the coroutine object back to the pool.
func (c *Coroutine) Release() {
	if c != nil && (c.next != nil || c.yield != nil) { // prevent double release
		_release_coroutine(c)
	}
}

// NewCoroutine creates a new coroutine that executes the provided function.
// The function receives the coroutine instance as a parameter, allowing it to
// control its own execution flow.
//
// The coroutine is initially in a paused state and can be started by calling the Switch method.
// The function can pass control back to the caller by calling the Switch method.
//
// Example:
//
//	c := NewCoroutine(func(c *Coroutine) {
//	    // Do some work
//	    c.Switch()
//	    // Continue after being resumed
//	})
//
//	c.Switch() // This will start the coroutine
func NewCoroutine(f func(c *Coroutine)) *Coroutine {
	next, _ := iter.Pull(func(yield func(struct{}) bool) {
		c := _acquire_coroutine()
		c.yield = yield
		c.next = nil
		f(c)
		c.Release()
	})

	c := _acquire_coroutine()
	c.next = next
	c.yield = nil
	return c
}
