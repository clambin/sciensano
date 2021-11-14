package sciensano_test

import (
	"context"
	"github.com/clambin/sciensano/apiclient/sciensano"
	"github.com/clambin/sciensano/apiclient/sciensano/fake"
	"github.com/clambin/sciensano/measurement"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestClient_GetHospitalisations(t *testing.T) {
	testServer := fake.Handler{}
	apiServer := httptest.NewServer(http.HandlerFunc(testServer.Handle))

	client := sciensano.Client{
		URL:        apiServer.URL,
		HTTPClient: &http.Client{},
	}

	ctx := context.Background()
	result, err := client.GetHospitalisations(ctx)

	require.NoError(t, err)
	require.Len(t, result, 2)

	assert.Equal(t, &sciensano.APIHospitalisationsResponseEntry{
		TimeStamp:   sciensano.TimeStamp{Time: time.Date(2020, time.March, 15, 0, 0, 0, 0, time.UTC)},
		Province:    "Brussels",
		Region:      "Brussels",
		TotalIn:     58,
		TotalInICU:  11,
		TotalInResp: 8,
		TotalInECMO: 0,
	}, result[0])
	assert.Equal(t, &sciensano.APIHospitalisationsResponseEntry{
		TimeStamp:   sciensano.TimeStamp{Time: time.Date(2020, time.March, 15, 0, 0, 0, 0, time.UTC)},
		Province:    "VlaamsBrabant",
		Region:      "Flanders",
		TotalIn:     13,
		TotalInICU:  2,
		TotalInResp: 0,
		TotalInECMO: 1,
	}, result[1])

	testServer.Fail = true
	_, err = client.GetHospitalisations(ctx)
	require.Error(t, err)

	apiServer.Close()
	_, err = client.GetHospitalisations(ctx)
	require.Error(t, err)
}

func TestClient_Hospitalisation_Measurement(t *testing.T) {
	entry := sciensano.APIHospitalisationsResponseEntry{
		TimeStamp:   sciensano.TimeStamp{Time: time.Now()},
		Region:      "Flanders",
		Province:    "VlaamsBrabant",
		TotalIn:     100,
		TotalInICU:  10,
		TotalInResp: 5,
		TotalInECMO: 1,
	}

	assert.NotZero(t, entry.GetTimestamp())
	assert.Equal(t, "Flanders", entry.GetGroupFieldValue(measurement.GroupByRegion))
	assert.Equal(t, "VlaamsBrabant", entry.GetGroupFieldValue(measurement.GroupByProvince))
	assert.Empty(t, entry.GetGroupFieldValue(measurement.GroupByAgeGroup))
	assert.Equal(t, 100.0, entry.GetTotalValue())
	assert.Equal(t, []string{"in", "inICU", "inResp", "inECMO"}, entry.GetAttributeNames())
	assert.Equal(t, []float64{100, 10, 5, 1}, entry.GetAttributeValues())
}
