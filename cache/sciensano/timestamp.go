package sciensano

import (
	"fmt"
	"strconv"
	"time"
)

// TimeStamp represents a timestamp in the API response. Needed for parsing purposes
//
//easyjson:skip
type TimeStamp struct {
	time.Time
}

// UnmarshalJSON unmarshals a TimeStamp from the API responder.
func (ts *TimeStamp) UnmarshalJSON(b []byte) error {
	year, month, day, err := parseDate(b)
	if err == nil {
		ts.Time = time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	}
	return err
}

func parseDate(b []byte) (int, int, int, error) {
	if len(b) != 12 || b[0] != '"' && b[11] != '"' {
		return 0, 0, 0, fmt.Errorf("invalid timestamp: %s", b)
	}
	year, errYear := strconv.Atoi(string(b[1:5]))
	month, errMonth := strconv.Atoi(string(b[6:8]))
	day, errDay := strconv.Atoi(string(b[9:11]))
	if errYear != nil || errMonth != nil || errDay != nil {
		return 0, 0, 0, fmt.Errorf("invalid timestamp: %s", b)
	}
	return year, month, day, nil
}

// MarshalJSON marshals a TimeStamp to JSON
func (ts *TimeStamp) MarshalJSON() ([]byte, error) {
	return []byte(`"` + ts.Time.Format(time.RFC3339)[:10] + `"`), nil
}
