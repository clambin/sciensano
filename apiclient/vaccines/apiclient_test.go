package vaccines_test

import (
	"context"
	"encoding/json"
	"github.com/clambin/httpclient"
	"github.com/clambin/sciensano/apiclient/vaccines"
	"github.com/go-http-utils/headers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func TestClient_GetLastUpdates(t *testing.T) {
	s := &server{}
	testServer := httptest.NewServer(http.HandlerFunc(s.handle))

	c := vaccines.Client{Caller: &httpclient.InstrumentedClient{Application: "test"}}
	c.URL = testServer.URL

	ctx := context.Background()
	lastModified, err := c.GetLastUpdated(ctx, vaccines.TypeBatches)
	require.NoError(t, err)
	assert.NotZero(t, lastModified)

	_, err = c.GetLastUpdated(ctx, -1)
	require.Error(t, err)

	testServer.Close()
	_, err = c.GetLastUpdated(ctx, vaccines.TypeBatches)
	require.Error(t, err)
}

func TestClient_Fetch(t *testing.T) {
	s := &server{}
	testServer := httptest.NewServer(http.HandlerFunc(s.handle))

	c := vaccines.Client{Caller: &httpclient.InstrumentedClient{
		Application: "test",
	}}
	c.URL = testServer.URL
	ctx := context.Background()

	for i := vaccines.TypeBatches; i <= vaccines.TypeBatches; i++ {
		entries, err := c.Fetch(ctx, i)
		assert.NoError(t, err, i)
		assert.NotEmpty(t, entries)
	}

	_, err := c.Fetch(ctx, -1)
	require.Error(t, err)

	testServer.Close()
	_, err = c.Fetch(ctx, vaccines.TypeBatches)
	require.Error(t, err)
}

type server struct {
	cache map[string][]byte
	lock  sync.Mutex
}

func (s *server) handle(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodHead {
		w.Header().Set(headers.LastModified, time.Now().Format(time.RFC1123))
		return
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	if s.cache == nil {
		s.cache = make(map[string][]byte)
	}

	body, found := s.cache[req.URL.Path]
	if !found {
		switch req.URL.Path {
		case "/api/v1/delivered.json":
			output := struct {
				Result struct {
					Delivered []*vaccines.APIBatchResponse `json:"delivered"`
				} `json:"result"`
			}{}
			output.Result.Delivered = batchResponse
			body, _ = json.Marshal(output)
		default:
			http.Error(w, "path not found", http.StatusNotFound)
			return
		}
		s.cache[req.URL.Path] = body
	}
	_, _ = w.Write(body)
}
