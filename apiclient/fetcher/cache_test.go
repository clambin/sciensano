package fetcher

import (
	"context"
	"errors"
	"github.com/clambin/sciensano/apiclient"
	"github.com/clambin/sciensano/apiclient/fetcher/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
	"time"
)

func TestCache_Run(t *testing.T) {
	f := mocks.NewFetcher(t)
	f.On("Fetch", mock.AnythingOfType("*context.cancelCtx"), 1).Return([]apiclient.APIResponse{}, nil)
	c := &cache{
		dataType: 1,
		fetcher:  f,
		expiry:   time.Second,
	}
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(1)
	go func() { c.Run(ctx, time.Hour); wg.Done() }()

	assert.Eventually(t, func() bool {
		data := c.Get()
		return data != nil
	}, time.Second, 10*time.Millisecond)

	cancel()
	wg.Wait()
}

func TestCache_Refresh(t *testing.T) {
	f := mocks.NewFetcher(t)
	c := &cache{
		dataType: 1,
		fetcher:  f,
		expiry:   100 * time.Millisecond,
	}
	ctx := context.Background()

	timestamp := time.Now()
	f.On("GetLastUpdated", mock.AnythingOfType("*context.emptyCtx"), 1).Return(timestamp, nil).Once()
	f.On("Fetch", mock.AnythingOfType("*context.emptyCtx"), 1).Return([]apiclient.APIResponse{}, nil).Once()
	err := c.refresh(ctx)
	require.NoError(t, err)
	data := c.Get()
	assert.NotNil(t, data)

	err = c.refresh(ctx)
	require.NoError(t, err)
	data = c.Get()
	assert.NotNil(t, data)

	time.Sleep(200 * time.Millisecond)
	f.On("GetLastUpdated", mock.AnythingOfType("*context.emptyCtx"), 1).Return(timestamp, nil).Once()
	err = c.refresh(ctx)
	require.NoError(t, err)
	data = c.Get()
	assert.NotNil(t, data)

	time.Sleep(200 * time.Millisecond)
	timestamp = timestamp.Add(100 * time.Millisecond)
	f.On("GetLastUpdated", mock.AnythingOfType("*context.emptyCtx"), 1).Return(timestamp, nil).Once()
	f.On("Fetch", mock.AnythingOfType("*context.emptyCtx"), 1).Return([]apiclient.APIResponse{}, nil).Once()
	err = c.refresh(context.Background())
	require.NoError(t, err)
	data = c.Get()
	assert.NotNil(t, data)
}

func TestCache_Refresh_Errors(t *testing.T) {
	f := mocks.NewFetcher(t)
	c := &cache{
		dataType: 1,
		fetcher:  f,
		expiry:   100 * time.Millisecond,
	}
	ctx := context.Background()
	timestamp := time.Now()

	f.On("Fetch", mock.AnythingOfType("*context.emptyCtx"), 1).Return(nil, errors.New("fail")).Once()
	err := c.refresh(ctx)
	assert.Error(t, err)

	f.On("Fetch", mock.AnythingOfType("*context.emptyCtx"), 1).Return([]apiclient.APIResponse{}, nil).Once()
	err = c.refresh(ctx)
	assert.NoError(t, err)

	time.Sleep(200 * time.Millisecond)
	f.On("GetLastUpdated", mock.AnythingOfType("*context.emptyCtx"), 1).Return(time.Time{}, errors.New("fail")).Once()
	err = c.refresh(ctx)
	assert.Error(t, err)

	f.On("GetLastUpdated", mock.AnythingOfType("*context.emptyCtx"), 1).Return(timestamp, nil).Once()
	err = c.refresh(ctx)
	assert.NoError(t, err)
}
