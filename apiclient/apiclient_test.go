package apiclient_test

import (
	"github.com/clambin/sciensano/apiclient"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTimeStamp_UnmarshalJSON(t *testing.T) {
	testCases := []struct {
		input  []byte
		pass   bool
		output apiclient.TimeStamp
	}{
		{input: []byte(`"2021-10-06"`), pass: true, output: apiclient.TimeStamp{Time: time.Date(2021, 10, 6, 0, 0, 0, 0, time.UTC)}},
		{input: []byte(`"2021-13-06"`), pass: true, output: apiclient.TimeStamp{Time: time.Date(2022, 1, 6, 0, 0, 0, 0, time.UTC)}},
		{input: []byte(`"2021-09-31"`), pass: true, output: apiclient.TimeStamp{Time: time.Date(2021, 10, 1, 0, 0, 0, 0, time.UTC)}},
		{input: []byte(`2021-10-06`), pass: false},
		{input: []byte(``), pass: false},
		{input: []byte(`""`), pass: false},
		{input: []byte(`"2021-10"`), pass: false},
		{input: []byte(`"2021-AA-06"`), pass: false},
	}
	var ts apiclient.TimeStamp

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
	ts := &apiclient.TimeStamp{}

	for i := 0; i < 1000000; i++ {
		_ = ts.UnmarshalJSON([]byte("\"2021-03-02\""))
	}
}
