package store_test

import (
	"github.com/clambin/go-common/tabulator"
	"github.com/clambin/sciensano/v2/internal/reports/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log/slog"
	"testing"
)

func TestStore(t *testing.T) {
	s := store.Store{Logger: slog.Default()}

	_, err := s.Get("foo")
	assert.ErrorIs(t, err, store.ErrNotFound)

	s.Put("foo", tabulator.New("A"))

	report, err := s.Get("foo")
	require.NoError(t, err)
	assert.Equal(t, []string{"A"}, report.GetColumns())

	_, err = s.Get("bar")
	assert.ErrorIs(t, err, store.ErrNotFound)

	assert.Equal(t, []string{"foo"}, s.Keys())
}
