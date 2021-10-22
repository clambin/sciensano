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

	assert.NotZero(t, result[0].TimeStamp)
	assert.Equal(t, "VlaamsBrabant", result[0].Province)
	assert.Equal(t, 1, result[0].Cases)
	assert.Equal(t, "Brussels", result[1].Province)
	assert.Equal(t, 2, result[1].Cases)

	testServer.Fail = true
	_, err = client.GetCases(ctx)
	require.Error(t, err)

	apiServer.Close()
	_, err = client.GetCases(ctx)
	require.Error(t, err)
}
