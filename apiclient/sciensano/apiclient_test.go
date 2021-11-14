package sciensano_test

import (
	"context"
	"github.com/clambin/sciensano/apiclient/sciensano"
	"github.com/clambin/sciensano/measurement"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func TestClient_AutoRefresh(t *testing.T) {
	server := &Server{}
	testServer := httptest.NewServer(http.HandlerFunc(server.handle))
	defer testServer.Close()

	client := sciensano.Client{
		URL:        testServer.URL,
		HTTPClient: &http.Client{},
		Cache:      measurement.Cache{Retention: 15 * time.Minute},
	}
	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		client.AutoRefresh(ctx)
		wg.Done()
	}()

	assert.Eventually(t, func() bool { return len(server.calls()) == 5 }, 500*time.Millisecond, 10*time.Millisecond)
	for key, value := range server.calls() {
		assert.Equal(t, 1, value, key)
	}

	cancel()
	wg.Wait()
}

type Server struct {
	lock  sync.Mutex
	paths map[string]int
}

func (s *Server) init() {
	if s.paths == nil {
		s.paths = make(map[string]int)
	}
}

func (s *Server) handle(w http.ResponseWriter, req *http.Request) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.init()
	count, _ := s.paths[req.URL.Path]
	s.paths[req.URL.Path] = count + 1
	_, _ = w.Write([]byte(`[]`))
}

func (s *Server) calls() (calls map[string]int) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.init()
	calls = make(map[string]int)
	for key, value := range s.paths {
		calls[key] = value
	}
	return
}

func TestTimeStamp_UnmarshalJSON(t *testing.T) {
	testCases := []struct {
		input  []byte
		pass   bool
		output sciensano.TimeStamp
	}{
		{input: []byte(`"2021-10-06"`), pass: true, output: sciensano.TimeStamp{Time: time.Date(2021, 10, 6, 0, 0, 0, 0, time.UTC)}},
		{input: []byte(`"2021-13-06"`), pass: true, output: sciensano.TimeStamp{Time: time.Date(2022, 1, 6, 0, 0, 0, 0, time.UTC)}},
		{input: []byte(`"2021-09-31"`), pass: true, output: sciensano.TimeStamp{Time: time.Date(2021, 10, 1, 0, 0, 0, 0, time.UTC)}},
		{input: []byte(`2021-10-06`), pass: false},
		{input: []byte(``), pass: false},
		{input: []byte(`""`), pass: false},
		{input: []byte(`"2021-10"`), pass: false},
		{input: []byte(`"2021-AA-06"`), pass: false},
	}
	var ts sciensano.TimeStamp

	for _, testCase := range testCases {
		err := ts.UnmarshalJSON(testCase.input)
		if testCase.pass {
			assert.NoError(t, err, string(testCase.input))
			assert.Equal(t, testCase.output, ts, string(testCase.input))
		} else {
			assert.Error(t, err, string(testCase.input))
		}
	}
}

func BenchmarkTimeStamp_UnmarshalJSON(b *testing.B) {
	ts := &sciensano.TimeStamp{}

	for i := 0; i < 1000000; i++ {
		_ = ts.UnmarshalJSON([]byte("\"2021-03-02\""))
	}
}
