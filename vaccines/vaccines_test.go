package vaccines_test

import (
	vaccines2 "github.com/clambin/sciensano/vaccines"
	mock2 "github.com/clambin/sciensano/vaccines/mock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVaccines(t *testing.T) {
	srv := vaccines2.New()
	srv.HTTPClient = mock2.GetServer()

	batches, err := srv.GetBatches()

	if assert.Nil(t, err) && assert.Len(t, batches, 3) {
		// assert.Equal(t, "A", batches[0].Manufacturer)
		// assert.Equal(t, "B", batches[1].Manufacturer)
		// assert.Equal(t, "C", batches[2].Manufacturer)

		accu := vaccines2.AccumulateBatches(batches)

		if assert.Len(t, accu, 3) {
			assert.Equal(t, 300, accu[0].Amount)
			assert.Equal(t, 500, accu[1].Amount)
			assert.Equal(t, 600, accu[2].Amount)
		}
	}

}
