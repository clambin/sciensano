package sciensano_test

import (
	"github.com/clambin/gotools/httpstub"
	"github.com/clambin/sciensano/pkg/sciensano"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGetTests(t *testing.T) {
	client := sciensano.Client{CacheDuration: 1 * time.Hour, HTTPClient: httpstub.NewTestClient(server)}
	firstDay := time.Date(2021, 03, 10, 0, 0, 0, 0, time.UTC)
	result, err := client.GetTests(firstDay)

	if assert.Nil(t, err) && assert.Len(t, result, 2) {
		assert.Equal(t, firstDay, result[1].Timestamp)
		assert.Equal(t, 11, result[1].Total)
		assert.Equal(t, 5, result[1].Positive)
	}
}
