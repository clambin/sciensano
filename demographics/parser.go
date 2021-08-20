package demographics

import (
	"encoding/csv"
	log "github.com/sirupsen/logrus"
	"io"
	"math"
	"os"
	"strconv"
)

func groupPopulation(filename, mapField string) (output map[string]int, err error) {
	var f *os.File
	f, err = os.Open(filename)

	if err != nil {
		return
	}

	reader := csv.NewReader(f)
	reader.Comma = '|'

	var fields map[string]int
	output = make(map[string]int)
	first := true
	for err == nil {
		var record []string
		record, err = reader.Read()

		if err == nil {
			if first {
				fields = parseFields(record)
				first = false
			} else if len(record) != len(fields) {
				log.Warning("record mismatch. skipping entry")
			} else if len(record) > 0 {
				var count int
				count, err = strconv.Atoi(record[fields["MS_POPULATION"]])
				if err == nil {
					output[record[fields[mapField]]] += count
				}
			}
		}
	}

	if err == io.EOF {
		err = nil
	}

	return
}

func parseFields(record []string) (columns map[string]int) {
	columns = make(map[string]int)
	for index, field := range record {
		columns[field] = index
	}
	return
}

func groupPopulationByAge(input map[string]int, ranges []float64) (output map[Bracket]int) {
	output = make(map[Bracket]int)

	for _, bracket := range buildAgeBrackets(ranges) {
		output[bracket] = 0
	}

	for ageString, count := range input {
		age, err := strconv.ParseFloat(ageString, 64)

		if err != nil {
			log.WithError(err).Warning("could not parse age: " + ageString)
			continue
		}
		for bracket, total := range output {
			if age >= bracket.Low && age <= bracket.High {
				output[bracket] = total + count
			}
		}

	}
	return
}

func buildAgeBrackets(ranges []float64) (brackets []Bracket) {
	low := 0.0
	for _, high := range ranges {
		brackets = append(brackets, Bracket{Low: low, High: high - 1})
		low = high
	}
	brackets = append(brackets, Bracket{Low: low, High: math.Inf(+1)})
	return
}

func groupPopulationByRegion(input map[string]int) (output map[string]int) {
	translation := []struct {
		from string
		to   string
	}{
		{from: "Vlaams Gewest", to: "Flanders"},
		{from: "Waals Gewest", to: "Wallonia"},
		{from: "Brussels Hoofdstedelijk Gewest", to: "Brussels"},
	}

	output = make(map[string]int)
	for _, entry := range translation {
		output[entry.to] = input[entry.from]
	}

	return
}
