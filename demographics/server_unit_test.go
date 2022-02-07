package demographics

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStore_IsUpdated(t *testing.T) {
	s := Server{}

	mtime, updated, err := s.isUpdated()
	require.Error(t, err)

	s.Path = "server_unit_test.go"

	mtime, updated, err = s.isUpdated()
	require.NoError(t, err)
	require.True(t, updated)
	assert.NotZero(t, mtime)

	s.mtime = mtime
	_, updated, err = s.isUpdated()
	require.NoError(t, err)
	assert.False(t, updated)
}
