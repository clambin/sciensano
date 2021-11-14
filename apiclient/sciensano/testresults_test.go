package sciensano_test

import (
	"context"
	"github.com/clambin/sciensano/apiclient/sciensano"
	"github.com/clambin/sciensano/apiclient/sciensano/fake"
	"github.com/clambin/sciensano/measurement"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestClient_GetTestResults(t *testing.T) {
	testServer := fake.Handler{}
	apiServer := httptest.NewServer(http.HandlerFunc(testServer.Handle))

	client := sciensano.Client{
		URL:        apiServer.URL,
		HTTPClient: &http.Client{},
	}

	ctx := context.Background()
	result, err := client.GetTestResults(ctx)

	require.NoError(t, err)
	require.Len(t, result, 3)
	assert.Equal(t, &sciensano.APITestResultsResponseEntry{
		TimeStamp: sciensano.TimeStamp{Time: time.Date(2021, 3, 11, 0, 0, 0, 0, time.UTC)},
		Region:    "Flanders",
		Province:  "",
		Total:     15,
		Positive:  10,
	}, result[2])

	testServer.Fail = true
	_, err = client.GetTestResults(ctx)
	require.Error(t, err)

	apiServer.Close()
	_, err = client.GetTestResults(ctx)
	require.Error(t, err)
}

func TestClient_TestResult_Measurement(t *testing.T) {
	entry := sciensano.APITestResultsResponseEntry{
		TimeStamp: sciensano.TimeStamp{Time: time.Now()},
		Region:    "Flanders",
		Province:  "VlaamsBrabant",
		Total:     100,
		Positive:  10,
	}

	assert.NotZero(t, entry.GetTimestamp())
	assert.Equal(t, "Flanders", entry.GetGroupFieldValue(measurement.GroupByRegion))
	assert.Equal(t, "VlaamsBrabant", entry.GetGroupFieldValue(measurement.GroupByProvince))
	assert.Empty(t, entry.GetGroupFieldValue(measurement.GroupByAgeGroup))
	assert.Equal(t, 100.0, entry.GetTotalValue())
	assert.Equal(t, []string{"total", "positive"}, entry.GetAttributeNames())
	assert.Equal(t, []float64{100, 10}, entry.GetAttributeValues())
}
