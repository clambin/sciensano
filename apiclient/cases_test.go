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

func TestClient_GetCases(t *testing.T) {
	testServer := fake.Handler{}
	apiServer := httptest.NewServer(http.HandlerFunc(testServer.Handle))

	client := apiclient.Client{
		URL:        apiServer.URL,
		HTTPClient: &http.Client{},
	}

	ctx := context.Background()
	result, err := client.GetCases(ctx)

	require.NoError(t, err)
	require.Len(t, result, 2)

	assert.Equal(t, &apiclient.APICasesResponseEntry{
		TimeStamp: apiclient.TimeStamp{Time: time.Date(2020, 3, 1, 0, 0, 0, 0, time.UTC)},
		Province:  "VlaamsBrabant",
		Region:    "Flanders",
		AgeGroup:  "40-49",
		Cases:     1,
	}, result[0])

	assert.Equal(t, &apiclient.APICasesResponseEntry{
		TimeStamp: apiclient.TimeStamp{Time: time.Date(2020, 3, 1, 0, 0, 0, 0, time.UTC)},
		Province:  "Brussels",
		Region:    "Brussels",
		AgeGroup:  "40-49",
		Cases:     2,
	}, result[1])

	testServer.Fail = true
	_, err = client.GetCases(ctx)
	require.Error(t, err)

	apiServer.Close()
	_, err = client.GetCases(ctx)
	require.Error(t, err)
}

func TestAPICasesResponseEntry_GetTimestamp(t *testing.T) {
	timestamp := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	entry := apiclient.APICasesResponseEntry{
		TimeStamp: apiclient.TimeStamp{Time: timestamp},
	}

	assert.Equal(t, timestamp, entry.GetTimestamp())
}

func TestAPICasesResponseEntry_GetGroupFieldValue(t *testing.T) {
	entry := apiclient.APICasesResponseEntry{
		Province: "VlaamsBrabant",
		Region:   "Flanders",
		AgeGroup: "85+",
	}

	assert.Equal(t, "Flanders", entry.GetGroupFieldValue(apiclient.GroupByRegion))
	assert.Equal(t, "VlaamsBrabant", entry.GetGroupFieldValue(apiclient.GroupByProvince))
	assert.Equal(t, "85+", entry.GetGroupFieldValue(apiclient.GroupByAgeGroup))
}
