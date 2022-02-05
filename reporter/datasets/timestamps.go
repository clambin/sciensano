package datasets

import (
	"sort"
	"time"
)

type Timestamps struct {
	timestamps     []time.Time
	timestampIndex map[time.Time]int
}

func MakeTimestamps() *Timestamps {
	return &Timestamps{
		timestamps:     make([]time.Time, 0),
		timestampIndex: make(map[time.Time]int),
	}
}

func (t Timestamps) GetIndex(timestamp time.Time) int {
	return t.timestampIndex[timestamp]
}

func (t Timestamps) Count() int {
	return len(t.timestamps)
}

func (t Timestamps) List() (timestamps []time.Time) {
	return t.timestamps
}

func (t *Timestamps) Add(timestamp time.Time) (index int, added bool) {
	_, found := t.timestampIndex[timestamp]

	if found {
		return 0, false
	}

	added = true
	index = len(t.timestamps)
	t.timestampIndex[timestamp] = index

	mustSort := len(t.timestamps) > 0 && timestamp.After(t.timestamps[len(t.timestamps)-1]) == false
	t.timestamps = append(t.timestamps, timestamp)
	if mustSort {
		sort.Slice(t.timestamps, func(i, j int) bool { return t.timestamps[i].Before(t.timestamps[j]) })
	}

	return
}

func (t Timestamps) Copy() (clone *Timestamps) {
	clone = &Timestamps{
		timestamps:     make([]time.Time, len(t.timestamps)),
		timestampIndex: make(map[time.Time]int),
	}
	copy(clone.timestamps, t.timestamps)
	for key, val := range t.timestampIndex {
		clone.timestampIndex[key] = val
	}
	return
}
