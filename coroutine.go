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

func (c *Coroutine) Switch() {
	if c.next != nil {
		c.next()
	} else {
		c.yield(struct{}{})
	}
}

func (c *Coroutine) Release() {
	if c != nil && (c.next != nil || c.yield != nil) { // prevent double release
		_release_coroutine(c)
	}
}

// NewCoroutine creates a new coroutine with the given function.
func NewCoroutine(f func(*Coroutine)) *Coroutine {
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
