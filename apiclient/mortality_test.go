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

func TestClient_GetMortality(t *testing.T) {
	testServer := fake.Handler{}
	apiServer := httptest.NewServer(http.HandlerFunc(testServer.Handle))

	client := &apiclient.Client{
		URL:        apiServer.URL,
		HTTPClient: &http.Client{},
	}

	ctx := context.Background()
	result, err := client.GetMortality(ctx)

	require.NoError(t, err)

	assert.Equal(t, &apiclient.APIMortalityResponseEntry{
		TimeStamp: apiclient.TimeStamp{Time: time.Date(2020, time.March, 10, 0, 0, 0, 0, time.UTC)}, Region: "Brussels", AgeGroup: "85+", Deaths: 1,
	}, result[0])

	assert.Equal(t, &apiclient.APIMortalityResponseEntry{
		TimeStamp: apiclient.TimeStamp{Time: time.Date(2020, time.March, 10, 0, 0, 0, 0, time.UTC)}, Region: "Brussels", AgeGroup: "85+", Deaths: 2,
	}, result[1])

	testServer.Fail = true
	_, err = client.GetMortality(ctx)
	require.Error(t, err)

	apiServer.Close()
	_, err = client.GetMortality(ctx)
	require.Error(t, err)
}

func TestAPINMortalityResponseEntry_GetTimestamp(t *testing.T) {
	timestamp := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	entry := apiclient.APIMortalityResponseEntry{
		TimeStamp: apiclient.TimeStamp{Time: timestamp},
	}

	assert.Equal(t, timestamp, entry.GetTimestamp())
}

func TestAPIMortalilityResponseEntry_GetGroupFieldValue(t *testing.T) {
	entry := apiclient.APIMortalityResponseEntry{
		Region:   "Flanders",
		AgeGroup: "85+",
	}

	assert.Equal(t, "Flanders", entry.GetGroupFieldValue(apiclient.GroupByRegion))
	assert.Equal(t, "85+", entry.GetGroupFieldValue(apiclient.GroupByAgeGroup))
}
