package datasource_test

import (
	"context"
	"github.com/clambin/sciensano/internal/reports/datasource"
	"github.com/clambin/sciensano/internal/reports/datasource/mocks"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"testing"
	"time"
)

func TestDataSource(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	f := mocks.NewFetcher[int](t)
	f.EXPECT().GetLastModified(ctx).Return(time.Now(), nil)
	f.EXPECT().Fetch(ctx).Return(100, nil)

	ds := datasource.DataSource[int]{
		Fetcher:         f,
		PollingInterval: time.Millisecond,
		Logger:          slog.Default().With("datasource", "test"),
	}

	dataCh := make(chan int)
	ds.Register(dataCh)
	errCh := make(chan error)

	go func() {
		errCh <- ds.Run(ctx)
	}()

	assert.Equal(t, 100, <-dataCh)

	ds.Unregister(dataCh)
	close(dataCh)

	cancel()
	assert.ErrorIs(t, <-errCh, context.Canceled)
}
