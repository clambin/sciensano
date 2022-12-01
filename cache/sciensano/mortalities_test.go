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

	if *update {
		f, err := os.OpenFile(filepath.Join("input", "mortalities.json"), os.O_TRUNC|os.O_WRONLY, 0644)
		require.NoError(t, err)
		decoder := json.NewEncoder(f)
		decoder.SetIndent("", "  ")
		err = decoder.Encode(m)
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
	return makeResponse[sciensano.Mortality](count, func(timestamp time.Time, region, province, ageGroup, manufacturer string, _ sciensano.DoseType) *sciensano.Mortality {
		return &sciensano.Mortality{
			TimeStamp: sciensano.TimeStamp{Time: timestamp},
			Region:    region,
			AgeGroup:  ageGroup,
			Deaths:    1,
		}
	})
}
