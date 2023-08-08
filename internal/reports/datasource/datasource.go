package datasource

import (
	"context"
	"github.com/clambin/go-common/taskmanager"
	"golang.org/x/exp/slog"
	"sync"
	"time"
)

//go:generate mockery --name Fetcher --with-expecter=true
type Fetcher[T any] interface {
	GetLastModified(ctx context.Context) (time.Time, error)
	Fetch(ctx context.Context) (T, error)
	GetTarget() string
}

var _ taskmanager.Task = &DataSource[int]{}

type DataSource[T any] struct {
	Publisher[T]
	Fetcher         Fetcher[T]
	PollingInterval time.Duration
	Logger          *slog.Logger
	currentData     T
	currentAge      time.Time
	lock            sync.RWMutex
	//clients         map[chan T]time.Time
}

func (d *DataSource[T]) GetCurrentAge() time.Time {
	d.lock.Lock()
	defer d.lock.Unlock()
	return d.currentAge
}

func (d *DataSource[T]) Run(ctx context.Context) error {
	if err := d.fetchData(ctx); err != nil {
		d.Logger.Error("failed to collect data", "err", err)
	}
	d.sendData()

	ticker := time.NewTicker(d.PollingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := d.fetchData(ctx); err != nil {
				d.Logger.Error("failed to collect data", "err", err)
				continue
			}
			d.sendData()
		}
	}
}

func (d *DataSource[T]) fetchData(ctx context.Context) error {
	d.lock.Lock()
	defer d.lock.Unlock()
	timestamp, err := d.Fetcher.GetLastModified(ctx)
	if !timestamp.After(d.currentAge) || err != nil {
		return err
	}
	data, err := d.Fetcher.Fetch(ctx)
	if err == nil {
		d.currentData = data
		d.currentAge = timestamp
	}
	return err
}

func (d *DataSource[T]) sendData() {
	d.Publisher.lock.RLock()
	defer d.Publisher.lock.RUnlock()
	//d.lock.Lock()
	//defer d.lock.Unlock()

	for ch, lastSent := range d.Publisher.clients {
		if lastSent.Before(d.currentAge) {
			d.Logger.Debug("new data found")
			ch <- d.currentData
			d.Publisher.clients[ch] = d.currentAge
		}
	}
}
