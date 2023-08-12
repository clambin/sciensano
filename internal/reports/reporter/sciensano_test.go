package reporter_test

import (
	"context"
	"github.com/clambin/go-common/taskmanager"
	"github.com/clambin/sciensano/internal/reports/datasource"
	"github.com/clambin/sciensano/internal/reports/datasource/mocks"
	"github.com/clambin/sciensano/internal/reports/reporter"
	mocks2 "github.com/clambin/sciensano/internal/reports/reporter/mocks"
	"github.com/clambin/sciensano/internal/reports/store"
	"github.com/clambin/sciensano/internal/sciensano"
	"github.com/clambin/sciensano/internal/sciensano/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/exp/slices"
	"golang.org/x/exp/slog"
	"net/http"
	"os"
	"sort"
	"testing"
	"time"
)

func TestSciensanoReporters(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	popStore := mocks2.NewFetcher(t)
	popStore.EXPECT().GetByRegion().Return(nil)
	popStore.EXPECT().GetByAgeBracket(mock.AnythingOfType("bracket.Bracket")).Return(1)
	popStore.EXPECT().WaitTillReady(mock.AnythingOfType("*context.timerCtx")).Return(nil)

	datasources := datasource.NewSciensanoDatastore("", 15*time.Second, http.DefaultClient, logger)
	casesFetcher := mocks.NewFetcher[sciensano.Cases](t)
	casesFetcher.EXPECT().GetLastModified(mock.AnythingOfType("*context.cancelCtx")).Return(time.Now(), nil)
	casesFetcher.EXPECT().Fetch(mock.AnythingOfType("*context.cancelCtx")).Return(testutil.Cases(), nil)
	datasources.Cases.Fetcher = casesFetcher

	hospFetcher := mocks.NewFetcher[sciensano.Hospitalisations](t)
	hospFetcher.EXPECT().GetLastModified(mock.AnythingOfType("*context.cancelCtx")).Return(time.Now(), nil)
	hospFetcher.EXPECT().Fetch(mock.AnythingOfType("*context.cancelCtx")).Return(testutil.Hospitalisations(), nil)
	datasources.Hospitalisations.Fetcher = hospFetcher

	mortFetcher := mocks.NewFetcher[sciensano.Mortalities](t)
	mortFetcher.EXPECT().GetLastModified(mock.AnythingOfType("*context.cancelCtx")).Return(time.Now(), nil)
	mortFetcher.EXPECT().Fetch(mock.AnythingOfType("*context.cancelCtx")).Return(testutil.Mortalities(), nil)
	datasources.Mortalities.Fetcher = mortFetcher

	testFetcher := mocks.NewFetcher[sciensano.TestResults](t)
	testFetcher.EXPECT().GetLastModified(mock.AnythingOfType("*context.cancelCtx")).Return(time.Now(), nil)
	testFetcher.EXPECT().Fetch(mock.AnythingOfType("*context.cancelCtx")).Return(testutil.TestResults(), nil)
	datasources.TestResults.Fetcher = testFetcher

	vaccFetcher := mocks.NewFetcher[sciensano.Vaccinations](t)
	vaccFetcher.EXPECT().GetLastModified(mock.AnythingOfType("*context.cancelCtx")).Return(time.Now(), nil)
	vaccFetcher.EXPECT().Fetch(mock.AnythingOfType("*context.cancelCtx")).Return(testutil.Vaccinations(), nil)
	datasources.Vaccinations.Fetcher = vaccFetcher

	mgr := taskmanager.New(datasources)

	s := store.Store{Logger: logger.With("component", "store")}
	reporters := reporter.NewSciensanoReporters(datasources, &s, popStore, logger)
	_ = mgr.Add(reporters...)

	ctx, cancel := context.WithCancel(context.Background())

	ch := make(chan error)
	go func() { ch <- mgr.Run(ctx) }()

	assert.Eventually(t, func() bool {
		keys := s.Keys()
		sort.Strings(keys)
		return slices.Equal(keys, []string{
			"cases-ByAgeGroup", "cases-ByProvince", "cases-ByRegion", "cases-Total",
			"hospitalisations-ByCategory", "hospitalisations-ByProvince", "hospitalisations-ByRegion", "hospitalisations-Total",
			"mortalities-ByAgeGroup", "mortalities-ByRegion", "mortalities-Total",
			"tests-ByCategory", "tests-Total",
			"vaccination-rate-Full-ByAgeGroup", "vaccination-rate-Full-ByRegion", "vaccination-rate-Partial-ByAgeGroup", "vaccination-rate-Partial-ByRegion",
			"vaccinations-ByAgeGroup", "vaccinations-ByManufacturer", "vaccinations-ByRegion", "vaccinations-ByVaccinationType", "vaccinations-Total",
		})
	}, time.Minute, time.Second)

	cancel()
	assert.ErrorIs(t, <-ch, context.Canceled)
}
