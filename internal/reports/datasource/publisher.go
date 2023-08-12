package datasource

import (
	"sync"
	"time"
)

type Publisher[T any] struct {
	lock    sync.RWMutex
	clients map[chan T]time.Time
}

func (p *Publisher[T]) Register(ch chan T) {
	p.lock.Lock()
	defer p.lock.Unlock()
	if p.clients == nil {
		p.clients = make(map[chan T]time.Time)
	}
	p.clients[ch] = time.Time{}
}

func (p *Publisher[T]) Unregister(ch chan T) {
	p.lock.Lock()
	defer p.lock.Unlock()
	delete(p.clients, ch)
}

/*
func (p *Publisher[T]) Send(value T) {
	p.lock.RLock()
	defer p.lock.RUnlock()

	currentAge := time.Now()
	for ch, lastSent := range p.clients {
		if lastSent.Before(currentAge) {
			ch <- value
			p.clients[ch] = currentAge
		}
	}

}
*/
