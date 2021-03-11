package sciensano_test

import (
	"github.com/clambin/sciensano/internal/sciensano"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGetTests(t *testing.T) {
	client := sciensano.APIClient{}
	firstDay := time.Date(2020, 03, 01, 0, 0, 0, 0, time.UTC)
	result, err := client.GetTests(firstDay)

	if assert.Nil(t, err) {
		assert.Len(t, result, 1)
		assert.Equal(t, firstDay, result[0].Timestamp)
		assert.Equal(t, 82, result[0].Total)
		assert.Equal(t, 0, result[0].Positive)
	}
}
