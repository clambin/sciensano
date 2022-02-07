package demographics

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

func groupPopulation(filename string) (byRegion map[string]int, byAge map[int]int, err error) {
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
	byAge = make(map[int]int)

	var line int
	for reader.Scan() {
		if len(record.Count) > 0 && record.Count[len(record.Count)-1] == '\r' {
			record.Count = record.Count[:len(record.Count)-1]
		}

		var count int
		count, err = strconv.Atoi(string(record.Count))
		if err != nil {
			err = fmt.Errorf("invalid number for Count on line %d: %w", line, err)
			return
		}

		region := translateRegion(string(record.Region))

		var age int
		age, err = strconv.Atoi(string(record.Age))
		if err != nil {
			err = fmt.Errorf("invalid number for Age on line %d: %w", line, err)
			return
		}

		byRegionCount, _ := byRegion[region]
		byRegion[region] = byRegionCount + count

		byAgeCount, _ := byAge[age]
		byAge[age] = byAgeCount + count
	}

	return
}

func translateRegion(input string) (output string) {
	translation := map[string]string{
		"Vlaams Gewest":                  "Flanders",
		"Waals Gewest":                   "Wallonia",
		"Brussels Hoofdstedelijk Gewest": "Brussels",
	}

	var ok bool
	output, ok = translation[input]
	if ok == false {
		output = input
	}
	return
}
