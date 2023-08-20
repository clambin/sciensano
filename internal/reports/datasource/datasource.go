package datasource

import (
	"context"
	"github.com/clambin/go-common/taskmanager"
	"log/slog"
	"math/rand"
	"sync"
	"time"
)

type Fetcher[T any] interface {
	GetLastModified(ctx context.Context) (time.Time, error)
	Fetch(ctx context.Context) (T, error)
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
}

func (d *DataSource[T]) GetCurrentAge() time.Time {
	d.lock.Lock()
	defer d.lock.Unlock()
	return d.currentAge
}

func (d *DataSource[T]) Run(ctx context.Context) error {
	if err := d.fetchData(ctx); err != nil {
		d.Logger.Error("failed to collect data", "err", err)
	} else {
		d.sendData()
	}

	ticker := time.NewTicker(jitter(d.PollingInterval, 0.04, rand.Float64()))
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

func jitter(interval time.Duration, width float64, randFactor float64) time.Duration {
	total := float64(interval) * width
	lowest := float64(interval) - total/2
	j := total * randFactor
	return time.Duration(lowest + j)
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
		d.Logger.Info("new data found")
	}
	return err
}

func (d *DataSource[T]) sendData() {
	d.Publisher.Send(d.currentData, d.currentAge)
}
