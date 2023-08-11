package server

import (
	"context"
	"github.com/clambin/go-common/tabulator"
	grafanaJSONServer "github.com/clambin/grafana-json-server"
	"github.com/clambin/sciensano/internal/sciensano"
	"github.com/clambin/sciensano/server/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewSummaryMetric(t *testing.T) {
	metric, _ := newSummaryMetric(nil, "foo", []sciensano.SummaryColumn{sciensano.ByRegion, sciensano.ByAgeGroup})

	assert.Equal(t, "foo", metric.Label)
	assert.Equal(t, "foo", metric.Value)
	require.Len(t, metric.Payloads, 2)
	assert.Equal(t, "Summary", metric.Payloads[0].Name)
	assert.Len(t, metric.Payloads[0].Options, 2)
	assert.Equal(t, "Accumulate", metric.Payloads[1].Name)
	assert.Len(t, metric.Payloads[1].Options, 2)
}

func TestSummaryMetric_Query(t *testing.T) {
	s := mocks.NewReportsStorer(t)
	table := tabulator.New("A", "B")
	s.EXPECT().Get("foo-ByRegion").Return(table, nil)
	_, query := newSummaryMetric(s, "foo", []sciensano.SummaryColumn{sciensano.ByRegion, sciensano.ByAgeGroup})

	ctx := context.Background()

	testCases := []struct {
		name    string
		payload []byte
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "valid",
			payload: []byte(`{ "summary": "ByRegion", "accumulate": "no" }`),
			wantErr: assert.NoError,
		},
		{
			name:    "invalid accumulate",
			payload: []byte(`{ "summary": "ByRegion", "accumulate": "false" }`),
			wantErr: assert.Error,
		},
		{
			name:    "missing summary",
			payload: []byte(`{ "accumulate": "yes" }`),
			wantErr: assert.Error,
		},
		{
			name:    "invalid payload",
			payload: []byte(`not a json object`),
			wantErr: assert.Error,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			req := grafanaJSONServer.QueryRequest{
				Targets: []grafanaJSONServer.QueryRequestTarget{{Payload: tt.payload, Target: "foo"}},
				Range:   grafanaJSONServer.Range{},
			}

			resp, err := query.Query(ctx, "foo", req)
			tt.wantErr(t, err)
			if err != nil {
				return
			}
			assert.Equal(t, 1+len(table.GetColumns()), len(resp.(grafanaJSONServer.TableResponse).Columns))

		})
	}
}

func TestNewVaccinationDoseTypeMetric(t *testing.T) {
	metric, _ := newVaccinationDoseTypeMetric(nil, "foo", []sciensano.SummaryColumn{sciensano.ByRegion, sciensano.ByAgeGroup}, []sciensano.DoseType{sciensano.Partial})

	assert.Equal(t, "foo", metric.Label)
	assert.Equal(t, "foo", metric.Value)
	require.Len(t, metric.Payloads, 3)
	assert.Equal(t, "Summary", metric.Payloads[0].Name)
	assert.Len(t, metric.Payloads[0].Options, 2)
	assert.Equal(t, "DoseType", metric.Payloads[1].Name)
	assert.Len(t, metric.Payloads[1].Options, 1)
	assert.Equal(t, "Accumulate", metric.Payloads[2].Name)
	assert.Len(t, metric.Payloads[2].Options, 2)
}

func TestVaccinationDoseTypeMetric_Query(t *testing.T) {
	s := mocks.NewReportsStorer(t)
	table := tabulator.New("A", "B")
	s.EXPECT().Get("foo-Partial-ByRegion").Return(table, nil)

	_, query := newVaccinationDoseTypeMetric(s, "foo", []sciensano.SummaryColumn{sciensano.ByRegion, sciensano.ByAgeGroup}, []sciensano.DoseType{sciensano.Partial})

	ctx := context.Background()

	testCases := []struct {
		name    string
		payload []byte
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "valid",
			payload: []byte(`{ "summary": "ByRegion", "doseType": "Partial", "accumulate": "no" }`),
			wantErr: assert.NoError,
		},
		{
			name:    "invalid accumulate",
			payload: []byte(`{ "summary": "ByRegion", "doseType": "Partial", "accumulate": "false" }`),
			wantErr: assert.Error,
		},
		{
			name:    "missing summary",
			payload: []byte(`{ "accumulate": "yes" }`),
			wantErr: assert.Error,
		},
		{
			name:    "invalid payload",
			payload: []byte(`not a json object`),
			wantErr: assert.Error,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			req := grafanaJSONServer.QueryRequest{
				Targets: []grafanaJSONServer.QueryRequestTarget{{Payload: tt.payload, Target: "foo"}},
				Range:   grafanaJSONServer.Range{},
			}

			resp, err := query.Query(ctx, "foo", req)
			tt.wantErr(t, err)
			if err != nil {
				return
			}
			assert.Equal(t, 1+len(table.GetColumns()), len(resp.(grafanaJSONServer.TableResponse).Columns))

		})
	}
}
