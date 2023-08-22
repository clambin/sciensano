package reporter_test

import (
	"context"
	"errors"
	"github.com/clambin/sciensano/v2/internal/reports/reporter"
	"github.com/clambin/sciensano/v2/internal/reports/reporter/mocks"
	"github.com/clambin/sciensano/v2/internal/reports/store"
	"github.com/clambin/sciensano/v2/internal/sciensano"
	"github.com/clambin/sciensano/v2/internal/sciensano/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"log/slog"
	"testing"
	"time"
)

func TestRater(t *testing.T) {
	f := mocks.NewPopulationFetcher(t)
	f.EXPECT().GetByRegion().Return(nil)
	f.EXPECT().WaitTillReady(mock.AnythingOfType("*context.timerCtx")).Return(nil)

	dataChCh := make(chan chan sciensano.Vaccinations)
	p := mocks.NewPublisher[sciensano.Vaccinations](t)
	p.EXPECT().Register(mock.AnythingOfType("chan sciensano.Vaccinations")).Run(func(ch chan sciensano.Vaccinations) {
		dataChCh <- ch
	})
	p.EXPECT().Unregister(mock.AnythingOfType("chan sciensano.Vaccinations"))

	s := &store.Store{Logger: slog.Default().With("component", "reportsStore")}

	r := reporter.ProRater{
		Name:     "vaccinations-rate-Full-ByRegion",
		Source:   p,
		PopStore: f,
		Mode:     sciensano.ByRegion,
		DoseType: sciensano.Full,
		Store:    s,
		Logger:   slog.Default().With("component", "rate"),
	}

	ctx, cancel := context.WithCancel(context.Background())

	ch2 := make(chan error)
	go func() {
		ch2 <- r.Run(ctx)
	}()

	dataCh := <-dataChCh
	dataCh <- testutil.Vaccinations()

	assert.Eventually(t, func() bool {
		_, err := r.Store.Get("vaccinations-rate-Full-ByRegion")
		return !errors.Is(err, store.ErrNotFound)
	}, time.Minute, time.Second)

	cancel()
	assert.ErrorIs(t, <-ch2, context.Canceled)
}
