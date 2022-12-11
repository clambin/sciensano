package cache

import (
	"context"
	mockFetcher "github.com/clambin/sciensano/cache/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

/*
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
			expiry:  time.Hour,
			Fetcher: &fetcher[response]{client: c},
		}

		ctx := context.Background()
		s.refresh(ctx)

		result := s.Get(ctx)
		require.Len(t, result, 1)
		assert.Equal(t, "foo", result[0].Name)
	}
*/
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

	// First refresh: cacher should check last modified data & fetch new data
	lastModified := time.Now()
	f.On("GetLastModified", mock.AnythingOfType("*context.emptyCtx")).Return(lastModified, nil).Once()
	f.On("Fetch", mock.AnythingOfType("*context.emptyCtx")).Return(responses{response{Name: "foo"}}, nil).Once()
	s.refresh(ctx)
	assert.Equal(t, responses{response{Name: "foo"}}, s.Get(ctx))

	// Next refresh: expiry hasn't passed, so no calls should be made
	s.refresh(ctx)

	// Fake expiry. GetLastModified should be called. lastModified isn't changed, so Fetch should not be called
	s.lastChecked = time.Time{}
	f.On("GetLastModified", mock.AnythingOfType("*context.emptyCtx")).Return(lastModified, nil).Once()
	s.refresh(ctx)

	// Fake expire + update lastModified. Both GetLastModified and Fetch should be called
	s.lastChecked = time.Time{}
	lastModified = time.Now()
	f.On("GetLastModified", mock.AnythingOfType("*context.emptyCtx")).Return(lastModified, nil).Once()
	f.On("Fetch", mock.AnythingOfType("*context.emptyCtx")).Return(responses{response{Name: "bar"}}, nil).Once()
	s.refresh(ctx)
	assert.Equal(t, responses{response{Name: "bar"}}, s.Get(ctx))
}

func TestJitter(t *testing.T) {
	target := 100 * time.Minute
	for i := 0; i < 1000; i++ {
		realInterval := jitter(target)
		require.GreaterOrEqual(t, realInterval, 99*time.Minute)
		require.LessOrEqual(t, realInterval, 101*time.Minute)
	}
}
