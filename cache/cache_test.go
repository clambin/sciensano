package cache

import (
	"context"
	"github.com/clambin/sciensano/cache/sciensano"
	"github.com/go-http-utils/headers"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func TestNewSciensanoCache(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(handler))
	defer s.Close()

	c := NewSciensanoCache(s.URL)
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		c.AutoRefresh(ctx, time.Hour)
	}()

	assert.Eventually(t, func() bool {
		return len(c.Vaccinations.Get()) > 0
	}, time.Minute, 100*time.Millisecond)

	assert.NotZero(t, len(c.Cases.Get()))
	assert.NotZero(t, len(c.Hospitalisations.Get()))
	assert.NotZero(t, len(c.Mortalities.Get()))
	assert.NotZero(t, len(c.TestResults.Get()))
	assert.NotZero(t, len(c.Vaccinations.Get()))

	cancel()
	wg.Wait()
}

func handler(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodHead {
		w.Header().Set(headers.LastModified, time.Now().Format(time.RFC1123))
		return
	}
	if req.Method != http.MethodGet {
		http.Error(w, "invalid method", http.StatusMethodNotAllowed)
		return
	}

	reqPath := req.URL.Path
	var reqName string
	for name, path := range sciensano.Routes {
		if path == reqPath {
			reqName = name
			break
		}
	}
	if reqName == "" {
		http.Error(w, "invalid route", http.StatusNotFound)
		return
	}

	f, err := os.Open(filepath.Join("sciensano", "input", reqName+".json"))
	if err != nil {
		http.Error(w, "could not open file "+reqName+": "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer func() { err = f.Close() }()

	_, err = io.Copy(w, f)
	if err != nil {
		http.Error(w, "could not open file "+reqName+": "+err.Error(), http.StatusInternalServerError)
		return
	}
}
