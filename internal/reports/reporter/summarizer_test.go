package reporter_test

import (
	"context"
	"errors"
	"github.com/clambin/sciensano/internal/reports/datasource"
	"github.com/clambin/sciensano/internal/reports/datasource/mocks"
	"github.com/clambin/sciensano/internal/reports/reporter"
	"github.com/clambin/sciensano/internal/reports/store"
	"github.com/clambin/sciensano/internal/sciensano"
	"github.com/clambin/sciensano/internal/sciensano/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"log/slog"
	"testing"
	"time"
)

func TestSummarizer(t *testing.T) {
	casesFetcher := mocks.NewFetcher[sciensano.Cases](t)
	casesFetcher.EXPECT().GetLastModified(mock.AnythingOfType("*context.cancelCtx")).Return(time.Now(), nil)
	casesFetcher.EXPECT().Fetch(mock.AnythingOfType("*context.cancelCtx")).Return(testutil.Cases(), nil)

	d := &datasource.DataSource[sciensano.Cases]{
		Fetcher:         casesFetcher,
		PollingInterval: time.Second,
		Logger:          slog.Default().With("datasource", "cases"),
	}
	r := reporter.Summary[sciensano.Cases]{
		Name:   "cases-Total",
		Source: d,
		Mode:   sciensano.Total,
		Store:  &store.Store{Logger: slog.Default().With("component", "store")},
		Logger: slog.Default().With("reporter", "cases-Total"),
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
		_, err := r.Store.Get("cases-Total")
		return !errors.Is(err, store.ErrNotFound)
	}, time.Minute, time.Second)

	cancel()
	assert.ErrorIs(t, <-ch, context.Canceled)
	assert.ErrorIs(t, <-ch, context.Canceled)
}
