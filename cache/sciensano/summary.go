package sciensano

import (
	"fmt"
)

type SummaryColumn int

const (
	Total SummaryColumn = iota
	ByRegion
	ByProvince
	ByAgeGroup
	ByManufacturer
)

func (s SummaryColumn) String() string {
	switch s {
	case Total:
		return "Total"
	case ByRegion:
		return "ByRegion"
	case ByProvince:
		return "ByProvince"
	case ByAgeGroup:
		return "ByAgeGroup"
	case ByManufacturer:
		return "ByManufacturer"
	}

	panic(fmt.Sprintf("unknown summary column: %d", int(s)))
}
