package sciensano_test

import (
	"github.com/clambin/sciensano/sciensano"
	"github.com/clambin/sciensano/sciensano/server"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetTests(t *testing.T) {
	testServer := server.Handler{}
	apiServer := httptest.NewServer(http.HandlerFunc(testServer.Handle))

	client := sciensano.Client{CacheDuration: 1 * time.Hour, URL: apiServer.URL}
	firstDay := time.Date(2021, 03, 10, 0, 0, 0, 0, time.UTC)
	result, err := client.GetTests(firstDay)

	if assert.Nil(t, err) && assert.Len(t, result, 2) {
		assert.Equal(t, firstDay, result[1].Timestamp)
		assert.Equal(t, 11, result[1].Total)
		assert.Equal(t, 5, result[1].Positive)
	}
}
