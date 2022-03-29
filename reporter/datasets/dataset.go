package datasets

import (
	"github.com/clambin/sciensano/apiclient"
	"time"
)

type Dataset struct {
	data       [][]float64
	timestamps *Indexer[time.Time]
	columns    *Indexer[string]
}

func New() *Dataset {
	return &Dataset{
		timestamps: MakeIndexer[time.Time](),
		columns:    MakeIndexer[string](),
	}
}

func NewFromAPIResponse(response []apiclient.APIResponse) (d *Dataset) {
	d = New()
	for _, entry := range response {
		ts := entry.GetTimestamp()
		attribs := entry.GetAttributeNames()
		values := entry.GetAttributeValues()

		for index, attrib := range attribs {
			value := values[index]
			d.Add(ts, attrib, value)
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
	d.ensureColumnExists(column)

	row, tsAdded := d.timestamps.Add(timestamp)
	if tsAdded {
		d.data = append(d.data, make([]float64, d.columns.Count()))
	}
	col, _ := d.columns.GetIndex(column)
	d.data[row][col] += value
}

func (d *Dataset) ensureColumnExists(column string) {
	_, ok := d.columns.Add(column)
	if ok == false {
		return
	}

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
		for _, c := range d.columns.List() {
			idx, _ := d.columns.GetIndex(c)
			values[c] = row[idx]
		}

		newVal := processor(values)
		d.data[index] = append(row, newVal)
	}
	d.columns.Add(column)
}

func (d Dataset) GetTimestamps() (timestamps []time.Time) {
	timestamps = make([]time.Time, d.timestamps.Count())
	copy(timestamps, d.timestamps.List())
	return
}

func (d Dataset) GetColumns() (columns []string) {
	columns = make([]string, d.columns.Count())
	copy(columns, d.columns.List())
	return
}

func (d Dataset) GetValues(column string) (values []float64, ok bool) {
	var index int
	index, ok = d.columns.GetIndex(column)

	if ok == false {
		return
	}

	values = make([]float64, len(d.data))
	for i, timestamp := range d.timestamps.List() {
		rowIndex, _ := d.timestamps.GetIndex(timestamp)
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
	ts := MakeIndexer[time.Time]()
	for index, timestamp := range timestamps {
		i, _ := d.timestamps.GetIndex(timestamp)
		data[index] = d.data[i]
		ts.Add(timestamp)
	}
	d.data = data
	d.timestamps = ts
}

func (d *Dataset) Accumulate() {
	accumulated := make([]float64, d.columns.Count())

	for _, timestamp := range d.timestamps.List() {
		row, _ := d.timestamps.GetIndex(timestamp)
		for index, value := range d.data[row] {
			accumulated[index] += value
		}
		copy(d.data[row], accumulated)
	}
}

func (d Dataset) Copy() (clone *Dataset) {
	clone = &Dataset{
		data:       make([][]float64, len(d.data)),
		timestamps: d.timestamps.Copy(),
		columns:    d.columns.Copy(),
	}
	for index, row := range d.data {
		clone.data[index] = make([]float64, len(row))
		copy(clone.data[index], row)
	}
	return
}
