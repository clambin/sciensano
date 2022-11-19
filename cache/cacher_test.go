package cache

import (
	"bytes"
	"context"
	"github.com/clambin/httpclient/mocks"
	mockFetcher "github.com/clambin/sciensano/cache/mocks"
	"github.com/go-http-utils/headers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"sync"
	"testing"
	"time"
)

func TestCacher_Get(t *testing.T) {
	type response []struct {
		Name string
	}

	body := `[ { "name": "foo" } ]`
	resp := &http.Response{
		Status:        "OK",
		StatusCode:    http.StatusOK,
		Header:        map[string][]string{headers.LastModified: {time.Now().Format(time.RFC1123)}},
		Body:          io.NopCloser(bytes.NewBufferString(body)),
		ContentLength: int64(len(body)),
	}

	c := mocks.NewCaller(t)
	c.On("Do", mock.AnythingOfType("*http.Request")).Return(resp, nil).Twice()

	s := cacher[response]{
		Fetcher: &fetcher[response]{client: c},
	}

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		s.AutoRefresh(ctx, time.Second)
	}()

	assert.Eventually(t, func() bool {
		return len(s.Get()) > 0
	}, time.Second, 10*time.Millisecond)

	result := s.Get()
	require.Len(t, result, 1)
	assert.Equal(t, "foo", result[0].Name)

	cancel()
	wg.Wait()
}

func TestCacher_Refresh(t *testing.T) {
	type response struct {
		Name string
	}
	type responses []response

	f := mockFetcher.NewFetcher[responses](t)
	s := cacher[responses]{
		Fetcher: f,
		expiry:  time.Hour,
	}
	ctx := context.Background()

	// First call: cacher should check last modified data & fetch new data
	lastModified := time.Now()
	f.On("GetLastModified", mock.AnythingOfType("*context.emptyCtx")).Return(lastModified, nil).Once()
	f.On("Fetch", mock.AnythingOfType("*context.emptyCtx")).Return(responses{response{Name: "foo"}}, nil).Once()
	err := s.refresh(ctx)
	require.NoError(t, err)
	assert.Equal(t, responses{response{Name: "foo"}}, s.Get())

	// Next call: expiry hasn't passed, so no calls should be made
	err = s.refresh(ctx)
	require.NoError(t, err)

	// Fake expiry. GetLastModified should be called. lastModified isn't changed, so Fetch should not be called
	s.lastChecked = time.Time{}
	f.On("GetLastModified", mock.AnythingOfType("*context.emptyCtx")).Return(lastModified, nil).Once()
	err = s.refresh(ctx)
	require.NoError(t, err)

	// Fake expire + update lastModified. Both GetLastModified and Fetch should be called
	s.lastChecked = time.Time{}
	lastModified = time.Now()
	f.On("GetLastModified", mock.AnythingOfType("*context.emptyCtx")).Return(lastModified, nil).Once()
	f.On("Fetch", mock.AnythingOfType("*context.emptyCtx")).Return(responses{response{Name: "bar"}}, nil).Once()
	err = s.refresh(ctx)
	require.NoError(t, err)
	assert.Equal(t, responses{response{Name: "bar"}}, s.Get())
}
