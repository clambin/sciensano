package sciensano_test

import (
	"github.com/clambin/sciensano/apiclient"
	"github.com/clambin/sciensano/apiclient/sciensano"
	"github.com/mailru/easyjson"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestAPICasesResponses(t *testing.T) {
	gp := filepath.Join("testdata", t.Name()+".golden")
	if *update {
		body, err := easyjson.Marshal(casesResponses)
		require.NoError(t, err)

		err = os.WriteFile(gp, body, 0644)
		require.NoError(t, err)
	}

	body, err := os.ReadFile(gp)
	require.NoError(t, err)

	var output sciensano.APICasesResponses
	err = easyjson.Unmarshal(body, &output)
	require.NoError(t, err)
	require.Len(t, output, len(casesResponses))
}

func TestAPICasesResponse_Attributes(t *testing.T) {
	groups := []string{"55-54", "45-54"}
	for idx, entry := range casesResponses {
		assert.Equal(t, time.Date(2022, time.June, 19, 0, 0, 0, 0, time.UTC), entry.GetTimestamp(), idx)
		assert.Equal(t, []string{"total"}, entry.GetAttributeNames())
		assert.Equal(t, []float64{float64(idx + 1)}, entry.GetAttributeValues())
		assert.Equal(t, float64(idx+1), entry.GetTotalValue())
		assert.Equal(t, "Flanders", entry.GetGroupFieldValue(apiclient.GroupByRegion))
		assert.Equal(t, "VlaamsBrabant", entry.GetGroupFieldValue(apiclient.GroupByProvince))
		assert.Equal(t, groups[idx], entry.GetGroupFieldValue(apiclient.GroupByAgeGroup))
	}
}

var (
	casesResponses = sciensano.APICasesResponses{
		&sciensano.APICasesResponse{
			TimeStamp: sciensano.TimeStamp{Time: time.Date(2022, time.June, 19, 0, 0, 0, 0, time.UTC)},
			Province:  "VlaamsBrabant",
			Region:    "Flanders",
			AgeGroup:  "55-54",
			Cases:     1,
		},
		&sciensano.APICasesResponse{
			TimeStamp: sciensano.TimeStamp{Time: time.Date(2022, time.June, 19, 0, 0, 0, 0, time.UTC)},
			Province:  "VlaamsBrabant",
			Region:    "Flanders",
			AgeGroup:  "45-54",
			Cases:     2,
		},
	}
)
