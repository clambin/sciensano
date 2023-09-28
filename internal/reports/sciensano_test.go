package reports_test

import (
	"context"
	"github.com/clambin/go-common/taskmanager"
	"github.com/clambin/sciensano/v2/internal/reports"
	"github.com/clambin/sciensano/v2/internal/reports/datasource"
	"github.com/clambin/sciensano/v2/internal/reports/reporter/mocks"
	"github.com/clambin/sciensano/v2/internal/reports/store"
	"github.com/clambin/sciensano/v2/internal/sciensano/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"log/slog"
	"net/http"
	"os"
	"slices"
	"testing"
	"time"
)

func TestSciensanoReporters(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	server := testutil.NewTestServer()
	defer server.Close()
	datasources := datasource.NewSciensanoDatastore(server.URL, 15*time.Second, http.DefaultClient, logger)
	mgr := taskmanager.New(datasources)

	s := store.Store{Logger: logger.With("component", "store")}

	popStore := mocks.NewPopulationFetcher(t)
	popStore.EXPECT().GetForRegion(mock.AnythingOfType("string")).Return(1)
	popStore.EXPECT().GetForAgeBracket(mock.AnythingOfType("bracket.Bracket")).Return(1)
	popStore.EXPECT().WaitTillReady(mock.AnythingOfType("*context.timerCtx")).Return(nil)

	reporters := reports.NewSciensanoReporters(datasources, &s, popStore, logger)
	_ = mgr.Add(reporters...)

	ctx, cancel := context.WithCancel(context.Background())

	ch := make(chan error)
	go func() { ch <- mgr.Run(ctx) }()

	assert.Eventually(t, func() bool {
		keys := s.Keys()
		slices.Sort(keys)
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
