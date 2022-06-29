package sciensano_test

import (
	"github.com/clambin/sciensano/apiclient/sciensano"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTimeStamp_MarshalJSON(t *testing.T) {
	tests := []struct {
		timestamp time.Time
		want      string
		wantErr   assert.ErrorAssertionFunc
	}{
		{timestamp: time.Date(2022, time.June, 29, 0, 0, 0, 0, time.UTC), want: `"2022-06-29"`, wantErr: assert.NoError},
	}
	for _, tt := range tests {
		ts := sciensano.TimeStamp{Time: tt.timestamp}
		got, err := ts.MarshalJSON()
		if !tt.wantErr(t, err, "MarshalJSON()") {
			return
		}
		assert.Equalf(t, tt.want, string(got), "MarshalJSON()")
	}
}

func TestTimeStamp_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		body    string
		want    sciensano.TimeStamp
		wantErr assert.ErrorAssertionFunc
	}{
		{body: `"2022-06-29"`, want: sciensano.TimeStamp{Time: time.Date(2022, time.June, 29, 0, 0, 0, 0, time.UTC)}, wantErr: assert.NoError},
		{body: `2022-06-29`, wantErr: assert.Error},
		{body: `"2022-AA-29"`, wantErr: assert.Error},
		{body: ``, wantErr: assert.Error},
	}
	for _, tt := range tests {
		var ts sciensano.TimeStamp
		err := ts.UnmarshalJSON([]byte(tt.body))
		if !tt.wantErr(t, err, "UnmarshalJSON()") {
			return
		}
		assert.Equalf(t, tt.want, ts, "UnmarshalJSON()")
	}
}
