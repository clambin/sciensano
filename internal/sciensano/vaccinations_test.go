package sciensano_test

import (
	"encoding/json"
	"github.com/clambin/sciensano/internal/sciensano"
	"github.com/clambin/sciensano/internal/sciensano/testutil"
	"github.com/mailru/easyjson"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func TestDoseType_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		input   sciensano.DoseType
		encoded string
	}{
		{input: sciensano.Partial, encoded: `"A"`},
		{input: sciensano.Full, encoded: `"B"`},
		{input: sciensano.SingleDose, encoded: `"C"`},
		{input: sciensano.Booster, encoded: `"E"`},
		{input: sciensano.Booster2, encoded: `"E2"`},
		{input: sciensano.Booster3, encoded: `"E3"`},
	}

	for _, tt := range tests {
		t.Run(tt.input.String(), func(t *testing.T) {
			body, err := json.Marshal(tt.input)
			require.NoError(t, err)
			assert.Equal(t, tt.encoded, string(body))

			var d sciensano.DoseType
			err = json.Unmarshal(body, &d)
			require.NoError(t, err)

			assert.Equal(t, tt.input, d)
		})
	}
}

func TestVaccinations_Unmarshal(t *testing.T) {
	f, err := os.Open(filepath.Join("testutil", "testdata", "vaccinations.json"))
	require.NoError(t, err)

	var input sciensano.Vaccinations
	err = json.NewDecoder(f).Decode(&input)
	require.NoError(t, err)
	assert.NotZero(t, len(input))
}

func BenchmarkVaccinations_Unmarshal_Json(b *testing.B) {
	content, err := os.ReadFile(filepath.Join("testutil", "testdata", "vaccinations.json"))
	require.NoError(b, err)

	for i := 0; i < b.N; i++ {
		var vaccinations sciensano.Vaccinations
		err = json.Unmarshal(content, &vaccinations)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkVaccinations_Unmarshal_Easyjson(b *testing.B) {
	content, err := os.ReadFile(filepath.Join("testutil", "testdata", "vaccinations.json"))
	require.NoError(b, err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var vaccinations sciensano.Vaccinations
		err = easyjson.Unmarshal(content, &vaccinations)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func TestVaccinations_Summarize(t *testing.T) {
	testCases := []struct {
		summaryColumn sciensano.SummaryColumn
		wantErr       assert.ErrorAssertionFunc
		want          []string
	}{
		{
			summaryColumn: sciensano.Total,
			wantErr:       assert.NoError,
			want:          []string{"Total"},
		},
		{
			summaryColumn: sciensano.ByRegion,
			wantErr:       assert.NoError,
			want:          []string{"(unknown)", "Brussels", "Flanders", "Ostbelgien", "Wallonia"},
		},
		{
			summaryColumn: sciensano.ByManufacturer,
			wantErr:       assert.NoError,
			want:          []string{"AstraZeneca-Oxford", "Johnson&Johnson", "Moderna Bivalent BA1", "Moderna Original", "Novavax", "Other", "Pfizer-BioNTech Bivalent BA1", "Pfizer-BioNTech Bivalent BA4-5", "Pfizer-BioNTech Original"},
		},
		{
			summaryColumn: sciensano.ByAgeGroup,
			wantErr:       assert.NoError,
			want:          []string{"00-04", "05-11", "12-15", "16-17", "18-24", "25-34", "35-44", "45-54", "55-64", "65-74", "75-84", "85+"},
		},
		{
			summaryColumn: sciensano.ByProvince,
			wantErr:       assert.Error,
		},
		{
			summaryColumn: sciensano.ByVaccinationType,
			wantErr:       assert.NoError,
			want:          []string{"booster", "booster2", "booster3", "full", "partial"},
		},
	}

	vaccinations := testutil.Vaccinations()
	for _, tt := range testCases {
		t.Run(tt.summaryColumn.String(), func(t *testing.T) {
			table, err := vaccinations.Summarize(tt.summaryColumn)
			tt.wantErr(t, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.want, table.GetColumns())
		})
	}
}

func BenchmarkVaccinations_Summarize_Total(b *testing.B) {
	vaccinations := testutil.Vaccinations()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := vaccinations.Summarize(sciensano.Total)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkVaccinations_Summarize_ByRegion(b *testing.B) {
	vaccinations := testutil.Vaccinations()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := vaccinations.Summarize(sciensano.ByRegion)
		if err != nil {
			b.Fatal(err)
		}
	}
}
