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

func TestCases_Unmarshal(t *testing.T) {
	f, err := os.Open(filepath.Join("input", "cases.json"))
	require.NoError(t, err)

	var input sciensano.Cases
	err = json.NewDecoder(f).Decode(&input)
	require.NoError(t, err)
	assert.NotZero(t, len(input))
}

func TestCases_Summarize(t *testing.T) {
	const dayCount = 1
	cases := sciensano.Cases(makeResponse[sciensano.Case](1, func(timestamp time.Time, region, province, ageGroup, _ string, _ sciensano.DoseType) *sciensano.Case {
		return &sciensano.Case{
			TimeStamp: sciensano.TimeStamp{Time: timestamp},
			Province:  province,
			Region:    region,
			AgeGroup:  ageGroup,
			Cases:     1,
		}
	}))

	if *update {
		f, err := os.OpenFile(filepath.Join("input", "cases.json"), os.O_TRUNC|os.O_WRONLY, 0644)
		require.NoError(t, err)
		decoder := json.NewEncoder(f)
		decoder.SetIndent("", "  ")
		err = decoder.Encode(cases)
		require.NoError(t, err)
		_ = f.Close()
	}

	testCases := []struct {
		summaryColumn   sciensano.SummaryColumn
		pass            bool
		expectedColumns int
	}{
		{
			summaryColumn:   sciensano.Total,
			pass:            true,
			expectedColumns: 1,
		},
		{
			summaryColumn:   sciensano.ByAgeGroup,
			pass:            true,
			expectedColumns: len(ageGroups),
		},
		{
			summaryColumn:   sciensano.ByRegion,
			pass:            true,
			expectedColumns: len(regions),
		},
		{
			summaryColumn:   sciensano.ByProvince,
			pass:            true,
			expectedColumns: len(provinces),
		},
		{
			summaryColumn: sciensano.ByManufacturer,
			pass:          false,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.summaryColumn.String(), func(t *testing.T) {
			d, err := cases.Summarize(tt.summaryColumn)
			if !tt.pass {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Len(t, d.GetTimestamps(), dayCount)
			assert.Len(t, d.GetColumns(), tt.expectedColumns)
		})
	}
}
