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

func TestCases_Unmarshal(t *testing.T) {
	f, err := os.Open(filepath.Join("testutil", "testdata", "cases.json"))
	require.NoError(t, err)

	var input sciensano.Cases
	err = json.NewDecoder(f).Decode(&input)
	require.NoError(t, err)
	assert.NotZero(t, len(input))
}

func BenchmarkCases_Unmarshal_JSON(b *testing.B) {
	content, err := os.ReadFile(filepath.Join("testutil", "testdata", "cases.json"))
	require.NoError(b, err)

	for i := 0; i < b.N; i++ {
		var records sciensano.Cases
		err = json.Unmarshal(content, &records)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCases_Unmarshal_EasyJSON(b *testing.B) {
	content, err := os.ReadFile(filepath.Join("testutil", "testdata", "cases.json"))
	require.NoError(b, err)

	for i := 0; i < b.N; i++ {
		var records sciensano.Cases
		err = easyjson.Unmarshal(content, &records)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func TestCases_Summarize(t *testing.T) {
	cases := testutil.Cases()

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
			wantErr:       assert.NoError,
			want:          wantAgeGroups,
		},
		{
			summaryColumn: sciensano.ByRegion,
			wantErr:       assert.NoError,
			want:          wantRegions,
		},
		{
			summaryColumn: sciensano.ByProvince,
			wantErr:       assert.NoError,
			want:          wantProvinces,
		},
		{
			summaryColumn: sciensano.ByManufacturer,
			wantErr:       assert.Error,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.summaryColumn.String(), func(t *testing.T) {
			d, err := cases.Summarize(tt.summaryColumn)
			tt.wantErr(t, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.want, d.GetColumns())
		})
	}
}
