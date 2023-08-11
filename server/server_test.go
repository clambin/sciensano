package server_test

import (
	"context"
	"fmt"
	"github.com/clambin/go-common/tabulator"
	grafanaJSONServer "github.com/clambin/grafana-json-server"
	"github.com/clambin/sciensano/internal/sciensano"
	"github.com/clambin/sciensano/internal/sciensano/testutil"
	"github.com/clambin/sciensano/server"
	"github.com/clambin/sciensano/server/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/slog"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	store := makeStore(t)
	s := server.New(store, slog.Default())
	ctx := context.Background()

	for target, handler := range s.Handlers {
		t.Run(target, func(t *testing.T) {
			payload := fmt.Sprintf(`{"summary":"%s", "accumulate": "yes"}`, sciensano.Total.String())
			if target == "vaccination-rate" {
				payload = fmt.Sprintf(`{"summary":"%s", "doseType": "%s", "accumulate": "yes"}`, sciensano.Total.String(), sciensano.Partial.String())
			}

			req := grafanaJSONServer.QueryRequest{Targets: []grafanaJSONServer.QueryRequestTarget{
				{Target: target, Payload: []byte(payload)},
			}, Range: grafanaJSONServer.Range{To: time.Now()}}

			resp, err := handler.Query(ctx, target, req)
			require.NoError(t, err)
			if target != "vaccination-rate" {
				assert.NotZero(t, len(resp.(grafanaJSONServer.TableResponse).Columns[0].Data.(grafanaJSONServer.TimeColumn)))
			}
		})
	}
}

func makeStore(t *testing.T) server.ReportsStorer {
	s := mocks.NewReportsStorer(t)

	cases, _ := testutil.Cases().Summarize(sciensano.Total)
	s.EXPECT().Get("cases-Total").Return(cases, nil)
	mortalities, _ := testutil.Mortalities().Summarize(sciensano.Total)
	s.EXPECT().Get("mortalities-Total").Return(mortalities, nil)
	hospitalisations, _ := testutil.Hospitalisations().Summarize(sciensano.Total)
	s.EXPECT().Get("hospitalisations-Total").Return(hospitalisations, nil)
	tests, _ := testutil.Cases().Summarize(sciensano.Total)
	s.EXPECT().Get("tests-Total").Return(tests, nil)
	vaccinations, _ := testutil.Vaccinations().Summarize(sciensano.Total)
	s.EXPECT().Get("vaccinations-Total").Return(vaccinations, nil)
	s.EXPECT().Get("vaccination-rate-Partial-Total").Return(tabulator.New(), nil)
	return s
}
