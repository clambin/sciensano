package bracket

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Bracket indicates a range of ages by which demographics data is grouped
type Bracket struct {
	Low  float64
	High float64
}

// FromString constructs a bracket from a string
func FromString(input string) (output Bracket, err error) {
	if strings.HasPrefix(input, "-") {
		input = strings.TrimPrefix(input, "-")
		return makeBracket("", input)
	}
	if strings.HasSuffix(input, "-") {
		input = strings.TrimSuffix(input, "-")
		return makeBracket(input, "")
	}
	if strings.HasSuffix(input, "+") {
		input = strings.TrimSuffix(input, "+")
		return makeBracket(input, "")
	}
	values := strings.Split(input, "-")
	if len(values) != 2 {
		return output, fmt.Errorf("invalid bracket: only 2 entries supported")
	}
	return makeBracket(values[0], values[1])
}

func makeBracket(low, high string) (b Bracket, err error) {
	b.Low, err = convert(low, math.Inf(-1))
	if err != nil {
		return
	}
	b.High, err = convert(high, math.Inf(+1))
	return
}

func convert(value string, fallback float64) (output float64, err error) {
	if value == "" {
		return fallback, nil
	}
	var valueAsInt int
	valueAsInt, err = strconv.Atoi(value)
	return float64(valueAsInt), err
}

// String returns a string representation of a Bracket
func (b Bracket) String() string {
	if b.High == math.Inf(+1) {
		return fmt.Sprintf("%02.0f+", b.Low)
	}
	return fmt.Sprintf("%02.0f-%02.0f", b.Low, b.High)
}
