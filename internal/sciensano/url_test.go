package sciensano_test

import (
	"github.com/clambin/sciensano/v2/internal/sciensano"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetURL(t *testing.T) {
	url := sciensano.MustGetURL("", sciensano.CasesEndpoint)
	assert.Equal(t, "https://epistat.sciensano.be/Data/COVID19BE_CASES_AGESEX.json", url)

	url = sciensano.MustGetURL("https://localhost", sciensano.CasesEndpoint)
	assert.Equal(t, "https://localhost/Data/COVID19BE_CASES_AGESEX.json", url)

	assert.Panics(t, func() {
		_ = sciensano.MustGetURL("", -1)
	})
}
