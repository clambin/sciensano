package sciensano_test

import (
	"github.com/clambin/sciensano/internal/sciensano"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGetVaccines(t *testing.T) {
	client := sciensano.APIClient{}
	firstDay := time.Date(2020, 12, 28, 0, 0, 0, 0, time.UTC)
	result, err := client.GetVaccines(firstDay)

	if assert.Nil(t, err) {
		if assert.Len(t, result, 1) {
			assert.Equal(t, firstDay, result[0].Timestamp)
			assert.Equal(t, 298, result[0].FirstDose)
			assert.Equal(t, 0, result[0].SecondDose)
		}
	}
}
