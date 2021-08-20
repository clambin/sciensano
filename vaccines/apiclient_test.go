package vaccines_test

import (
	"context"
	"github.com/clambin/sciensano/vaccines"
	"github.com/clambin/sciensano/vaccines/fake"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_GetBatches(t *testing.T) {
	server := fake.Server{}
	apiServer := httptest.NewServer(http.HandlerFunc(server.Handler))

	client := vaccines.Client{
		URL:        apiServer.URL,
		HTTPClient: &http.Client{},
	}

	batches, err := client.GetBatches(context.Background())
	require.NoError(t, err)
	require.Len(t, batches, 3)
	assert.Equal(t, 300, batches[0].Amount)
	assert.Equal(t, 200, batches[1].Amount)
	assert.Equal(t, 100, batches[2].Amount)

	accumulated := vaccines.AccumulateBatches(batches)
	require.Len(t, accumulated, 3)
	assert.Equal(t, 300, accumulated[0].Amount)
	assert.Equal(t, 500, accumulated[1].Amount)
	assert.Equal(t, 600, accumulated[2].Amount)

	server.Fail = true
	_, err = client.GetBatches(context.Background())
	require.Error(t, err)

	apiServer.Close()
	_, err = client.GetBatches(context.Background())
	require.Error(t, err)
}

func BenchmarkClient_GetBatches(b *testing.B) {
	server := fake.Server{}
	apiServer := httptest.NewServer(http.HandlerFunc(server.Handler))

	client := vaccines.Client{
		URL:        apiServer.URL,
		HTTPClient: &http.Client{},
	}

	for i := 0; i < 5000; i++ {
		batches, err := client.GetBatches(context.Background())
		require.NoError(b, err)
		accumulated := vaccines.AccumulateBatches(batches)
		require.Len(b, accumulated, 3)
	}
}
