package cache

import (
	"bytes"
	"context"
	"github.com/clambin/sciensano/cache/sciensano"
	"github.com/go-http-utils/headers"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
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

	r := prometheus.NewRegistry()
	c := NewSciensanoCache(s.URL)
	r.MustRegister(c)
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		c.AutoRefresh(ctx, time.Hour)
	}()

	assert.Eventually(t, func() bool {
		return len(c.Vaccinations.Get(ctx)) > 0
	}, time.Minute, 100*time.Millisecond)

	assert.NotZero(t, len(c.Cases.Get(ctx)))
	assert.NotZero(t, len(c.Hospitalisations.Get(ctx)))
	assert.NotZero(t, len(c.Mortalities.Get(ctx)))
	assert.NotZero(t, len(c.TestResults.Get(ctx)))
	assert.NotZero(t, len(c.Vaccinations.Get(ctx)))

	cancel()
	wg.Wait()

	assert.NoError(t, testutil.CollectAndCompare(c, bytes.NewBufferString(`
# HELP sciensano_api_errors_total Number of failed Reporter API calls
# TYPE sciensano_api_errors_total counter
sciensano_api_errors_total{application="sciensano",method="GET",path="/Data/COVID19BE_CASES_AGESEX.json"} 0
sciensano_api_errors_total{application="sciensano",method="GET",path="/Data/COVID19BE_HOSP.json"} 0
sciensano_api_errors_total{application="sciensano",method="GET",path="/Data/COVID19BE_MORT.json"} 0
sciensano_api_errors_total{application="sciensano",method="GET",path="/Data/COVID19BE_VACC.json"} 0
sciensano_api_errors_total{application="sciensano",method="GET",path="/Data/COVID19BE_tests.json"} 0
sciensano_api_errors_total{application="sciensano",method="HEAD",path="/Data/COVID19BE_CASES_AGESEX.json"} 0
sciensano_api_errors_total{application="sciensano",method="HEAD",path="/Data/COVID19BE_HOSP.json"} 0
sciensano_api_errors_total{application="sciensano",method="HEAD",path="/Data/COVID19BE_MORT.json"} 0
sciensano_api_errors_total{application="sciensano",method="HEAD",path="/Data/COVID19BE_VACC.json"} 0
sciensano_api_errors_total{application="sciensano",method="HEAD",path="/Data/COVID19BE_tests.json"} 0
# HELP sciensano_api_max_inflight Maximum number of requests in flight
# TYPE sciensano_api_max_inflight gauge
sciensano_api_max_inflight{application="sciensano"} 3
`), "sciensano_api_errors_total", "sciensano_api_max_inflight"))
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
