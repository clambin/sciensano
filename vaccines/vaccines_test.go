package vaccines_test

import (
	"github.com/clambin/sciensano/vaccines"
	"github.com/clambin/sciensano/vaccines/mock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestVaccines(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(mock.Handler))
	defer server.Close()

	client := vaccines.New()
	client.URL = server.URL

	batches, err := client.GetBatches()

	if assert.Nil(t, err) && assert.Len(t, batches, 3) {
		// assert.Equal(t, "A", batches[0].Manufacturer)
		// assert.Equal(t, "B", batches[1].Manufacturer)
		// assert.Equal(t, "C", batches[2].Manufacturer)

		accu := vaccines.AccumulateBatches(batches)

		if assert.Len(t, accu, 3) {
			assert.Equal(t, 300, accu[0].Amount)
			assert.Equal(t, 500, accu[1].Amount)
			assert.Equal(t, 600, accu[2].Amount)
		}
	}

}
