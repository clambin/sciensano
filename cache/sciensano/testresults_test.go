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

func TestTestResults_Unmarshal(t *testing.T) {
	f, err := os.Open(filepath.Join("input", "testResults.json"))
	require.NoError(t, err)

	var input sciensano.TestResults
	err = json.NewDecoder(f).Decode(&input)
	require.NoError(t, err)
	assert.NotZero(t, len(input))
}

func TestTestResults_Summarize(t *testing.T) {
	testResults := makeTestResults(1)

	if *update {
		f, err := os.OpenFile(filepath.Join("input", "testResults.json"), os.O_TRUNC|os.O_WRONLY, 0644)
		require.NoError(t, err)
		decoder := json.NewEncoder(f)
		decoder.SetIndent("", "  ")
		err = decoder.Encode(testResults)
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
		{summaryColumn: sciensano.ByManufacturer, err: "testResults: invalid summary column: ByManufacturer"},
		{summaryColumn: sciensano.ByAgeGroup, err: "testResults: invalid summary column: ByAgeGroup"},
		{summaryColumn: sciensano.ByProvince, expectedColumns: []string{"bar", "foo", "snafu"}},
	}

	for _, tt := range testCases {
		t.Run(tt.summaryColumn.String(), func(t *testing.T) {
			table, err := testResults.Summarize(tt.summaryColumn)
			if tt.err != "" {
				assert.Equal(t, tt.err, err.Error())
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.expectedColumns, table.GetColumns())
		})
	}
}

func TestTestResults_Categorize(t *testing.T) {
	testResults := makeTestResults(1)
	c := testResults.Categorize()
	assert.Equal(t, []string{"positive", "rate", "total"}, c.GetColumns())
	records := len(c.GetTimestamps())
	for _, col := range c.GetColumns() {
		values, ok := c.GetValues(col)
		require.True(t, ok)
		assert.Len(t, values, records)
	}
}

func makeTestResults(count int) sciensano.TestResults {
	return makeResponse[sciensano.TestResult](count, func(timestamp time.Time, region, province, ageGroup, _ string, _ sciensano.DoseType) *sciensano.TestResult {
		return &sciensano.TestResult{
			TimeStamp: sciensano.TimeStamp{Time: timestamp},
			Province:  province,
			Region:    region,
			Total:     2,
			Positive:  1,
		}
	})
}
