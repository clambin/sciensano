package population

import (
	"fmt"
	csv "github.com/rovaughn/fastcsv"
	"strconv"
)

type populationRecord struct {
	Region []byte `csv:"TX_RGN_DESCR_NL"`
	Age    []byte `csv:"CD_AGE"`
	Count  []byte `csv:"MS_POPULATION\r"`
}

func groupPopulation(filename string) (map[string]int, map[int]int, error) {
	var record populationRecord
	reader, err := csv.NewFileReader(filename, '|', &record)
	if err != nil {
		return nil, nil, err
	}

	defer func() {
		_ = reader.Close()
	}()

	byRegion := make(map[string]int)
	byAge := make(map[int]int)

	var line int
	for reader.Scan() {
		if len(record.Count) > 0 && record.Count[len(record.Count)-1] == '\r' {
			record.Count = record.Count[:len(record.Count)-1]
		}

		var count int
		count, err = strconv.Atoi(string(record.Count))
		if err != nil {
			return nil, nil, fmt.Errorf("invalid number for Count on line %d: %w", line, err)
		}

		var age int
		age, err = strconv.Atoi(string(record.Age))
		if err != nil {
			return nil, nil, fmt.Errorf("invalid number for Age on line %d: %w", line, err)
		}

		region := translateRegion(string(record.Region))

		byRegion[region] += count
		byAge[age] += count
	}

	return byRegion, byAge, nil
}

var regionTranslationTable = map[string]string{
	"Vlaams Gewest":                  "Flanders",
	"Waals Gewest":                   "Wallonia",
	"Brussels Hoofdstedelijk Gewest": "Brussels",
}

func translateRegion(input string) string {
	output, ok := regionTranslationTable[input]
	if !ok {
		output = input
	}
	return output
}
