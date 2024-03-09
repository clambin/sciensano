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
	ByVaccinationType
	ByCategory
)

var SummaryColumnNames map[string]SummaryColumn

func init() {
	SummaryColumnNames = make(map[string]SummaryColumn)

	for i := range ByCategory + 1 {
		SummaryColumnNames[i.String()] = i
	}
}

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
	case ByVaccinationType:
		return "ByVaccinationType"
	case ByCategory:
		return "ByCategory"
	}

	panic(fmt.Sprintf("unknown summary column: %d", int(s)))
}
