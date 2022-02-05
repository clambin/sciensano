package sciensano_test

import (
	"context"
	"github.com/clambin/sciensano/apiclient/sciensano"
	"github.com/clambin/sciensano/apiclient/sciensano/fake"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestClient_Update(t *testing.T) {
	server := &fake.Handler{}
	testServer := httptest.NewServer(http.HandlerFunc(server.Handle))
	defer testServer.Close()

	client := sciensano.Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
	}
	ctx := context.Background()
	results, err := client.Update(ctx)
	require.NoError(t, err)
	assert.Len(t, results, 5)
	assert.Contains(t, results, "TestResults")
	assert.Contains(t, results, "Vaccinations")
	assert.Contains(t, results, "Hospitalisations")
	assert.Contains(t, results, "Cases")
	assert.Contains(t, results, "Mortality")

	testServer.Close()

	_, err = client.Update(ctx)
	assert.Error(t, err)
}

func TestTimeStamp_UnmarshalJSON(t *testing.T) {
	testCases := []struct {
		input  []byte
		pass   bool
		output sciensano.TimeStamp
	}{
		{input: []byte(`"2021-10-06"`), pass: true, output: sciensano.TimeStamp{Time: time.Date(2021, 10, 6, 0, 0, 0, 0, time.UTC)}},
		{input: []byte(`"2021-13-06"`), pass: true, output: sciensano.TimeStamp{Time: time.Date(2022, 1, 6, 0, 0, 0, 0, time.UTC)}},
		{input: []byte(`"2021-09-31"`), pass: true, output: sciensano.TimeStamp{Time: time.Date(2021, 10, 1, 0, 0, 0, 0, time.UTC)}},
		{input: []byte(`2021-10-06`), pass: false},
		{input: []byte(``), pass: false},
		{input: []byte(`""`), pass: false},
		{input: []byte(`"2021-10"`), pass: false},
		{input: []byte(`"2021-AA-06"`), pass: false},
	}
	var ts sciensano.TimeStamp

	for _, testCase := range testCases {
		err := ts.UnmarshalJSON(testCase.input)
		if testCase.pass {
			assert.NoError(t, err, string(testCase.input))
			assert.Equal(t, testCase.output, ts, string(testCase.input))
		} else {
			assert.Error(t, err, string(testCase.input))
		}
	}
}

func BenchmarkTimeStamp_UnmarshalJSON(b *testing.B) {
	ts := &sciensano.TimeStamp{}

	for i := 0; i < b.N; i++ {
		_ = ts.UnmarshalJSON([]byte("\"2021-03-02\""))
	}
}
