package apiclient_test

import (
	"context"
	"github.com/clambin/sciensano/apiclient"
	"github.com/clambin/sciensano/apiclient/fake"
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

	client := apiclient.Client{
		URL:        apiServer.URL,
		HTTPClient: &http.Client{},
	}

	ctx := context.Background()
	result, err := client.GetTestResults(ctx)

	require.NoError(t, err)
	require.Len(t, result, 3)
	assert.Equal(t, &apiclient.APITestResultsResponseEntry{
		TimeStamp: apiclient.TimeStamp{Time: time.Date(2021, 3, 11, 0, 0, 0, 0, time.UTC)},
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

func TestAPITestResultsResponseEntry_GetTimestamp(t *testing.T) {
	timestamp := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	entry := apiclient.APITestResultsResponseEntry{
		TimeStamp: apiclient.TimeStamp{Time: timestamp},
	}

	assert.Equal(t, timestamp, entry.GetTimestamp())
}

func TestAPITestResultsResponseEntry_GetGroupFieldValue(t *testing.T) {
	entry := apiclient.APITestResultsResponseEntry{
		Region:   "Flanders",
		Province: "VlaamsBrabant",
	}

	assert.Equal(t, "Flanders", entry.GetGroupFieldValue(apiclient.GroupByRegion))
	assert.Equal(t, "VlaamsBrabant", entry.GetGroupFieldValue(apiclient.GroupByProvince))
}
