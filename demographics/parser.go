package demographics

import (
	csv "github.com/rovaughn/fastcsv"
	log "github.com/sirupsen/logrus"
	"math"
	"strconv"
)

type populationRecord struct {
	Region []byte `csv:"TX_RGN_DESCR_NL"`
	Age    []byte `csv:"CD_AGE"`
	Count  []byte `csv:"MS_POPULATION\r"`
}

func groupPopulation(filename string) (byRegion map[string]int, byAge map[string]int, err error) {
	var record populationRecord
	var reader *csv.FileReader
	reader, err = csv.NewFileReader(filename, '|', &record)

	if err != nil {
		return
	}

	defer func() {
		_ = reader.Close()
	}()

	byRegion = make(map[string]int)
	byAge = make(map[string]int)

	var count int
	for reader.Scan() {
		if len(record.Count) > 0 && record.Count[len(record.Count)-1] == '\r' {
			record.Count = record.Count[:len(record.Count)-1]
		}
		count, err = strconv.Atoi(string(record.Count))
		region := string(record.Region)
		age := string(record.Age)

		if err != nil {
			return
		}

		byRegionCount, _ := byRegion[region]
		byRegion[region] = byRegionCount + count

		byAgeCount, _ := byAge[age]
		byAge[age] = byAgeCount + count
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
