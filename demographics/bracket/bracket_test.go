package bracket_test

import (
	"github.com/clambin/sciensano/demographics/bracket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math"
	"testing"
)

func TestBracket_String(t *testing.T) {
	b := bracket.Bracket{
		Low:  0,
		High: 11,
	}

	assert.Equal(t, "00-11", b.String())

	b.Low = 75
	b.High = math.Inf(+1)
	assert.Equal(t, "75+", b.String())
}

func TestBracketFromString(t *testing.T) {
	testCases := []struct {
		input    string
		expected bracket.Bracket
		pass     bool
	}{
		{input: "a-21", pass: false},
		{input: "-21", expected: bracket.Bracket{Low: math.Inf(-1), High: 21}, pass: true},
		{input: "21-", expected: bracket.Bracket{Low: 21, High: math.Inf(+1)}, pass: true},
		{input: "21+", expected: bracket.Bracket{Low: 21, High: math.Inf(+1)}, pass: true},
		{input: "21-65", expected: bracket.Bracket{Low: 21, High: 65}, pass: true},
		{input: "-", expected: bracket.Bracket{Low: math.Inf(-1), High: math.Inf(+1)}, pass: true},
		{input: "21-32-65", pass: false},
		{input: "21-a", pass: false},
	}

	for _, testCase := range testCases {
		output, err := bracket.FromString(testCase.input)

		if testCase.pass {
			require.NoError(t, err, testCase.input)
			assert.Equal(t, testCase.expected, output, testCase.input)
		} else {
			assert.Error(t, err, testCase.input)
		}
	}
}
