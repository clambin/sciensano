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

func TestMortalities_Unmarshal(t *testing.T) {
	f, err := os.Open(filepath.Join("testutil", "testdata", "mortalities.json"))
	require.NoError(t, err)

	var input sciensano.Mortalities
	err = json.NewDecoder(f).Decode(&input)
	require.NoError(t, err)
	assert.NotZero(t, len(input))
}

func BenchmarkMortalities_Unmarshal_JSON(b *testing.B) {
	content, err := os.ReadFile(filepath.Join("testutil", "testdata", "mortalities.json"))
	require.NoError(b, err)

	for range b.N {
		var records sciensano.Mortalities
		err = json.Unmarshal(content, &records)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func TestMortalities_Summarize(t *testing.T) {
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
			want:          []string{"Brussels", "Flanders", "Wallonia"},
		},
		{
			summaryColumn: sciensano.ByAgeGroup,
			wantErr:       assert.NoError,
			want:          []string{"(unknown)", "0-24", "25-44", "45-64", "65-74", "75-84", "85+"},
		},
		{
			summaryColumn: sciensano.ByProvince,
			wantErr:       assert.Error,
		},
	}

	mortalities := testutil.Mortalities()
	for _, tt := range testCases {
		t.Run(tt.summaryColumn.String(), func(t *testing.T) {
			table, err := mortalities.Summarize(tt.summaryColumn)
			tt.wantErr(t, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.want, table.GetColumns())
		})
	}
}

func BenchmarkMortalities_Summarize_Total(b *testing.B) {
	mortalities := testutil.Mortalities()
	b.ResetTimer()
	for range b.N {
		_, err := mortalities.Summarize(sciensano.Total)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMortalities_Summarize_ByAgeGroup(b *testing.B) {
	mortalities := testutil.Mortalities()
	b.ResetTimer()
	for range b.N {
		_, err := mortalities.Summarize(sciensano.ByAgeGroup)
		if err != nil {
			b.Fatal(err)
		}
	}
}
