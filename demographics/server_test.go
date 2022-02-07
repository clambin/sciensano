package demographics

import (
	"context"
	"github.com/clambin/sciensano/demographics/bracket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
	"time"
)

func TestServer_GetByRegion(t *testing.T) {
	s := Server{Path: "../data/demographics.txt"}
	err := s.update()
	require.NoError(t, err)

	regions := s.GetByRegion()
	assert.Len(t, regions, 4)

	count, ok := regions["Ostbelgien"]
	require.True(t, ok)
	assert.Equal(t, ostbelgienPopulation, count)

	_, ok = regions["invalid"]
	assert.False(t, ok)
}

func TestServer_GetByAgeBracket(t *testing.T) {
	s := Server{Path: "../data/demographics.txt"}
	err := s.update()
	require.NoError(t, err)

	testCases := []struct {
		arguments bracket.Bracket
		expected  int
	}{
		{
			arguments: bracket.Bracket{},
			expected:  107103,
		},
		{
			arguments: bracket.Bracket{Low: 21, High: 64},
			expected:  58199,
		},
		{
			arguments: bracket.Bracket{Low: 85},
			expected:  741,
		},
	}

	for _, testCase := range testCases {
		count := s.GetByAgeBracket(testCase.arguments)
		assert.Equal(t, testCase.expected, count)

	}
}

func TestStore_Run(t *testing.T) {
	s := Server{
		Path:     "../data/demographics.txt",
		Interval: time.Hour,
	}

	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		s.Run(ctx)
		wg.Done()
	}()

	assert.Eventually(t, func() bool {
		population := s.GetByRegion()
		return len(population) > 0
	}, 5*time.Second, 100*time.Millisecond)

	cancel()
	wg.Wait()
}
