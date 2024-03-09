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

func TestHospitalisations_Unmarshal(t *testing.T) {
	f, err := os.Open(filepath.Join("testutil", "testdata", "hospitalisations.json"))
	require.NoError(t, err)

	var input sciensano.Hospitalisations
	err = json.NewDecoder(f).Decode(&input)
	require.NoError(t, err)
	assert.NotZero(t, len(input))
}

func BenchmarkHospitalisations_Unmarshal_JSON(b *testing.B) {
	content, err := os.ReadFile(filepath.Join("testutil", "testdata", "hospitalisations.json"))
	require.NoError(b, err)

	for range b.N {
		var records sciensano.Hospitalisations
		err = json.Unmarshal(content, &records)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func TestHospitalisations_Summarize(t *testing.T) {
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
			summaryColumn: sciensano.ByAgeGroup,
			wantErr:       assert.Error,
		},
		{
			summaryColumn: sciensano.ByRegion,
			wantErr:       assert.NoError,
			want:          []string{"Brussels", "Flanders", "Wallonia"},
		},
		{
			summaryColumn: sciensano.ByProvince,
			wantErr:       assert.NoError,
			want:          []string{"Antwerpen", "BrabantWallon", "Brussels", "Hainaut", "Limburg", "Liège", "Luxembourg", "Namur", "OostVlaanderen", "VlaamsBrabant", "WestVlaanderen"},
		},
		{
			summaryColumn: sciensano.ByCategory,
			wantErr:       assert.NoError,
			want:          []string{"in", "inECMO", "inICU", "inResp"},
		},
		{
			summaryColumn: sciensano.ByManufacturer,
			wantErr:       assert.Error,
		},
	}

	hospitalisations := testutil.Hospitalisations()
	for _, tt := range testCases {
		t.Run(tt.summaryColumn.String(), func(t *testing.T) {
			d, err := hospitalisations.Summarize(tt.summaryColumn)
			tt.wantErr(t, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.want, d.GetColumns())
		})
	}
}
