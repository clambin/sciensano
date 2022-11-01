package fetcher

import (
	"context"
	"github.com/clambin/sciensano/apiclient"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/semaphore"
	"time"
)

type Limiter struct {
	Fetcher
	limit *semaphore.Weighted
}

var _ Fetcher = &Limiter{}

func NewLimiter(fetcher Fetcher, maxParallel int64) *Limiter {
	return &Limiter{
		Fetcher: fetcher,
		limit:   semaphore.NewWeighted(maxParallel),
	}
}

func (l *Limiter) Fetch(ctx context.Context, dataType int) (results []apiclient.APIResponse, err error) {
	log.Debugf("limiter-fetch: attempting to call API for %s", l.DataTypes()[dataType])
	err = l.limit.Acquire(ctx, 1)
	if err != nil {
		return
	}
	log.Debugf("limiter-fetch: calling API for %s", l.DataTypes()[dataType])

	results, err = l.Fetcher.Fetch(ctx, dataType)

	l.limit.Release(1)
	log.Debugf("limiter-fetch: API called for %s", l.DataTypes()[dataType])
	return
}

func (l *Limiter) GetLastUpdated(ctx context.Context, dataType int) (lastUpdated time.Time, err error) {
	log.Debugf("limiter-getlastupdates: attempting to call API for %s", l.DataTypes()[dataType])
	err = l.limit.Acquire(ctx, 1)
	if err != nil {
		return
	}
	log.Debugf("limiter-getlastupdates: calling API for %s", l.DataTypes()[dataType])

	lastUpdated, err = l.Fetcher.GetLastUpdated(ctx, dataType)

	l.limit.Release(1)
	log.Debugf("limiter-getlastupdates: API called for %s", l.DataTypes()[dataType])
	return
}
