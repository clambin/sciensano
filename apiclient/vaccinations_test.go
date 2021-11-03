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
