package datasource

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_jitter(t *testing.T) {
	assert.Equal(t, 475*time.Second, jitter(500*time.Second, 0.1, 0))
	//assert.Equal(t, 500*time.Second, jitter(500*time.Second, 0.1, 0.5))
	assert.Equal(t, 525*time.Second, jitter(500*time.Second, 0.1, 1))

	assert.Equal(t, 490*time.Millisecond, jitter(500*time.Millisecond, 0.04, 0))
	//assert.Equal(t, 500*time.Millisecond, jitter(500*time.Millisecond, 0.1, 0.5))
	assert.Equal(t, 510*time.Millisecond, jitter(500*time.Millisecond, 0.04, 1))
}
