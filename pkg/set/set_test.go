package set_test

import (
	"github.com/clambin/sciensano/pkg/set"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSet_IsNew(t *testing.T) {
	var s set.Set
	assert.True(t, s.IsNew("foo"))
	assert.False(t, s.IsNew("foo"))
	assert.True(t, s.IsNew("bar"))
	assert.False(t, s.IsNew("bar"))
}
