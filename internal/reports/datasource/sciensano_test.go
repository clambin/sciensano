package datasource_test

import (
	"context"
	"github.com/clambin/sciensano/internal/reports/datasource"
	"github.com/clambin/sciensano/internal/reports/datasource/mocks"
	"github.com/clambin/sciensano/internal/sciensano"
	"github.com/clambin/sciensano/internal/sciensano/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"log/slog"
	"net/http"
	"testing"
	"time"
)

func TestNewSciensanoDatastore(t *testing.T) {
	s := datasource.NewSciensanoDatastore("", time.Second, http.DefaultClient, slog.Default())

	casesFetcher := mocks.NewFetcher[sciensano.Cases](t)
	casesFetcher.EXPECT().GetLastModified(mock.AnythingOfType("*context.cancelCtx")).Return(time.Now(), nil)
	casesFetcher.EXPECT().Fetch(mock.AnythingOfType("*context.cancelCtx")).Return(testutil.Cases(), nil)
	s.Cases.Fetcher = casesFetcher

	hospFetcher := mocks.NewFetcher[sciensano.Hospitalisations](t)
	hospFetcher.EXPECT().GetLastModified(mock.AnythingOfType("*context.cancelCtx")).Return(time.Now(), nil)
	hospFetcher.EXPECT().Fetch(mock.AnythingOfType("*context.cancelCtx")).Return(testutil.Hospitalisations(), nil)
	s.Hospitalisations.Fetcher = hospFetcher

	mortFetcher := mocks.NewFetcher[sciensano.Mortalities](t)
	mortFetcher.EXPECT().GetLastModified(mock.AnythingOfType("*context.cancelCtx")).Return(time.Now(), nil)
	mortFetcher.EXPECT().Fetch(mock.AnythingOfType("*context.cancelCtx")).Return(testutil.Mortalities(), nil)
	s.Mortalities.Fetcher = mortFetcher

	testFetcher := mocks.NewFetcher[sciensano.TestResults](t)
	testFetcher.EXPECT().GetLastModified(mock.AnythingOfType("*context.cancelCtx")).Return(time.Now(), nil)
	testFetcher.EXPECT().Fetch(mock.AnythingOfType("*context.cancelCtx")).Return(testutil.TestResults(), nil)
	s.TestResults.Fetcher = testFetcher

	vaccFetcher := mocks.NewFetcher[sciensano.Vaccinations](t)
	vaccFetcher.EXPECT().GetLastModified(mock.AnythingOfType("*context.cancelCtx")).Return(time.Now(), nil)
	vaccFetcher.EXPECT().Fetch(mock.AnythingOfType("*context.cancelCtx")).Return(testutil.Vaccinations(), nil)
	s.Vaccinations.Fetcher = vaccFetcher

	ch := make(chan sciensano.Vaccinations)
	s.Vaccinations.Register(ch)

	errCh := make(chan error)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		errCh <- s.Run(ctx)
	}()

	result := <-ch
	assert.NotEmpty(t, result)

	assert.Eventually(t, func() bool {
		return !s.Cases.GetCurrentAge().IsZero() &&
			!s.Hospitalisations.GetCurrentAge().IsZero() &&
			!s.Mortalities.GetCurrentAge().IsZero() &&
			!s.TestResults.GetCurrentAge().IsZero() &&
			!s.Vaccinations.GetCurrentAge().IsZero()
	}, time.Minute, time.Second)

	cancel()
	assert.ErrorIs(t, <-errCh, context.Canceled)
}
