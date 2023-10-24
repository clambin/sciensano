package sciensano

import "fmt"

// DoseType is the type of vaccination: partial, full, single dose, booster, etc.
type DoseType int

const (
	Partial DoseType = iota
	Full
	SingleDose
	Booster
	Booster2
	Booster3
	Booster4
)

var doseTypeStrings = map[DoseType]string{
	Partial:    "Partial",
	Full:       "Full",
	SingleDose: "SingleDose",
	Booster:    "Booster",
	Booster2:   "Booster 2",
	Booster3:   "Booster 3",
	Booster4:   "Booster 4",
}

var DoseTypeNames map[string]DoseType

func init() {
	DoseTypeNames = make(map[string]DoseType)

	for i := Partial; i <= Booster4; i++ {
		DoseTypeNames[i.String()] = i
	}
}

func (d DoseType) String() string {
	value, ok := doseTypeStrings[d]
	if !ok {
		value = "(unknown)"
	}
	return value
}

func (d *DoseType) UnmarshalJSON(body []byte) (err error) {
	switch string(body) {
	case `"A"`:
		*d = Partial
	case `"B"`:
		*d = Full
	case `"C"`:
		*d = SingleDose
	case `"E"`:
		*d = Booster
	case `"E2"`:
		*d = Booster2
	case `"E3"`:
		*d = Booster3
	case `"E4+"`:
		*d = Booster4
	default:
		err = fmt.Errorf("invalid Dose: %s", string(body))
	}
	return
}

func (d DoseType) MarshalJSON() (body []byte, err error) {
	switch d {
	case Partial:
		body = []byte(`"A"`)
	case Full:
		body = []byte(`"B"`)
	case SingleDose:
		body = []byte(`"C"`)
	case Booster:
		body = []byte(`"E"`)
	case Booster2:
		body = []byte(`"E2"`)
	case Booster3:
		body = []byte(`"E3"`)
	case Booster4:
		body = []byte(`"E4+"`)
	default:
		err = fmt.Errorf("invalid Dose: %d", d)
	}
	return
}
