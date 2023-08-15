package sciensano

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"path/filepath"
	"testing"
	"time"
)

func TestFetcher_Cases(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(handler))
	defer s.Close()

	f := Fetcher[Cases]{
		Target: MustGetURL(s.URL, CasesEndpoint),
		Client: http.DefaultClient,
	}

	assert.Equal(t, "COVID19BE_CASES_AGESEX.json", filepath.Base(f.Target))

	ctx := context.Background()

	timestamp, err := f.GetLastModified(ctx)
	require.NoError(t, err)
	assert.NotZero(t, timestamp)

	entries, err := f.Fetch(context.Background())
	require.NoError(t, err)
	assert.NotEmpty(t, entries)
}

func TestFetcher_Vaccinations(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(handler))
	defer s.Close()

	f := Fetcher[Vaccinations]{
		Target: MustGetURL(s.URL, VaccinationsEndpoint),
		Client: http.DefaultClient,
	}

	assert.Equal(t, "COVID19BE_VACC.json", filepath.Base(f.Target))

	ctx := context.Background()

	timestamp, err := f.GetLastModified(ctx)
	require.NoError(t, err)
	assert.NotZero(t, timestamp)

	entries, err := f.Fetch(context.Background())
	require.NoError(t, err)
	assert.NotEmpty(t, entries)
}

func TestFetcher_TestResults(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(handler))
	defer s.Close()

	f := Fetcher[TestResults]{
		Target: MustGetURL(s.URL, TestResultsEndpoint),
		Client: http.DefaultClient,
	}

	assert.Equal(t, "COVID19BE_tests.json", filepath.Base(f.Target))

	ctx := context.Background()

	timestamp, err := f.GetLastModified(ctx)
	require.NoError(t, err)
	assert.NotZero(t, timestamp)

	entries, err := f.Fetch(context.Background())
	require.NoError(t, err)
	assert.NotEmpty(t, entries)
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodHead {
		w.Header().Add("Last-Modified", time.Now().Format(time.RFC1123))
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "invalid http method", http.StatusBadRequest)
		return
	}
	endpoint, ok := getEndpoint(r.URL.Path)
	if !ok {
		http.Error(w, r.URL.Path+" not found", http.StatusNotFound)
		return
	}
	filename, ok := filenames[endpoint]
	if !ok {
		http.Error(w, "input file not found", http.StatusInternalServerError)
		return
	}
	content, err := os.ReadFile(path.Join("testutil", "testdata", filename))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(content)
}

var filenames = map[Endpoint]string{
	CasesEndpoint:            "cases.json",
	HospitalisationsEndpoint: "hospitalisations.json",
	MortalitiesEndpoint:      "mortalities.json",
	TestResultsEndpoint:      "testResults.json",
	VaccinationsEndpoint:     "vaccinations.json",
}

func getEndpoint(s string) (Endpoint, bool) {
	for endpoint, pathName := range routes {
		if s == pathName {
			return endpoint, true
		}
	}
	return -1, false
}
