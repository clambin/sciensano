package sciensano_test

import (
	"encoding/json"
	"github.com/clambin/sciensano/cache/sciensano"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestVaccinations_Unmarshal(t *testing.T) {
	f, err := os.Open(filepath.Join("input", "vaccinations.json"))
	require.NoError(t, err)

	var input sciensano.Vaccinations
	err = json.NewDecoder(f).Decode(&input)
	require.NoError(t, err)
	assert.NotZero(t, len(input))
}

func TestVaccinations_Summarize(t *testing.T) {
	vaccinations := makeVaccinations(1)
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

func makeVaccinations(count int) sciensano.Vaccinations {
	vaccinations := make(sciensano.Vaccinations, 0)
	timestamp := sciensano.TimeStamp{Time: time.Date(2022, time.November, 19, 0, 0, 0, 0, time.UTC)}
	for i := 0; i < count; i++ {
		for _, region := range regions {
			for _, manufacturer := range manufacturers {
				for _, ageGroup := range ageGroups {
					for _, dose := range doses {
						vaccinations = append(vaccinations, sciensano.Vaccination{
							TimeStamp:    timestamp,
							Manufacturer: manufacturer,
							Region:       region,
							AgeGroup:     ageGroup,
							Dose:         dose,
							Count:        1,
						})
					}
				}
			}
		}
		timestamp.Time = timestamp.Time.Add(24 * time.Hour)
	}
	return vaccinations
}
