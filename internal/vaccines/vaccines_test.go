package vaccines_test

import (
	"github.com/clambin/sciensano/internal/vaccines"
	"github.com/clambin/sciensano/internal/vaccines/mock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVaccines(t *testing.T) {
	srv := vaccines.New()
	srv.HTTPClient = mock.GetServer()

	batches, err := srv.GetBatches()

	if assert.Nil(t, err) && assert.Len(t, batches, 3) {
		assert.Equal(t, "A", batches[0].Manufacturer)
		assert.Equal(t, "B", batches[1].Manufacturer)
		assert.Equal(t, "C", batches[2].Manufacturer)

		accu := vaccines.AccumulateBatches(batches)

		if assert.Len(t, accu, 3) {
			assert.Equal(t, int64(300), accu[0].Amount)
			assert.Equal(t, int64(500), accu[1].Amount)
			assert.Equal(t, int64(600), accu[2].Amount)
		}
	}

}
