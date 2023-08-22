package population

import (
	"context"
	"github.com/clambin/sciensano/v2/internal/population/bracket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log/slog"
	"path"
	"testing"
	"time"
)

func TestServer_GetByRegion(t *testing.T) {
	s := Server{Path: path.Join(tmpDir, "demographics.txt"), Logger: slog.Default()}
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
	s := Server{Path: path.Join(tmpDir, "demographics.txt"), Logger: slog.Default()}
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
		Path:     path.Join(tmpDir, "demographics.txt"),
		Interval: time.Hour,
		Logger:   slog.Default(),
	}

	ctx, cancel := context.WithCancel(context.Background())

	ch := make(chan error)
	go func() {
		ch <- s.Run(ctx)
	}()

	ctx2, cancel2 := context.WithTimeout(ctx, 5*time.Second)
	defer cancel2()
	assert.NoError(t, s.WaitTillReady(ctx2))

	assert.NotEmpty(t, s.GetByRegion())

	cancel()
	assert.NoError(t, <-ch)
}
