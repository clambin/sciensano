package demographics

import (
	"fmt"
	"math"
)

// Bracket indicates a range of ages by which demographics data is grouped
type Bracket struct {
	Low  float64
	High float64
}

// String returns a string representation of a Bracket
func (bracket Bracket) String() string {
	if bracket.High == math.Inf(+1) {
		return fmt.Sprintf("%02.0f+", bracket.Low)
	}
	return fmt.Sprintf("%02.0f-%02.0f", bracket.Low, bracket.High)
}
