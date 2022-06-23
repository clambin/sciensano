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

func TestAPIHospitalisationsResponses(t *testing.T) {
	gp := filepath.Join("testdata", t.Name()+".golden")
	if *update {
		body, err := json.Marshal(hospitalisationResponses)
		require.NoError(t, err)

		err = os.WriteFile(gp, body, 0644)
		require.NoError(t, err)
	}

	body, err := os.ReadFile(gp)
	require.NoError(t, err)

	var output []sciensano.APIHospitalisationsResponse
	err = json.Unmarshal(body, &output)
	require.NoError(t, err)
	require.Len(t, output, len(hospitalisationResponses))
}

func TestAPIHospitalisationsResponse_Attributes(t *testing.T) {
	assert.Equal(t, time.Date(2022, time.June, 19, 0, 0, 0, 0, time.UTC), hospitalisationResponses[0].GetTimestamp())
	assert.Equal(t, []string{"in", "inICU", "inResp", "inECMO"}, hospitalisationResponses[0].GetAttributeNames())
	assert.Equal(t, []float64{100, 20, 10, 30}, hospitalisationResponses[0].GetAttributeValues())
	assert.Equal(t, "Flanders", hospitalisationResponses[0].GetGroupFieldValue(apiclient.GroupByRegion))
	assert.Equal(t, "VlaamsBrabant", hospitalisationResponses[0].GetGroupFieldValue(apiclient.GroupByProvince))
	assert.Equal(t, 100.0, hospitalisationResponses[0].GetTotalValue())
}

var (
	hospitalisationResponses = []sciensano.APIHospitalisationsResponse{
		{
			TimeStamp:   sciensano.TimeStamp{Time: time.Date(2022, time.June, 19, 0, 0, 0, 0, time.UTC)},
			Province:    "VlaamsBrabant",
			Region:      "Flanders",
			TotalIn:     100,
			TotalInECMO: 30,
			TotalInICU:  20,
			TotalInResp: 10,
		},
	}
)
