package apiclient_test

import (
	"context"
	"github.com/clambin/sciensano/apiclient"
	"github.com/clambin/sciensano/apiclient/fake"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestClient_GetVaccinations(t *testing.T) {
	testServer := fake.Handler{}
	apiServer := httptest.NewServer(http.HandlerFunc(testServer.Handle))
	defer apiServer.Close()

	client := apiclient.Client{
		URL:        apiServer.URL,
		HTTPClient: &http.Client{},
	}

	ctx := context.Background()
	result, err := client.GetVaccinations(ctx)

	require.NoError(t, err)
	require.Len(t, result, 7)
	assert.Equal(t, &apiclient.APIVaccinationsResponseEntry{
		TimeStamp: apiclient.TimeStamp{Time: time.Date(2021, time.March, 11, 0, 0, 0, 0, time.UTC)},
		Region:    "Flanders",
		AgeGroup:  "45-54",
		Dose:      "B",
		Count:     50,
	}, result[6])

	testServer.Fail = true
	_, err = client.GetVaccinations(ctx)
	require.Error(t, err)

	apiServer.Close()
	_, err = client.GetVaccinations(ctx)
	require.Error(t, err)

}

func TestAPIVaccinationsResponseEntry_GetTimestamp(t *testing.T) {
	timestamp := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	entry := apiclient.APIVaccinationsResponseEntry{
		TimeStamp: apiclient.TimeStamp{Time: timestamp},
	}

	assert.Equal(t, timestamp, entry.GetTimestamp())
}

func TestAPIVaccinationsResponseEntry_GetGroupFieldValue(t *testing.T) {
	entry := apiclient.APIVaccinationsResponseEntry{
		Region:   "Flanders",
		AgeGroup: "85+",
	}

	assert.Equal(t, "Flanders", entry.GetGroupFieldValue(apiclient.GroupByRegion))
	assert.Equal(t, "85+", entry.GetGroupFieldValue(apiclient.GroupByAgeGroup))
}

func BenchmarkClient_GetVaccinations(b *testing.B) {
	testServer := httptest.NewServer(http.HandlerFunc(handleVaccinationResponse))
	defer testServer.Close()

	client := apiclient.Client{
		HTTPClient: &http.Client{},
		URL:        testServer.URL,
	}
	_, err := client.GetVaccinations(context.Background())
	require.NoError(b, err)
}

var bigFile []byte

func handleVaccinationResponse(w http.ResponseWriter, _ *http.Request) {
	var err error
	if bigFile == nil {
		bigFile, err = os.ReadFile("../data/vaccinations.json")
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(bigFile)
}
