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

func TestMortalities_Unmarshal(t *testing.T) {
	f, err := os.Open(filepath.Join("input", "mortalities.json"))
	require.NoError(t, err)

	var input sciensano.Mortalities
	err = json.NewDecoder(f).Decode(&input)
	require.NoError(t, err)
	assert.NotZero(t, len(input))
}

func TestMortalities_Summarize(t *testing.T) {
	m := makeMortalities(1)

	testCases := []struct {
		summaryColumn   sciensano.SummaryColumn
		err             string
		expectedColumns []string
	}{
		{summaryColumn: sciensano.Total, expectedColumns: []string{"Total"}},
		{summaryColumn: sciensano.ByRegion, expectedColumns: regions},
		{summaryColumn: sciensano.ByAgeGroup, expectedColumns: ageGroups},
		{summaryColumn: sciensano.ByProvince, err: "mortalities: invalid summary column: ByProvince"},
	}

	for _, tt := range testCases {
		t.Run(tt.summaryColumn.String(), func(t *testing.T) {
			table, err := m.Summarize(tt.summaryColumn)
			if tt.err != "" {
				assert.Equal(t, tt.err, err.Error())
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.expectedColumns, table.GetColumns())
		})
	}
}

func BenchmarkMortalities_Summarize_Total(b *testing.B) {
	m := makeMortalities(365)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := m.Summarize(sciensano.Total)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMortalities_Summarize_ByAgeGroup(b *testing.B) {
	m := makeMortalities(365)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := m.Summarize(sciensano.ByAgeGroup)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func makeMortalities(count int) sciensano.Mortalities {
	mortalities := make(sciensano.Mortalities, 0)

	timestamp := sciensano.TimeStamp{Time: time.Date(2022, time.November, 19, 0, 0, 0, 0, time.UTC)}
	for i := 0; i < count; i++ {
		for _, region := range regions {
			for _, ageGroup := range ageGroups {
				mortalities = append(mortalities, sciensano.Mortality{
					TimeStamp: timestamp,
					Region:    region,
					AgeGroup:  ageGroup,
					Deaths:    1,
				})
			}
		}
		timestamp.Time = timestamp.Time.Add(24 * time.Hour)
	}
	return mortalities
}
