package reports_test

import (
	"github.com/clambin/go-common/tabulator"
	"strconv"
	"time"
)

func createBigDataSet() (*tabulator.Tabulator, error) {
	d := tabulator.New("0", "1", "2", "3", "4", "5", "6", "7", "8", "9")
	timestamp := time.Now()
	for r := 0; r < 500; r++ {
		for c := 0; c < 10; c++ {
			d.Add(timestamp, strconv.Itoa(c), float64(r))
		}
		timestamp = timestamp.Add(24 * time.Hour)
	}
	return d, nil
}

func createSimpleDataSet() (*tabulator.Tabulator, error) {
	d := tabulator.New("A", "B", "C")
	d.Add(time.Date(2023, time.July, 31, 0, 0, 0, 0, time.UTC), "C", 3)
	d.Add(time.Date(2023, time.July, 30, 0, 0, 0, 0, time.UTC), "C", 2)
	d.Add(time.Date(2023, time.July, 29, 0, 0, 0, 0, time.UTC), "C", 1)
	return d, nil
}
