package sciensano_test

import (
	"encoding/json"
	"github.com/clambin/sciensano/apiclient"
	"github.com/clambin/sciensano/apiclient/sciensano"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestAPIMortalityResponses(t *testing.T) {
	gp := filepath.Join("testdata", t.Name()+".golden")
	if *update {
		body, err := json.Marshal(mortalityResponses)
		require.NoError(t, err)

		err = os.WriteFile(gp, body, 0644)
		require.NoError(t, err)
	}

	body, err := os.ReadFile(gp)
	require.NoError(t, err)

	var output []sciensano.APIMortalityResponse
	err = json.Unmarshal(body, &output)
	require.NoError(t, err)
	require.Len(t, output, len(mortalityResponses))
}

func TestAPIMortalityResponses_Attributes(t *testing.T) {
	assert.Equal(t, time.Date(2022, time.June, 19, 0, 0, 0, 0, time.UTC), mortalityResponses[0].GetTimestamp())
	assert.Equal(t, []string{"total"}, mortalityResponses[0].GetAttributeNames())
	assert.Equal(t, []float64{1}, mortalityResponses[0].GetAttributeValues())
	assert.Equal(t, "Flanders", mortalityResponses[0].GetGroupFieldValue(apiclient.GroupByRegion))
	assert.Equal(t, "21-30", mortalityResponses[0].GetGroupFieldValue(apiclient.GroupByAgeGroup))
	assert.Equal(t, 1.0, mortalityResponses[0].GetTotalValue())
}

var (
	mortalityResponses = []sciensano.APIMortalityResponse{
		{
			TimeStamp: sciensano.TimeStamp{Time: time.Date(2022, time.June, 19, 0, 0, 0, 0, time.UTC)},
			Region:    "Flanders",
			AgeGroup:  "21-30",
			Deaths:    1,
		},
	}
)
