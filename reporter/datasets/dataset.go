package datasets

import (
	"github.com/clambin/sciensano/apiclient"
	"sort"
	"time"
)

type Dataset struct {
	data        [][]float64
	timestamps  *Timestamps
	columns     []string
	columnIndex map[string]int
}

func New() *Dataset {
	return &Dataset{
		timestamps:  MakeTimestamps(),
		columnIndex: make(map[string]int),
	}
}

func NewFromAPIResponse(response []apiclient.APIResponse) (d *Dataset) {
	d = New()
	for _, entry := range response {
		attribs := entry.GetAttributeNames()
		values := entry.GetAttributeValues()

		for index, attrib := range attribs {
			d.Add(entry.GetTimestamp(), attrib, values[index])
		}
	}
	return
}

func NewGroupedFromAPIResponse(response []apiclient.APIResponse, groupField int) (d *Dataset) {
	d = New()
	for _, entry := range response {
		d.Add(entry.GetTimestamp(), entry.GetGroupFieldValue(groupField), entry.GetTotalValue())
	}
	return
}

func (d *Dataset) Add(timestamp time.Time, column string, value float64) {
	_, added := d.timestamps.Add(timestamp)
	if added {
		d.data = append(d.data, make([]float64, len(d.columns)))
	}
	d.ensureColumnExists(column)

	row := d.timestamps.GetIndex(timestamp)
	col := d.columnIndex[column]
	d.data[row][col] += value
}

func (d *Dataset) ensureColumnExists(column string) {
	// TODO: having a dedicated map would be faster
	_, ok := d.columnIndex[column]
	if ok == true {
		return
	}

	d.columnIndex[column] = len(d.columns)
	d.columns = append(d.columns, column)
	sort.Strings(d.columns)

	for key, entry := range d.data {
		entry = append(entry, 0)
		d.data[key] = entry
	}
}

func (d Dataset) Size() int {
	return d.timestamps.Count()
}

func (d *Dataset) AddColumn(column string, processor func(values map[string]float64) float64) {
	for index, row := range d.data {
		values := make(map[string]float64)
		for _, c := range d.columns {
			values[c] = row[d.columnIndex[c]]
		}

		newVal := processor(values)
		d.data[index] = append(row, newVal)
	}
	d.columnIndex[column] = len(d.columns)
	d.columns = append(d.columns, column)
	sort.Strings(d.columns)
}

func (d Dataset) GetTimestamps() (timestamps []time.Time) {
	timestamps = make([]time.Time, d.timestamps.Count())
	copy(timestamps, d.timestamps.List())
	return
}

func (d Dataset) GetColumns() (columns []string) {
	columns = make([]string, len(d.columns))
	copy(columns, d.columns)
	return
}

func (d Dataset) GetValues(column string) (values []float64, ok bool) {
	var index int
	index, ok = d.columnIndex[column]

	if ok == false {
		return
	}

	values = make([]float64, len(d.data))
	for i, timestamp := range d.timestamps.List() {
		rowIndex := d.timestamps.GetIndex(timestamp)
		values[i] = d.data[rowIndex][index]
	}
	return
}

func (d *Dataset) FilterByRange(from, to time.Time) {
	// make a list of all records to be removed, and the remaining timestamps
	timestamps := make([]time.Time, 0, d.timestamps.Count())
	var remove bool
	for _, timestamp := range d.timestamps.List() {
		if from.IsZero() == false && timestamp.Before(from) {
			remove = true
			continue
		} else if to.IsZero() == false && timestamp.After(to) {
			remove = true
			continue
		}
		timestamps = append(timestamps, timestamp)
	}

	// nothing to do here?
	if remove == false {
		return
	}

	// create a new data list from the timestamps we want to keep
	data := make([][]float64, len(timestamps))
	ts := MakeTimestamps()
	for index, timestamp := range timestamps {
		data[index] = d.data[d.timestamps.GetIndex(timestamp)]
		ts.Add(timestamp)
	}
	d.data = data
	d.timestamps = ts
}

func (d *Dataset) Accumulate() {
	accumulated := make([]float64, len(d.columns))

	for _, timestamp := range d.timestamps.List() {
		row := d.timestamps.GetIndex(timestamp)
		for index, value := range d.data[row] {
			accumulated[index] += value
		}
		copy(d.data[row], accumulated)
	}
}

func (d *Dataset) Copy() (clone *Dataset) {
	clone = &Dataset{
		data:        make([][]float64, len(d.data)),
		timestamps:  d.timestamps.Copy(),
		columns:     make([]string, len(d.columns)),
		columnIndex: make(map[string]int),
	}
	for index, row := range d.data {
		clone.data[index] = make([]float64, len(row))
		copy(clone.data[index], row)
	}
	copy(clone.columns, d.columns)
	for key, value := range d.columnIndex {
		clone.columnIndex[key] = value
	}
	return
}
