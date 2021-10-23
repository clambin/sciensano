package apihandler_test

import (
	"github.com/clambin/sciensano/apihandler"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreate(t *testing.T) {
	h := apihandler.Create()

	assert.Len(t, h.GetHandlers(), 4)
}
