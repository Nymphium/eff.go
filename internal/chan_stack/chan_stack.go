// chan_stack controls call stacks for coroutines.
// It provides push/pop operations
package chan_stack

import (
	"sync"
)

type ch[A any] chan A

type t[A any] struct {
	chans []ch[A]
	mu    sync.Mutex
}

func New[A any]() *t[A] {
	return &t[A]{chans: make([]ch[A], 0),
		mu: sync.Mutex{},
	}
}

func (t *t[A]) Push(ch ch[A]) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.chans = append(t.chans, ch)
}

func (t *t[A]) Pop() ch[A] {
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

// DeleteBySelf deletes the specified channel from the stack.
func (t *t[A]) DeleteBySelf(ch ch[A]) {
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
