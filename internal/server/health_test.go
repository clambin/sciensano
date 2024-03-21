package server

import (
	"github.com/clambin/sciensano/v2/internal/server/mocks"
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
	s := New(r, nil, slog.Default())

	req, _ := http.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	s.JSONServer.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, `{
  "DataSources": [
    "cases",
    "hospitalisations",
    "mortalities",
    "tests",
    "vaccination-rate",
    "vaccinations"
  ],
  "ReporterCache": [
    "foo",
    "bar"
  ]
}
`, w.Body.String())
}
