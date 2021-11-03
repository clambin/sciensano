package apiclient_test

import (
	"github.com/clambin/sciensano/apiclient"
	"testing"
)

/*
import (
	"context"
	"github.com/clambin/sciensano/apiclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestCache_GetCases_Real(t *testing.T) {
	c := apiclient.Client{HTTPClient: &http.Client{}}

	cases, err := c.GetCases(context.Background())
	require.NoError(t, err)

	for _, entry := range cases {
		if assert.Greater(t, entry.TimeStamp.Time.Year(), 2000) == false {
			assert.Zero(t, entry.TimeStamp)
		}
	}
}

*/

func BenchmarkTimeStamp_UnmarshalJSON(b *testing.B) {
	apiclient.ManualParse = true
	ts := &apiclient.TimeStamp{}

	for i := 0; i < 1000000; i++ {
		_ = ts.UnmarshalJSON([]byte("\"2021-03-02\""))
		// require.NoError(b, err)
	}
}

func BenchmarkTimeStamp_UnmarshalJSON_Old(b *testing.B) {
	apiclient.ManualParse = false
	ts := &apiclient.TimeStamp{}

	for i := 0; i < 1000000; i++ {
		_ = ts.UnmarshalJSON([]byte("\"2021-03-02\""))
		// require.NoError(b, err)
	}
}
