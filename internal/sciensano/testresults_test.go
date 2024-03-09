package sciensano_test

import (
	"encoding/json"
	"github.com/clambin/sciensano/v2/internal/sciensano"
	"github.com/clambin/sciensano/v2/internal/sciensano/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func TestTestResults_Unmarshal(t *testing.T) {
	f, err := os.Open(filepath.Join("testutil", "testdata", "testResults.json"))
	require.NoError(t, err)

	var input sciensano.TestResults
	err = json.NewDecoder(f).Decode(&input)
	require.NoError(t, err)
	assert.NotZero(t, len(input))
}

func BenchmarkTestResults_Unmarshal_JSON(b *testing.B) {
	content, err := os.ReadFile(filepath.Join("testutil", "testdata", "testResults.json"))
	require.NoError(b, err)

	for range b.N {
		var records sciensano.TestResults
		err = json.Unmarshal(content, &records)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func TestTestResults_Summarize(t *testing.T) {
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
			want:          wantRegions,
		},
		{
			summaryColumn: sciensano.ByManufacturer,
			wantErr:       assert.Error,
		},
		{
			summaryColumn: sciensano.ByAgeGroup,
			wantErr:       assert.Error,
		},
		{
			summaryColumn: sciensano.ByProvince,
			wantErr:       assert.NoError,
			want:          wantProvinces,
		},
		{
			summaryColumn: sciensano.ByCategory,
			wantErr:       assert.NoError,
			want:          []string{"positive", "total"},
		},
	}

	testResults := testutil.TestResults()
	for _, tt := range testCases {
		t.Run(tt.summaryColumn.String(), func(t *testing.T) {
			table, err := testResults.Summarize(tt.summaryColumn)
			tt.wantErr(t, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.want, table.GetColumns())
		})
	}
}
