package apiclient_test

import (
	"context"
	"github.com/clambin/sciensano/sciensano/apiclient"
	"github.com/clambin/sciensano/sciensano/apiclient/fake"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
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
	assert.NotZero(t, result[6].TimeStamp)

	testServer.Fail = true
	_, err = client.GetVaccinations(ctx)
	require.Error(t, err)

	apiServer.Close()
	_, err = client.GetVaccinations(ctx)
	require.Error(t, err)

}
