package demographics_test

import (
	"github.com/clambin/sciensano/demographics"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestBracket_String(t *testing.T) {
	bracket := demographics.Bracket{
		Low:  0,
		High: 11,
	}

	assert.Equal(t, "00-11", bracket.String())

	bracket.Low = 75
	bracket.High = math.Inf(+1)
	assert.Equal(t, "75+", bracket.String())
}
