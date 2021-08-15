package sciensano_test

import (
	"context"
	"github.com/clambin/sciensano/sciensano"
	"github.com/clambin/sciensano/sciensano/mock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetTests(t *testing.T) {
	testServer := mock.Handler{}
	apiServer := httptest.NewServer(http.HandlerFunc(testServer.Handle))
	defer apiServer.Close()

	client := sciensano.NewClient(time.Hour)
	client.SetURL(apiServer.URL)

	firstDay := time.Date(2021, 03, 10, 0, 0, 0, 0, time.UTC)
	result, err := client.GetTests(context.Background(), firstDay)

	if assert.Nil(t, err) && assert.Len(t, result, 2) {
		assert.Equal(t, firstDay, result[1].Timestamp)
		assert.Equal(t, 11, result[1].Total)
		assert.Equal(t, 5, result[1].Positive)
	}

	_, err = client.GetTests(context.Background(), firstDay)
	assert.NoError(t, err)

	assert.Equal(t, 1, testServer.Count)

}
