package fetcher_test

import (
	"context"
	"errors"
	"github.com/clambin/sciensano/apiclient"
	"github.com/clambin/sciensano/apiclient/fetcher"
	"github.com/clambin/sciensano/apiclient/fetcher/mocks"
	"github.com/clambin/sciensano/apiclient/sciensano"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
	"time"
)

func TestCacher_GetLastUpdates(t *testing.T) {
	timestamp := time.Date(2022, time.June, 20, 0, 0, 0, 0, time.UTC)
	f := &mocks.Fetcher{}
	f.On("GetLastUpdates", mock.AnythingOfType("*context.emptyCtx"), 0).Return(timestamp, nil)
	c := fetcher.NewCacher(f)

	ctx := context.Background()
	ts, err := c.GetLastUpdates(ctx, 0)
	require.NoError(t, err)
	assert.Equal(t, timestamp, ts)
}

func TestCacher_Fetch(t *testing.T) {
	f := &mocks.Fetcher{}
	c := fetcher.NewCacher(f)

	ctx := context.Background()

	// first call triggers a fetch
	f.On("Fetch", mock.AnythingOfType("*context.emptyCtx"), 0).Return(nil, nil).Once()
	f.On("DataTypes").Return(map[int]string{0: "test"})
	_, err := c.Fetch(ctx, 0)
	require.NoError(t, err)

	// calls within grace period should be cached
	_, err = c.Fetch(ctx, 0)
	require.NoError(t, err)

	mock.AssertExpectationsForObjects(t, f)
}

func TestCacher_Fetch2(t *testing.T) {
	f := &mocks.Fetcher{}
	c := fetcher.NewCacher(f)
	c.GracePeriod = 0

	ctx := context.Background()

	timestamp := time.Now()

	// first call triggers a fetch
	f.On("Fetch", mock.AnythingOfType("*context.emptyCtx"), 0).Return(nil, nil).Once()
	f.On("DataTypes").Return(map[int]string{0: "test"})
	_, err := c.Fetch(ctx, 0)
	require.NoError(t, err)

	// call beyond grace period will trigger a getLastUpdates. If not recent, no call to fetch is triggered.
	f.On("GetLastUpdates", mock.AnythingOfType("*context.emptyCtx"), 0).Return(timestamp, nil).Once()
	_, err = c.Fetch(ctx, 0)
	require.NoError(t, err)

	// If data is more recent, fetch is triggered.
	f.On("GetLastUpdates", mock.AnythingOfType("*context.emptyCtx"), 0).Return(time.Now(), nil).Once()
	f.On("Fetch", mock.AnythingOfType("*context.emptyCtx"), 0).Return(nil, nil).Once()
	_, err = c.Fetch(ctx, 0)
	require.NoError(t, err)

	mock.AssertExpectationsForObjects(t, f)
}

func TestCacher_Failures(t *testing.T) {
	f := &mocks.Fetcher{}
	f.On("DataTypes").Return(map[int]string{0: "test"})

	c := fetcher.NewCacher(f)
	c.GracePeriod = 0
	ctx := context.Background()

	f.On("GetLastUpdates", mock.AnythingOfType("*context.emptyCtx"), 0).Return(time.Time{}, errors.New("fail")).Once()
	_, err := c.GetLastUpdates(ctx, 0)
	require.Error(t, err)

	timestamp := time.Now()
	f.On("Fetch", mock.AnythingOfType("*context.emptyCtx"), 0).Return(nil, errors.New("fail")).Once()
	_, err = c.Fetch(ctx, 0)
	require.Error(t, err)

	f.On("GetLastUpdates", mock.AnythingOfType("*context.emptyCtx"), 0).Return(time.Time{}, errors.New("fail")).Once()
	_, err = c.Fetch(ctx, 0)
	require.Error(t, err)

	f.On("GetLastUpdates", mock.AnythingOfType("*context.emptyCtx"), 0).Return(timestamp, nil).Once()
	f.On("Fetch", mock.AnythingOfType("*context.emptyCtx"), 0).Return(nil, errors.New("fail")).Once()
	_, err = c.Fetch(ctx, 0)
	require.Error(t, err)

	f.On("GetLastUpdates", mock.AnythingOfType("*context.emptyCtx"), 0).Return(timestamp, nil).Once()
	f.On("Fetch", mock.AnythingOfType("*context.emptyCtx"), 0).Return(nil, nil).Once()
	_, err = c.Fetch(ctx, 0)
	require.NoError(t, err)

	mock.AssertExpectationsForObjects(t, f)
}

func BenchmarkCacher_Fetch(b *testing.B) {
	f := &fakeClient{}
	c := fetcher.NewCacher(f)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := c.Fetch(ctx, sciensano.TypeVaccinations); err != nil {
			b.Fatal(err)
		}
	}
}

type fakeClient struct {
	response  []apiclient.APIResponse
	timestamp time.Time
	lock      sync.Mutex
}

var _ fetcher.Fetcher = &fakeClient{}

func (f *fakeClient) Fetch(_ context.Context, _ int) ([]apiclient.APIResponse, error) {
	f.lock.Lock()
	defer f.lock.Unlock()
	if len(f.response) == 0 {
		f.response = bigResponse()
	}
	return f.response, nil
}

func (f *fakeClient) GetLastUpdates(_ context.Context, _ int) (time.Time, error) {
	f.lock.Lock()
	defer f.lock.Unlock()
	if f.timestamp.IsZero() {
		f.timestamp = time.Now()
	}
	return f.timestamp, nil
}

func (f *fakeClient) DataTypes() map[int]string {
	return map[int]string{
		0: "test",
	}
}

func bigResponse() (output []apiclient.APIResponse) {
	timestamp := time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC)
	for d := 0; d < 2*365; d++ {
		for _, region := range []string{"Flanders", "Brussels", "Wallonia"} {
			for _, dose := range []string{"A", "B", "C", "E"} {
				output = append(output, &sciensano.APIVaccinationsResponse{
					TimeStamp:    sciensano.TimeStamp{Time: timestamp},
					Manufacturer: "",
					Region:       region,
					AgeGroup:     "",
					Gender:       "",
					Dose:         dose,
					Count:        1,
				})
			}
		}

		timestamp = timestamp.Add(24 * time.Hour)
	}
	return
}
