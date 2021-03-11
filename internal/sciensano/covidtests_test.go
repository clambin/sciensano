package sciensano_test

import (
	"github.com/clambin/sciensano/internal/sciensano"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGetTests(t *testing.T) {
	client := sciensano.APIClient{}
	firstDay := time.Date(2020, 03, 10, 0, 0, 0, 0, time.UTC)
	result, err := client.GetTests(firstDay)

	if assert.Nil(t, err) {
		assert.Len(t, result, 10)
		assert.Equal(t, firstDay, result[9].Timestamp)
		assert.Equal(t, 804, result[9].Total)
		assert.Equal(t, 57, result[9].Positive)
	}
}
