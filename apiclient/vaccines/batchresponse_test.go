package vaccines_test

import (
	"encoding/json"
	"flag"
	"github.com/clambin/sciensano/apiclient"
	"github.com/clambin/sciensano/apiclient/vaccines"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
	"time"
)

var update = flag.Bool("update", false, "update .golden files")

func TestAPIBatchResponses(t *testing.T) {
	gp := filepath.Join("testdata", t.Name()+".golden")
	if *update {
		body, err := json.Marshal(batchResponse)
		require.NoError(t, err)

		err = os.WriteFile(gp, body, 0644)
		require.NoError(t, err)
	}

	body, err := os.ReadFile(gp)
	require.NoError(t, err)

	var output []vaccines.APIBatchResponse
	err = json.Unmarshal(body, &output)
	require.NoError(t, err)
	require.Len(t, output, len(batchResponse))
}

func TestBatch_Measurement(t *testing.T) {
	b := batchResponse[0]

	assert.NotZero(t, b.GetTimestamp())
	assert.Empty(t, b.GetGroupFieldValue(apiclient.GroupByAgeGroup))
	assert.Equal(t, "A", b.GetGroupFieldValue(apiclient.GroupByManufacturer))
	assert.Equal(t, 200.0, b.GetTotalValue())
	assert.Equal(t, []string{"total"}, b.GetAttributeNames())
	assert.Equal(t, []float64{200}, b.GetAttributeValues())
}

var (
	batchResponse = []*vaccines.APIBatchResponse{
		{Date: vaccines.Timestamp{Time: time.Now()}, Manufacturer: "A", Amount: 200},
		{Date: vaccines.Timestamp{Time: time.Now().Add(-24 * time.Hour)}, Manufacturer: "A", Amount: 200},
	}
)
