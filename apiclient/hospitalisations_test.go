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

func TestClient_GetHospitalisations(t *testing.T) {
	testServer := fake.Handler{}
	apiServer := httptest.NewServer(http.HandlerFunc(testServer.Handle))

	client := apiclient.Client{
		URL:        apiServer.URL,
		HTTPClient: &http.Client{},
	}

	ctx := context.Background()
	result, err := client.GetHospitalisations(ctx)

	require.NoError(t, err)
	require.Len(t, result, 2)

	assert.Equal(t, &apiclient.APIHospitalisationsResponseEntry{
		TimeStamp:   apiclient.TimeStamp{Time: time.Date(2020, time.March, 15, 0, 0, 0, 0, time.UTC)},
		Province:    "Brussels",
		Region:      "Brussels",
		TotalIn:     58,
		TotalInICU:  11,
		TotalInResp: 8,
		TotalInECMO: 0,
	}, result[0])
	assert.Equal(t, &apiclient.APIHospitalisationsResponseEntry{
		TimeStamp:   apiclient.TimeStamp{Time: time.Date(2020, time.March, 15, 0, 0, 0, 0, time.UTC)},
		Province:    "VlaamsBrabant",
		Region:      "Flanders",
		TotalIn:     13,
		TotalInICU:  2,
		TotalInResp: 0,
		TotalInECMO: 1,
	}, result[1])

	testServer.Fail = true
	_, err = client.GetHospitalisations(ctx)
	require.Error(t, err)

	apiServer.Close()
	_, err = client.GetHospitalisations(ctx)
	require.Error(t, err)
}

func TestAPIHospitalisationsResponseEntry_GetTimestamp(t *testing.T) {
	timestamp := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	entry := apiclient.APIHospitalisationsResponseEntry{
		TimeStamp: apiclient.TimeStamp{Time: timestamp},
	}

	assert.Equal(t, timestamp, entry.GetTimestamp())
}

func TestAPIHospitalisationsResponseEntry_GetGroupFieldValue(t *testing.T) {
	entry := apiclient.APIHospitalisationsResponseEntry{
		Province: "VlaamsBrabant",
		Region:   "Flanders",
	}

	assert.Equal(t, "Flanders", entry.GetGroupFieldValue(apiclient.GroupByRegion))
	assert.Equal(t, "VlaamsBrabant", entry.GetGroupFieldValue(apiclient.GroupByProvince))
}
