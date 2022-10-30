package sciensano_test

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"github.com/clambin/httpclient"
	"github.com/clambin/sciensano/apiclient/sciensano"
	"github.com/go-http-utils/headers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
	"time"
)

var update = flag.Bool("update", false, "update .golden files")

func TestClient_GetLastUpdates(t *testing.T) {
	s := &server{}
	testServer := httptest.NewServer(http.HandlerFunc(s.handle))

	c := sciensano.Client{Caller: &httpclient.InstrumentedClient{
		Application: "test",
	}}
	c.URL = testServer.URL

	ctx := context.Background()
	lastModified, err := c.GetLastUpdates(ctx, sciensano.TypeTestResults)
	require.NoError(t, err)
	assert.NotZero(t, lastModified)

	_, err = c.GetLastUpdates(ctx, -1)
	require.Error(t, err)

	testServer.Close()
	_, err = c.GetLastUpdates(ctx, sciensano.TypeTestResults)
	require.Error(t, err)
}

func TestClient_Fetch(t *testing.T) {
	s := &server{}
	testServer := httptest.NewServer(http.HandlerFunc(s.handle))
	defer testServer.Close()

	c := sciensano.Client{Caller: &httpclient.InstrumentedClient{
		Application: "test",
	}}
	c.URL = testServer.URL
	ctx := context.Background()

	for i := sciensano.TypeTestResults; i <= sciensano.TypeHospitalisations; i++ {
		_, err := c.Fetch(ctx, i)
		assert.NoError(t, err, i)
	}

	_, err := c.Fetch(ctx, -1)
	require.Error(t, err)

	testServer.Close()
	_, err = c.Fetch(ctx, sciensano.TypeHospitalisations)
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
		case "/Data/COVID19BE_tests.json":
			body, _ = json.Marshal(testResultsResponses)
		case "/Data/COVID19BE_CASES_AGESEX.json":
			body, _ = json.Marshal(casesResponses)
		case "/Data/COVID19BE_HOSP.json":
			body, _ = json.Marshal(hospitalisationResponses)
		case "/Data/COVID19BE_MORT.json":
			body, _ = json.Marshal(mortalityResponses)
		case "/Data/COVID19BE_VACC.json":
			body, _ = json.Marshal(vaccinationResponses)
		default:
			http.Error(w, "path not found", http.StatusNotFound)
			return
		}
		s.cache[req.URL.Path] = body
	}
	_, _ = w.Write(body)
}

func BenchmarkClient_Fetch(b *testing.B) {
	c := sciensano.Client{
		Caller: &bigResponseCaller{},
	}
	ctx := context.Background()
	_, _ = c.Fetch(ctx, sciensano.TypeVaccinations)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := c.Fetch(ctx, sciensano.TypeVaccinations)
		if err != nil {
			b.Fatal(err)
		}
	}
}

type bigResponseCaller struct {
	body []byte
	lock sync.Mutex
}

func (b *bigResponseCaller) Do(_ *http.Request) (resp *http.Response, err error) {
	b.lock.Lock()
	defer b.lock.Unlock()
	if len(b.body) == 0 {
		if b.body, err = os.ReadFile("../../data/vaccinations.json"); err != nil {
			panic(err)
		}
	}
	resp = &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBuffer(b.body)),
	}
	return
}

func TestClient_DataTypes(t *testing.T) {
	c := sciensano.Client{}
	dataTypes := c.DataTypes()

	for _, code := range []int{sciensano.TypeCases, sciensano.TypeHospitalisations, sciensano.TypeMortality, sciensano.TypeTestResults, sciensano.TypeVaccinations} {
		_, found := dataTypes[code]
		assert.True(t, found, code)
	}
}
