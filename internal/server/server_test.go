package server_test

import (
	"context"
	"fmt"
	"github.com/clambin/go-common/tabulator"
	gjson "github.com/clambin/grafana-json-server"
	"github.com/clambin/sciensano/v2/internal/sciensano"
	"github.com/clambin/sciensano/v2/internal/sciensano/testutil"
	"github.com/clambin/sciensano/v2/internal/server"
	"github.com/clambin/sciensano/v2/internal/server/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log/slog"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	store := makeStore(t)
	s := server.New(store, nil, slog.Default())
	ctx := context.Background()

	for target, handler := range s.Handlers {
		t.Run(target, func(t *testing.T) {
			payload := fmt.Sprintf(`{"summary":"%s", "accumulate": "yes"}`, sciensano.Total.String())
			if target == "vaccination-rate" {
				payload = fmt.Sprintf(`{"summary":"%s", "doseType": "%s", "accumulate": "yes"}`, sciensano.Total.String(), sciensano.Partial.String())
			}

			req := gjson.QueryRequest{Targets: []gjson.QueryRequestTarget{
				{Target: target, Payload: []byte(payload)},
			}, Range: gjson.Range{To: time.Now()}}

			resp, err := handler.Query(ctx, target, req)
			require.NoError(t, err)
			if target != "vaccination-rate" {
				assert.NotZero(t, len(resp.(gjson.TableResponse).Columns[0].Data.(gjson.TimeColumn)))
			}
		})
	}
}

func makeStore(t *testing.T) server.ReportsStore {
	s := mocks.NewReportsStore(t)

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
