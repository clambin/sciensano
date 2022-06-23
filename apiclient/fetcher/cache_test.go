package fetcher

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCacheEntry_ShouldPoll(t *testing.T) {
	type testCase struct {
		updated time.Duration
		checked time.Duration
		poll    bool
	}

	const (
		gracePeriod = 5 * time.Minute
		expiration  = time.Hour
	)

	testCases := []testCase{
		{updated: time.Minute, checked: time.Minute, poll: false},
		{updated: 30 * time.Minute, checked: time.Minute, poll: false},
		{updated: 30 * time.Minute, checked: 10 * time.Minute, poll: true},
		{updated: 2 * time.Hour, checked: time.Minute, poll: true},
	}

	for idx, tc := range testCases {
		now := time.Now()
		entry := cacheEntry{
			updated: now.Add(-tc.updated),
			checked: now.Add(-tc.checked),
		}

		assert.Equal(t, tc.poll, entry.shouldPoll(gracePeriod, expiration), idx)
	}
}
