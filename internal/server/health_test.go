package server_test

import (
	"github.com/clambin/sciensano/internal/server"
	"github.com/clambin/sciensano/internal/server/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer_Health(t *testing.T) {
	r := mocks.NewReportsStore(t)
	r.EXPECT().Keys().Return([]string{"foo", "bar"})
	s := server.New(r, slog.Default())

	req, _ := http.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	s.JSONServer.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, `{
  "DataSources": 6,
  "ReporterCache": [
    "foo",
    "bar"
  ]
}
`, w.Body.String())
}
