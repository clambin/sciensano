package reporter

import (
	"github.com/clambin/sciensano/demographics/bracket"
	"github.com/clambin/sciensano/internal/reports/reporter/mocks"
	"github.com/clambin/sciensano/internal/reports/store"
	"github.com/clambin/sciensano/internal/sciensano"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/slog"
	"testing"
	"time"
)

func TestProRater_Prorate(t *testing.T) {
	l := slog.Default()

	f := mocks.NewFetcher(t)
	f.EXPECT().GetByRegion().Return(map[string]int{"Flanders": 10, "Wallonia": 5, "Brussels": 1})
	f.EXPECT().GetByAgeBracket(bracket.Bracket{Low: 20, High: 29}).Return(10)
	f.EXPECT().GetByAgeBracket(bracket.Bracket{Low: 30, High: 39}).Return(5)
	f.EXPECT().GetByAgeBracket(bracket.Bracket{Low: 40, High: 49}).Return(1)

	ts := time.Date(2023, time.August, 7, 0, 0, 0, 0, time.UTC)
	vaccinations := sciensano.Vaccinations{
		{TimeStamp: sciensano.TimeStamp{Time: ts}, Region: "Flanders", AgeGroup: "20-29", Dose: sciensano.Partial, Count: 100},
		{TimeStamp: sciensano.TimeStamp{Time: ts}, Region: "Wallonia", AgeGroup: "30-39", Dose: sciensano.Partial, Count: 50},
		{TimeStamp: sciensano.TimeStamp{Time: ts}, Region: "Brussels", AgeGroup: "40-49", Dose: sciensano.Partial, Count: 10},
		{TimeStamp: sciensano.TimeStamp{Time: ts}, Region: "Flanders", Dose: sciensano.Full, Count: 10},
		{TimeStamp: sciensano.TimeStamp{Time: ts}, Region: "Wallonia", Dose: sciensano.SingleDose, Count: 5},
		{TimeStamp: sciensano.TimeStamp{Time: ts}, Region: "Brussels", Dose: sciensano.Full, Count: 1},
		{TimeStamp: sciensano.TimeStamp{Time: ts.Add(24 * time.Hour)}, Region: "Flanders", Dose: sciensano.Full, Count: 10},
		{TimeStamp: sciensano.TimeStamp{Time: ts.Add(24 * time.Hour)}, Region: "Wallonia", Dose: sciensano.SingleDose, Count: 5},
		{TimeStamp: sciensano.TimeStamp{Time: ts.Add(24 * time.Hour)}, Region: "Brussels", Dose: sciensano.Full, Count: 1},
		{TimeStamp: sciensano.TimeStamp{Time: ts.Add(48 * time.Hour)}, Region: "Flanders", Dose: sciensano.Full, Count: 10},
		{TimeStamp: sciensano.TimeStamp{Time: ts.Add(48 * time.Hour)}, Region: "Wallonia", Dose: sciensano.SingleDose, Count: 5},
		{TimeStamp: sciensano.TimeStamp{Time: ts.Add(48 * time.Hour)}, Region: "Brussels", Dose: sciensano.Full, Count: 1},
	}

	testCases := []struct {
		name        string
		mode        sciensano.SummaryColumn
		doseType    sciensano.DoseType
		wantColumns []string
		wantValues  []float64
	}{
		{
			name:        "Partial-ByRegion",
			mode:        sciensano.ByRegion,
			doseType:    sciensano.Partial,
			wantColumns: []string{"Brussels", "Flanders", "Wallonia"},
			wantValues:  []float64{10},
		},
		{
			name:        "Full-ByRegion",
			mode:        sciensano.ByRegion,
			doseType:    sciensano.Full,
			wantColumns: []string{"Brussels", "Flanders", "Wallonia"},
			wantValues:  []float64{1, 1, 1},
		},
		{
			name:        "Partial-ByAgeGroup",
			mode:        sciensano.ByAgeGroup,
			doseType:    sciensano.Partial,
			wantColumns: []string{"20-29", "30-39", "40-49"},
			wantValues:  []float64{10},
		},
		{
			name:        "Full-ByAgeGroup",
			mode:        sciensano.ByAgeGroup,
			doseType:    sciensano.Full,
			wantColumns: []string{"(unknown)"},
			wantValues:  []float64{0, 0, 0},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			s := &store.Store{Logger: l}
			r := ProRater{
				Name:     tt.name,
				PopStore: f,
				Mode:     tt.mode,
				DoseType: tt.doseType,
				Store:    s,
				Logger:   l,
			}
			r.createReport(vaccinations)

			report, err := s.Get(tt.name)
			require.NoError(t, err)

			assert.Equal(t, tt.wantColumns, report.GetColumns())
			for _, columnName := range tt.wantColumns {
				values, ok := report.GetValues(columnName)
				require.True(t, ok)
				assert.Equal(t, tt.wantValues, values)
			}
		})
	}
}
