package population

import (
	"context"
	"sync"
)

type Waiter struct {
	ready   bool
	waiters []chan struct{}
	lock    sync.RWMutex
}

func (w *Waiter) Ready() {
	w.lock.Lock()
	defer w.lock.Unlock()
	w.ready = true
	for _, ch := range w.waiters {
		close(ch)
	}
	w.waiters = nil
}

func (w *Waiter) WaitTillReady(ctx context.Context) error {
	if w.isReady() {
		return nil
	}

	ch := w.addWaiter()
	select {
	case <-ch:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (w *Waiter) isReady() bool {
	w.lock.Lock()
	defer w.lock.Unlock()
	return w.ready
}

func (w *Waiter) addWaiter() chan struct{} {
	w.lock.Lock()
	defer w.lock.Unlock()
	ch := make(chan struct{})
	w.waiters = append(w.waiters, ch)
	return ch
}
