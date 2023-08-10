package server

import (
	"context"
	"fmt"
	"github.com/clambin/go-common/set"
	"github.com/clambin/go-common/tabulator"
	grafanaJSONServer "github.com/clambin/grafana-json-server"
	"github.com/clambin/sciensano/internal/sciensano"
	"github.com/clambin/sciensano/internal/sciensano/testutil"
	"github.com/clambin/sciensano/server/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestSummaryByDoseTypeHandler_Query(t *testing.T) {
	targets := []struct {
		sciensanoType string
		summaryTypes  set.Set[sciensano.SummaryColumn]
		doseType      sciensano.DoseType
	}{
		{sciensanoType: "vaccinations-rate", summaryTypes: sciensano.VaccinationsValidSummaryModes(), doseType: sciensano.Partial},
		{sciensanoType: "vaccinations-rate", summaryTypes: sciensano.VaccinationsValidSummaryModes(), doseType: sciensano.Full},
	}

	type summarizer interface {
		Summarize(summaryColumn sciensano.SummaryColumn) (*tabulator.Tabulator, error)
	}

	for _, target := range targets {
		var records summarizer
		var summaryTypes set.Set[sciensano.SummaryColumn]
		switch target.sciensanoType {
		case "vaccinations":
			records = testutil.Vaccinations()
			summaryTypes = sciensano.VaccinationsValidSummaryModes()
		}

		for _, summaryType := range summaryTypes.List() {
			t.Run(target.sciensanoType+"-"+summaryType.String(), func(t *testing.T) {
				s := mocks.NewReportsStorer(t)
				report, _ := records.Summarize(summaryType)
				expectedColumns := 1 + len(report.GetColumns())
				s.EXPECT().Get(target.sciensanoType+"-"+summaryType.String()+"-"+target.doseType.String()).Return(report, nil).Once()

				r := SummaryByDoseTypeHandler{ReportsStore: s}

				req := grafanaJSONServer.QueryRequest{
					Targets: []grafanaJSONServer.QueryRequestTarget{{
						Target:  target.sciensanoType,
						Payload: []byte(fmt.Sprintf(`{ "summary": "%s", "accumulate": "no" }`, summaryType.String())),
					}},
					Range: grafanaJSONServer.Range{From: time.Now().Add(-24 * time.Hour)},
				}

				resp, err := r.Query(context.Background(), target.sciensanoType, req)
				require.NoError(t, err)
				assert.Len(t, resp.(grafanaJSONServer.TableResponse).Columns, expectedColumns)
			})
		}
	}
}
