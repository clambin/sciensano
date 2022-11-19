package sciensano

import (
	"fmt"
	"strconv"
	"time"
)

// TimeStamp represents a timestamp in the API responder. Needed for parsing purposes
type TimeStamp struct {
	time.Time
}

// UnmarshalJSON unmarshals a TimeStamp from the API responder.
func (ts *TimeStamp) UnmarshalJSON(b []byte) (err error) {
	s := string(b)
	if len(s) != 12 || s[0] != '"' && s[11] != '"' {
		return fmt.Errorf("invalid timestamp: %s", s)
	}
	//var year, month, day int
	year, errYear := strconv.Atoi(s[1:5])
	month, errMonth := strconv.Atoi(s[6:8])
	day, errDay := strconv.Atoi(s[9:11])

	if errYear != nil || errMonth != nil || errDay != nil {
		return fmt.Errorf("invalid timestamp: %s", s)
	}
	ts.Time = time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	return
}

// MarshalJSON marshals a TimeStamp to JSON
func (ts *TimeStamp) MarshalJSON() ([]byte, error) {
	return []byte(`"` + ts.Time.Format(time.RFC3339)[:10] + `"`), nil
}
