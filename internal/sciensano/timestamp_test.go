package sciensano

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestTimeStamp_MarshalJSON(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		timestamp time.Time
		fail      bool
	}{
		{name: "valid", input: `"2022-11-23"`, timestamp: time.Date(2022, time.November, 23, 0, 0, 0, 0, time.UTC)},
		{name: "no quotes", input: `2022-11-23`, fail: true},
		{name: "invalid", input: `"20221123"`, fail: true},
		{name: "too long", input: `"2022-11-23T00:00:00"`, fail: true},
		{name: "empty", input: ``, fail: true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			var timestamp TimeStamp
			err := json.Unmarshal([]byte(tt.input), &timestamp)
			if tt.fail {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.timestamp, timestamp.Time)

			body, err := json.Marshal(&timestamp)
			require.NoError(t, err)
			assert.Equal(t, tt.input, string(body))
		})
	}
}
