// chan_stack controls global singleton call stack for coroutines.
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

// Global, singleton chanStack object.
var cs = &t{
	chans: make([]ch, 0),
	mu:    sync.Mutex{},
}

func PushNew() ch {
	ch := make(ch)
	Push(ch)
	return ch
}

func Push(ch ch) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	cs.chans = append(cs.chans, ch)
}

func Pop() ch {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	chans := cs.chans
	len := len(chans)

	if len == 0 {
		return nil
	}

	ch := chans[len-1]
	cs.chans = chans[:len-1]
	return ch
}

func DeleteBySelf(ch ch) {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	chans := cs.chans
	len := len(chans)

	for i := 0; i < len; i++ {
		if chans[i] == ch {
			cs.chans = append(chans[:i], chans[i+1:]...)
			return
		}
	}
}
