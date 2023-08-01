package reporter_test

import (
	"context"
	"errors"
	"github.com/clambin/sciensano/internal/reports/datasource"
	mocks2 "github.com/clambin/sciensano/internal/reports/datasource/mocks"
	"github.com/clambin/sciensano/internal/reports/reporter"
	"github.com/clambin/sciensano/internal/reports/reporter/mocks"
	"github.com/clambin/sciensano/internal/reports/store"
	"github.com/clambin/sciensano/internal/sciensano"
	"github.com/clambin/sciensano/internal/sciensano/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/exp/slog"
	"testing"
	"time"
)

func TestRater(t *testing.T) {
	p := mocks.NewFetcher(t)
	p.EXPECT().GetByRegion().Return(nil)

	vaccinationsFetcher := mocks2.NewFetcher[sciensano.Vaccinations](t)
	vaccinationsFetcher.EXPECT().GetLastModified(mock.AnythingOfType("*context.cancelCtx")).Return(time.Now(), nil)
	vaccinationsFetcher.EXPECT().Fetch(mock.AnythingOfType("*context.cancelCtx")).Return(testutil.Vaccinations(), nil)

	s := &store.Store{Logger: slog.Default().With("component", "reportsStore")}

	d := &datasource.DataSource[sciensano.Vaccinations]{
		Fetcher:         vaccinationsFetcher,
		PollingInterval: time.Second,
		Logger:          slog.Default().With("datasource", "vaccinations"),
	}

	r := reporter.ProRater{
		Name:       "vaccinations-rate-Full-ByRegion",
		Source:     d,
		PopStore:   p,
		Mode:       sciensano.ByRegion,
		DoseType:   sciensano.Full,
		Accumulate: true,
		Store:      s,
		Logger:     slog.Default().With("component", "rate"),
	}

	ctx, cancel := context.WithCancel(context.Background())

	ch := make(chan error)
	go func() {
		ch <- d.Run(ctx)
	}()
	go func() {
		ch <- r.Run(ctx)
	}()

	assert.Eventually(t, func() bool {
		_, err := r.Store.Get("vaccinations-rate-Full-ByRegion")
		return !errors.Is(err, store.ErrNotFound)
	}, time.Minute, time.Second)

	cancel()
	assert.ErrorIs(t, <-ch, context.Canceled)
	assert.ErrorIs(t, <-ch, context.Canceled)
}
