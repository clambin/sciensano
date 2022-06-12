package vaccines_test

import (
	"context"
	"github.com/clambin/go-metrics/client"
	"github.com/clambin/sciensano/apiclient"
	"github.com/clambin/sciensano/apiclient/cache"
	"github.com/clambin/sciensano/apiclient/vaccines"
	"github.com/clambin/sciensano/apiclient/vaccines/fake"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestClient_GetBatches(t *testing.T) {
	server := fake.Server{}
	apiServer := httptest.NewServer(http.HandlerFunc(server.Handler))

	c := vaccines.Client{
		Caller: &client.BaseClient{HTTPClient: http.DefaultClient},
		URL:    apiServer.URL,
	}

	batches, err := c.GetBatches(context.Background())
	require.NoError(t, err)
	require.Len(t, batches, 3)
	assert.Equal(t, 300, batches[0].(*vaccines.APIBatchResponse).Amount)
	assert.Equal(t, 200, batches[1].(*vaccines.APIBatchResponse).Amount)
	assert.Equal(t, 100, batches[2].(*vaccines.APIBatchResponse).Amount)

	server.Fail = true
	_, err = c.GetBatches(context.Background())
	require.Error(t, err)

	apiServer.Close()
	_, err = c.GetBatches(context.Background())
	require.Error(t, err)
}

func BenchmarkClient_GetBatches(b *testing.B) {
	server := fake.Server{}
	apiServer := httptest.NewServer(http.HandlerFunc(server.Handler))

	c := vaccines.Client{
		Caller: &client.BaseClient{HTTPClient: http.DefaultClient},
		URL:    apiServer.URL,
	}

	for i := 0; i < b.N; i++ {
		_, err := c.GetBatches(context.Background())
		if err != nil {
			b.Fatal(err)
		}
	}
}

func TestBatch_Measurement(t *testing.T) {
	b := vaccines.APIBatchResponse{
		Date:         vaccines.Timestamp{Time: time.Now()},
		Manufacturer: "A",
		Amount:       200,
	}

	assert.NotZero(t, b.GetTimestamp())
	assert.Empty(t, b.GetGroupFieldValue(apiclient.GroupByAgeGroup))
	assert.Equal(t, "A", b.GetGroupFieldValue(apiclient.GroupByManufacturer))
	assert.Equal(t, 200.0, b.GetTotalValue())
	assert.Equal(t, []string{"total"}, b.GetAttributeNames())
	assert.Equal(t, []float64{200}, b.GetAttributeValues())
}

func TestClient_Refresh(t *testing.T) {
	server := fake.Server{}
	apiServer := httptest.NewServer(http.HandlerFunc(server.Handler))
	defer apiServer.Close()

	c := vaccines.Client{
		Caller: &client.BaseClient{HTTPClient: http.DefaultClient},
		URL:    apiServer.URL,
	}

	ch := make(chan cache.FetcherResponse)
	go c.Fetch(context.Background(), ch)

	response := <-ch
	assert.Equal(t, "Vaccines", response.Name)
}
