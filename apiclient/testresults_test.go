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

func TestClient_GetTestResults(t *testing.T) {
	testServer := fake.Handler{}
	apiServer := httptest.NewServer(http.HandlerFunc(testServer.Handle))

	client := apiclient.Client{
		URL:        apiServer.URL,
		HTTPClient: &http.Client{},
	}

	ctx := context.Background()
	result, err := client.GetTestResults(ctx)

	require.NoError(t, err)
	require.Len(t, result, 3)

	assert.NotZero(t, result[2].TimeStamp)
	assert.Equal(t, 15, result[2].Total)
	assert.Equal(t, 10, result[2].Positive)

	testServer.Fail = true
	_, err = client.GetTestResults(ctx)
	require.Error(t, err)

	apiServer.Close()
	_, err = client.GetTestResults(ctx)
	require.Error(t, err)
}
