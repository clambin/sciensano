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

func TestHospitalisations_Unmarshal(t *testing.T) {
	f, err := os.Open(filepath.Join("input", "hospitalisations.json"))
	require.NoError(t, err)

	var input sciensano.Hospitalisations
	err = json.NewDecoder(f).Decode(&input)
	require.NoError(t, err)
	assert.NotZero(t, len(input))
}

func TestHospitalisations_Summarize(t *testing.T) {
	const dayCount = 1
	hospitalisations := makeHospitalisations(dayCount)

	if *update {
		f, err := os.OpenFile(filepath.Join("input", "hospitalisations.json"), os.O_TRUNC|os.O_WRONLY, 0644)
		require.NoError(t, err)
		decoder := json.NewEncoder(f)
		decoder.SetIndent("", "  ")
		err = decoder.Encode(hospitalisations)
		require.NoError(t, err)
		_ = f.Close()
		t.Log("updated")
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
			summaryColumn: sciensano.ByAgeGroup,
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
			d, err := hospitalisations.Summarize(tt.summaryColumn)
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

func TestHospitalisations_Categorize(t *testing.T) {
	hospitalisations := makeHospitalisations(1)
	c := hospitalisations.Categorize()
	assert.Equal(t, []string{"in", "inECMO", "inICU", "inResp"}, c.GetColumns())
	records := len(c.GetTimestamps())
	for _, col := range c.GetColumns() {
		values, ok := c.GetValues(col)
		require.True(t, ok)
		assert.Len(t, values, records)
	}
}

func makeHospitalisations(count int) sciensano.Hospitalisations {
	return makeResponse[sciensano.Hospitalisation](count, func(timestamp time.Time, region, province, ageGroup, _ string, _ sciensano.DoseType) *sciensano.Hospitalisation {
		return &sciensano.Hospitalisation{
			TimeStamp:   sciensano.TimeStamp{Time: timestamp},
			Province:    province,
			Region:      region,
			TotalIn:     10,
			TotalInICU:  3,
			TotalInResp: 2,
			TotalInECMO: 1,
		}
	})
}
