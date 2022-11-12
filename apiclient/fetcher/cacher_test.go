package fetcher_test

import (
	"context"
	"errors"
	"github.com/clambin/sciensano/apiclient"
	"github.com/clambin/sciensano/apiclient/fetcher"
	"github.com/clambin/sciensano/apiclient/fetcher/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
	"time"
)

func TestCacher_Fetch(t *testing.T) {
	ctx := context.Background()
	f := mocks.NewFetcher(t)
	f.On("DataTypes").Return(map[int]string{1: "test"})

	c := fetcher.NewCacher(f)
	timestamp := time.Now()

	f.On("GetLastUpdated", mock.AnythingOfType("*context.emptyCtx"), 1).Return(timestamp, nil).Once()
	f.On("Fetch", mock.AnythingOfType("*context.emptyCtx"), 1).Return(nil, errors.New("fail")).Once()
	_, err := c.Fetch(ctx, 1)
	assert.Error(t, err)

	f.On("GetLastUpdated", mock.AnythingOfType("*context.emptyCtx"), 1).Return(timestamp, nil).Once()
	f.On("Fetch", mock.AnythingOfType("*context.emptyCtx"), 1).Return([]apiclient.APIResponse{}, nil).Once()
	data, err := c.Fetch(ctx, 1)
	require.NoError(t, err)
	assert.Equal(t, []apiclient.APIResponse{}, data)

	_, err = c.Fetch(ctx, 2)
	assert.Error(t, err)
}

func TestCacher_GetLastUpdated(t *testing.T) {
	ctx := context.Background()
	f := mocks.NewFetcher(t)
	f.On("DataTypes").Return(map[int]string{1: "test"})

	c := fetcher.NewCacher(f)

	f.On("GetLastUpdated", mock.AnythingOfType("*context.emptyCtx"), 1).Return(time.Time{}, errors.New("fail")).Once()
	_, err := c.GetLastUpdated(ctx, 1)
	assert.Error(t, err)

	timestamp := time.Now()
	f.On("GetLastUpdated", mock.AnythingOfType("*context.emptyCtx"), 1).Return(timestamp, nil).Once()
	data, err := c.GetLastUpdated(ctx, 1)
	require.NoError(t, err)
	assert.Equal(t, timestamp, data)
}

func TestCacher_DataTypes(t *testing.T) {
	f := mocks.NewFetcher(t)
	f.On("DataTypes").Return(map[int]string{1: "test"})

	c := fetcher.NewCacher(f)

	data := c.DataTypes()
	assert.Equal(t, map[int]string{1: "test"}, data)
}

func TestCacher_AutoUpdate(t *testing.T) {
	f := mocks.NewFetcher(t)
	f.On("DataTypes").Return(map[int]string{1: "test"})

	c := fetcher.NewCacher(f)
	timestamp := time.Now()

	f.On("GetLastUpdated", mock.AnythingOfType("*context.cancelCtx"), 1).Return(timestamp, nil).Once()
	f.On("Fetch", mock.AnythingOfType("*context.cancelCtx"), 1).Return([]apiclient.APIResponse{}, nil).Once()

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(1)
	go func() { c.AutoUpdate(ctx, time.Hour); wg.Done() }()

	assert.Eventually(t, func() bool {
		data, err := c.Fetch(context.Background(), 1)
		return err == nil && data != nil
	}, time.Second, 10*time.Millisecond)

	cancel()
	wg.Wait()
}

func TestCacher_With_Limiter(t *testing.T) {
	f := mocks.NewFetcher(t)
	f.On("DataTypes").Return(map[int]string{1: "test"})
	timestamp := time.Now()

	c := fetcher.NewCacher(fetcher.NewLimiter(f, 1))
	f.On("GetLastUpdated", mock.AnythingOfType("*context.cancelCtx"), 1).Return(timestamp, nil).Once()
	f.On("Fetch", mock.AnythingOfType("*context.cancelCtx"), 1).Return([]apiclient.APIResponse{}, nil).Once()

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(1)
	go func() { c.AutoUpdate(ctx, time.Hour); wg.Done() }()

	assert.Eventually(t, func() bool {
		data, err := c.Fetch(context.Background(), 1)
		return err == nil && data != nil
	}, time.Second, 10*time.Millisecond)

	cancel()
	wg.Wait()
}
