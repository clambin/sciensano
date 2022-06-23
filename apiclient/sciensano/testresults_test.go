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

func TestAPITestResultsResponse(t *testing.T) {
	gp := filepath.Join("testdata", t.Name()+".golden")
	if *update {
		body, err := json.Marshal(testResultsResponses)
		require.NoError(t, err)

		err = os.WriteFile(gp, body, 0644)
		require.NoError(t, err)
	}

	body, err := os.ReadFile(gp)
	require.NoError(t, err)

	var output []sciensano.APITestResultsResponse
	err = json.Unmarshal(body, &output)
	require.NoError(t, err)
	require.Len(t, output, len(mortalityResponses))
}

func TestAPITestResultsResponse_Attributes(t *testing.T) {
	assert.Equal(t, time.Date(2022, time.June, 19, 0, 0, 0, 0, time.UTC), testResultsResponses[0].GetTimestamp())
	assert.Equal(t, []string{"total", "positive"}, testResultsResponses[0].GetAttributeNames())
	assert.Equal(t, []float64{10, 1}, testResultsResponses[0].GetAttributeValues())
	assert.Equal(t, "Flanders", testResultsResponses[0].GetGroupFieldValue(apiclient.GroupByRegion))
	assert.Equal(t, "VlaamsBrabant", testResultsResponses[0].GetGroupFieldValue(apiclient.GroupByProvince))
	assert.Equal(t, 10.0, testResultsResponses[0].GetTotalValue())
}

var (
	testResultsResponses = []sciensano.APITestResultsResponse{
		{
			TimeStamp: sciensano.TimeStamp{Time: time.Date(2022, time.June, 19, 0, 0, 0, 0, time.UTC)},
			Region:    "Flanders",
			Province:  "VlaamsBrabant",
			Total:     10,
			Positive:  1,
		},
	}
)
