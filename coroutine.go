package zu

import "iter"

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

func NewCoroutine(f func(*Coroutine)) *Coroutine {
	next, _ := iter.Pull(func(yield func(struct{}) bool) {
		f(&Coroutine{
			yield: yield,
		})
	})
	return &Coroutine{
		next: next,
	}
}
