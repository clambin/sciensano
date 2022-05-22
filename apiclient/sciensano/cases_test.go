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

func TestClient_GetCases(t *testing.T) {
	testServer := fake.Handler{}
	apiServer := httptest.NewServer(http.HandlerFunc(testServer.Handle))

	c := sciensano.Client{
		URL:    apiServer.URL,
		Caller: &client.BaseClient{HTTPClient: http.DefaultClient},
	}

	ctx := context.Background()
	result, err := c.GetCases(ctx)

	require.NoError(t, err)
	require.Len(t, result, 2)

	assert.Equal(t, &sciensano.APICasesResponse{
		TimeStamp: sciensano.TimeStamp{Time: time.Date(2020, 3, 1, 0, 0, 0, 0, time.UTC)},
		Province:  "VlaamsBrabant",
		Region:    "Flanders",
		AgeGroup:  "40-49",
		Cases:     1,
	}, result[0])

	assert.Equal(t, &sciensano.APICasesResponse{
		TimeStamp: sciensano.TimeStamp{Time: time.Date(2020, 3, 1, 0, 0, 0, 0, time.UTC)},
		Province:  "Brussels",
		Region:    "Brussels",
		AgeGroup:  "40-49",
		Cases:     2,
	}, result[1])

	testServer.Fail = true
	_, err = c.GetCases(ctx)
	require.Error(t, err)

	apiServer.Close()
	_, err = c.GetCases(ctx)
	require.Error(t, err)
}

func TestClient_Case_Measurement(t *testing.T) {
	c := sciensano.APICasesResponse{
		TimeStamp: sciensano.TimeStamp{Time: time.Now()},
		Region:    "Flanders",
		Province:  "VlaamsBrabant",
		AgeGroup:  "85+",
		Cases:     100,
	}

	assert.NotZero(t, c.GetTimestamp())
	assert.Equal(t, "Flanders", c.GetGroupFieldValue(apiclient.GroupByRegion))
	assert.Equal(t, "VlaamsBrabant", c.GetGroupFieldValue(apiclient.GroupByProvince))
	assert.Equal(t, "85+", c.GetGroupFieldValue(apiclient.GroupByAgeGroup))
	assert.Equal(t, 100.0, c.GetTotalValue())
	assert.Equal(t, []string{"total"}, c.GetAttributeNames())
	assert.Equal(t, []float64{100}, c.GetAttributeValues())
}

func BenchmarkClient_GetCases(b *testing.B) {
	b.StopTimer()
	testServer := httptest.NewServer(http.HandlerFunc(handleCasesResponse))
	defer testServer.Close()

	c := sciensano.Client{
		Caller: &client.BaseClient{HTTPClient: http.DefaultClient},
		URL:    testServer.URL,
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, err := c.GetCases(context.Background())
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkClient_GetCases_Cached(b *testing.B) {
	b.StopTimer()
	testServer := httptest.NewServer(http.HandlerFunc(handleCasesResponse))
	defer testServer.Close()

	// using caller's Cacher has little added value: main processing time goes to json unmarshal. So, apiclient.Cache is more effective here.
	c := sciensano.Client{
		Caller: client.NewCacher(http.DefaultClient, "sciensano", client.Options{}, []client.CacheTableEntry{
			{Endpoint: "/.*", IsRegExp: true, Expiry: time.Hour},
		}, time.Minute, time.Hour),
		URL: testServer.URL,
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, err := c.GetCases(context.Background())
		if err != nil {
			b.Fatal(err)
		}
	}
}

var bigCasesFile []byte

func handleCasesResponse(w http.ResponseWriter, _ *http.Request) {
	var err error
	if bigCasesFile == nil {
		bigCasesFile, err = os.ReadFile("../../data/cases.json")
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(bigCasesFile)
}
