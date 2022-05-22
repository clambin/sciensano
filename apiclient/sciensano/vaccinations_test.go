package sciensano_test

import (
	"context"
	"github.com/clambin/go-metrics/client"
	"github.com/clambin/sciensano/apiclient"
	"github.com/clambin/sciensano/apiclient/sciensano"
	"github.com/clambin/sciensano/apiclient/sciensano/fake"
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

	c := sciensano.Client{
		URL:    apiServer.URL,
		Caller: &client.BaseClient{HTTPClient: http.DefaultClient},
	}

	ctx := context.Background()
	result, err := c.GetVaccinations(ctx)

	require.NoError(t, err)
	require.Len(t, result, 7)
	assert.Equal(t, &sciensano.APIVaccinationsResponse{
		TimeStamp: sciensano.TimeStamp{Time: time.Date(2021, time.March, 11, 0, 0, 0, 0, time.UTC)},
		Region:    "Flanders",
		AgeGroup:  "45-54",
		Dose:      "B",
		Count:     50,
	}, result[6])

	testServer.Fail = true
	_, err = c.GetVaccinations(ctx)
	require.Error(t, err)

	apiServer.Close()
	_, err = c.GetVaccinations(ctx)
	require.Error(t, err)

}

func TestAPIVaccinationsResponseEntry_GetTimestamp(t *testing.T) {
	timestamp := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	entry := sciensano.APIVaccinationsResponse{
		TimeStamp: sciensano.TimeStamp{Time: timestamp},
	}

	assert.Equal(t, timestamp, entry.GetTimestamp())
}

func TestAPIVaccinationsResponseEntry_GetGroupFieldValue(t *testing.T) {
	entry := sciensano.APIVaccinationsResponse{
		Manufacturer: "Moderna",
		Region:       "Flanders",
		AgeGroup:     "85+",
	}

	assert.Equal(t, "Flanders", entry.GetGroupFieldValue(apiclient.GroupByRegion))
	assert.Equal(t, "85+", entry.GetGroupFieldValue(apiclient.GroupByAgeGroup))
	assert.Equal(t, "Moderna", entry.GetGroupFieldValue(apiclient.GroupByManufacturer))
}

func BenchmarkClient_GetVaccinations(b *testing.B) {
	testServer := httptest.NewServer(http.HandlerFunc(handleVaccinationResponse))
	defer testServer.Close()

	c := sciensano.Client{
		Caller: &client.BaseClient{HTTPClient: http.DefaultClient},
		URL:    testServer.URL,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := c.GetVaccinations(context.Background())
		if err != nil {
			b.Fatal(err)
		}
	}
}

var bigVaccinationsFile []byte

func handleVaccinationResponse(w http.ResponseWriter, _ *http.Request) {
	var err error
	if bigVaccinationsFile == nil {
		bigVaccinationsFile, err = os.ReadFile("../../data/vaccinations.json")
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(bigVaccinationsFile)
}

func TestClient_Vaccination_Measurement(t *testing.T) {
	entry := sciensano.APIVaccinationsResponse{
		TimeStamp: sciensano.TimeStamp{Time: time.Now()},
		Region:    "Flanders",
		AgeGroup:  "85+",
		Dose:      "A",
		Count:     10,
	}

	assert.NotZero(t, entry.GetTimestamp())
	assert.Equal(t, "Flanders", entry.GetGroupFieldValue(apiclient.GroupByRegion))
	assert.Empty(t, entry.GetGroupFieldValue(apiclient.GroupByProvince))
	assert.Equal(t, "85+", entry.GetGroupFieldValue(apiclient.GroupByAgeGroup))
	assert.Equal(t, 10.0, entry.GetTotalValue())
	assert.Equal(t, []string{"partial", "full", "singledose", "booster"}, entry.GetAttributeNames())
	assert.Equal(t, []float64{10, 0, 0, 0}, entry.GetAttributeValues())

	entry.Dose = "B"
	assert.Equal(t, []float64{0, 10, 0, 0}, entry.GetAttributeValues())
	entry.Dose = "C"
	assert.Equal(t, []float64{0, 0, 10, 0}, entry.GetAttributeValues())
	entry.Dose = "E"
	assert.Equal(t, []float64{0, 0, 0, 10}, entry.GetAttributeValues())
}
