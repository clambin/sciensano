package sciensano_test

import (
	"encoding/json"
	"github.com/clambin/sciensano/cache/sciensano"
	"github.com/mailru/easyjson"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
	"time"
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
	f, err := os.Open(filepath.Join("input", "vaccinations.json"))
	require.NoError(t, err)

	var input sciensano.Vaccinations
	err = json.NewDecoder(f).Decode(&input)
	require.NoError(t, err)
	assert.NotZero(t, len(input))
}

func BenchmarkVaccinations_Unmarshal_Json(b *testing.B) {
	content, err := os.ReadFile(filepath.Join("input", "vaccinations.json"))
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
	content, err := os.ReadFile(filepath.Join("input", "vaccinations.json"))
	require.NoError(b, err)

	for i := 0; i < b.N; i++ {
		var vaccinations sciensano.Vaccinations
		err = easyjson.Unmarshal(content, &vaccinations)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func TestVaccinations_Summarize(t *testing.T) {
	vaccinations := makeVaccinations(1)

	if *update {
		f, err := os.OpenFile(filepath.Join("input", "vaccinations.json"), os.O_TRUNC|os.O_WRONLY, 0644)
		require.NoError(t, err)
		decoder := json.NewEncoder(f)
		decoder.SetIndent("", "  ")
		err = decoder.Encode(vaccinations)
		require.NoError(t, err)
		_ = f.Close()
	}

	testCases := []struct {
		summaryColumn   sciensano.SummaryColumn
		err             string
		expectedColumns []string
	}{
		{summaryColumn: sciensano.Total, expectedColumns: []string{"Total"}},
		{summaryColumn: sciensano.ByRegion, expectedColumns: regions},
		{summaryColumn: sciensano.ByManufacturer, expectedColumns: manufacturers},
		{summaryColumn: sciensano.ByAgeGroup, expectedColumns: ageGroups},
		{summaryColumn: sciensano.ByProvince, err: "vaccinations: invalid summary column: ByProvince"},
	}

	for _, tt := range testCases {
		t.Run(tt.summaryColumn.String(), func(t *testing.T) {
			table, err := vaccinations.Summarize(tt.summaryColumn)
			if tt.err != "" {
				assert.Equal(t, tt.err, err.Error())
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.expectedColumns, table.GetColumns())
		})
	}
}

func BenchmarkVaccinations_Summarize_Total(b *testing.B) {
	vaccinations := makeVaccinations(365)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := vaccinations.Summarize(sciensano.Total)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkVaccinations_Summarize_ByRegion(b *testing.B) {
	vaccinations := makeVaccinations(365)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := vaccinations.Summarize(sciensano.ByRegion)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func TestVaccinations_Categorize(t *testing.T) {
	vaccinations := makeVaccinations(1)
	c := vaccinations.Categorize()
	assert.Equal(t, []string{"booster", "booster2", "booster3", "full", "partial"}, c.GetColumns())
	records := len(c.GetTimestamps())
	for _, col := range c.GetColumns() {
		values, ok := c.GetValues(col)
		require.True(t, ok)
		assert.Len(t, values, records)
	}
}

func makeVaccinations(count int) sciensano.Vaccinations {
	return makeResponse[sciensano.Vaccination](count, func(timestamp time.Time, region, province, ageGroup, manufacturer string, dose sciensano.DoseType) *sciensano.Vaccination {
		return &sciensano.Vaccination{
			TimeStamp:    sciensano.TimeStamp{Time: timestamp},
			Manufacturer: manufacturer,
			Region:       region,
			AgeGroup:     ageGroup,
			Dose:         dose,
			Count:        1,
		}
	})
}
