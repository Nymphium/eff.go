// chan_stack controls call stacks for coroutines.
// It provides push/pop operations
package chan_stack

import (
	"sync"
)

type ch = chan any

type t struct {
	chans []ch
	mu    sync.Mutex
}

// channel call stack for entering/returning goroutine.
var (
	Entering = &t{
		chans: make([]ch, 0),
		mu:    sync.Mutex{},
	}

	Returning = &t{
		chans: make([]ch, 0),
		mu:    sync.Mutex{},
	}
)

func (t *t) PushNew() ch {
	ch := make(ch)
	t.Push(ch)
	return ch
}

func (t *t) Push(ch ch) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.chans = append(t.chans, ch)
}

func (t *t) Pop() ch {
	t.mu.Lock()
	defer t.mu.Unlock()

	chans := t.chans
	len := len(chans)

	if len == 0 {
		return nil
	}

	ch := chans[len-1]
	t.chans = chans[:len-1]
	return ch
}

func (t *t) DeleteBySelf(ch ch) {
	t.mu.Lock()
	defer t.mu.Unlock()

	chans := t.chans
	len := len(chans)

	for i := 0; i < len; i++ {
		if chans[i] == ch {
			t.chans = append(chans[:i], chans[i+1:]...)
			return
		}
	}
}
