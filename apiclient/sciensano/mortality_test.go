package sciensano_test

import (
	"context"
	"github.com/clambin/go-metrics/caller"
	"github.com/clambin/sciensano/apiclient"
	"github.com/clambin/sciensano/apiclient/sciensano"
	"github.com/clambin/sciensano/apiclient/sciensano/fake"
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

	client := &sciensano.Client{
		URL:    apiServer.URL,
		Caller: &caller.BaseClient{HTTPClient: http.DefaultClient},
	}

	ctx := context.Background()
	result, err := client.GetMortality(ctx)

	require.NoError(t, err)

	assert.Equal(t, &sciensano.APIMortalityResponse{
		TimeStamp: sciensano.TimeStamp{Time: time.Date(2020, time.March, 10, 0, 0, 0, 0, time.UTC)}, Region: "Brussels", AgeGroup: "85+", Deaths: 1,
	}, result[0])

	assert.Equal(t, &sciensano.APIMortalityResponse{
		TimeStamp: sciensano.TimeStamp{Time: time.Date(2020, time.March, 10, 0, 0, 0, 0, time.UTC)}, Region: "Brussels", AgeGroup: "85+", Deaths: 2,
	}, result[1])

	testServer.Fail = true
	_, err = client.GetMortality(ctx)
	require.Error(t, err)

	apiServer.Close()
	_, err = client.GetMortality(ctx)
	require.Error(t, err)
}

func TestClient_Mortality_Measurement(t *testing.T) {
	entry := sciensano.APIMortalityResponse{
		TimeStamp: sciensano.TimeStamp{Time: time.Now()},
		Region:    "Flanders",
		AgeGroup:  "85+",
		Deaths:    10,
	}

	assert.NotZero(t, entry.GetTimestamp())
	assert.Equal(t, "Flanders", entry.GetGroupFieldValue(apiclient.GroupByRegion))
	assert.Empty(t, entry.GetGroupFieldValue(apiclient.GroupByProvince))
	assert.Equal(t, "85+", entry.GetGroupFieldValue(apiclient.GroupByAgeGroup))
	assert.Equal(t, 10.0, entry.GetTotalValue())
	assert.Equal(t, []string{"total"}, entry.GetAttributeNames())
	assert.Equal(t, []float64{10}, entry.GetAttributeValues())
}
