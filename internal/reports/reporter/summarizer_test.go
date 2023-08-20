package reporter_test

import (
	"context"
	"errors"
	"github.com/clambin/sciensano/internal/reports/reporter"
	"github.com/clambin/sciensano/internal/reports/reporter/mocks"
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
	dataChCh := make(chan chan sciensano.Cases)

	p := mocks.NewPublisher[sciensano.Cases](t)
	p.EXPECT().Register(mock.AnythingOfType("chan sciensano.Cases")).Run(func(ch chan sciensano.Cases) {
		dataChCh <- ch
	})
	p.EXPECT().Unregister(mock.AnythingOfType("chan sciensano.Cases"))

	r := reporter.Summary[sciensano.Cases]{
		Name:   "cases-Total",
		Source: p,
		Mode:   sciensano.Total,
		Store:  &store.Store{Logger: slog.Default().With("component", "store")},
		Logger: slog.Default().With("reporter", "cases-Total"),
	}
	ctx, cancel := context.WithCancel(context.Background())

	ch := make(chan error)
	go func() {
		ch <- r.Run(ctx)
	}()

	dataCh := <-dataChCh
	dataCh <- testutil.Cases()

	assert.Eventually(t, func() bool {
		_, err := r.Store.Get("cases-Total")
		return !errors.Is(err, store.ErrNotFound)
	}, time.Minute, time.Second)

	cancel()
	assert.ErrorIs(t, <-ch, context.Canceled)
}
