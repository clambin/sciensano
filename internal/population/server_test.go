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
	err = s.update()
	assert.NoError(t, err)

	testCases := []struct {
		region string
		want   int
	}{
		{
			region: "Ostbelgien",
			want:   ostbelgienPopulation,
		},
		{
			region: "Wallonia",
			want:   -61704,
		},
		{
			region: "Brussels",
			want:   87835,
		},
		{
			region: "Flanders",
			want:   2972,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.region, func(t *testing.T) {
			population := s.GetForRegion(tt.region)
			assert.Equal(t, tt.want, population)
		})
	}
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
		count := s.GetForAgeBracket(testCase.arguments)
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

	assert.NotZero(t, s.GetForRegion("Ostbelgien"))

	cancel()
	assert.NoError(t, <-ch)
}
