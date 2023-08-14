package population

import (
	"context"
	"sync"
)

type Waiter struct {
	waiters []chan struct{}
	lock    sync.RWMutex
}

func (w *Waiter) Ready() {
	w.lock.Lock()
	defer w.lock.Unlock()
	for _, ch := range w.waiters {
		close(ch)
	}
	w.waiters = nil
}

func (w *Waiter) WaitTillReady(ctx context.Context) error {
	ch := w.addWaiter()
	select {
	case <-ch:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (w *Waiter) addWaiter() chan struct{} {
	w.lock.Lock()
	defer w.lock.Unlock()
	ch := make(chan struct{})
	w.waiters = append(w.waiters, ch)
	return ch
}
