package tabulator

import (
	"github.com/clambin/simplejson/v3/data"
	grafanaData "github.com/grafana/grafana-plugin-sdk-go/data"
	"time"
)

// Tabulator tabulates a set of entries in rows by timestamp and columns by label.  For performance reasons, data must
// be added sequentially.
type Tabulator struct {
	timestamps *Indexer[time.Time]
	columns    *Indexer[string]
	data       [][]float64
}

// New creates a new Tabulator
func New(columns ...string) *Tabulator {
	t := &Tabulator{
		timestamps: MakeIndexer[time.Time](),
		columns:    MakeIndexer[string](),
	}
	t.RegisterColumn(columns...)
	return t
}

// Add adds a value for a specified timestamp and column to the table.  If there is already a value for that
// timestamp and column, the specified value is added to the existing value.
//
// Returns false if the column does not exist. Use RegisterColumn to add it first.
func (t *Tabulator) Add(timestamp time.Time, column string, value float64) bool {
	col, found := t.columns.GetIndex(column)
	if !found {
		return false
	}

	row, added := t.timestamps.Add(timestamp)
	if added {
		t.data = append(t.data, make([]float64, t.columns.Count()))
	}

	t.data[row][col] += value
	return true
}

func (t *Tabulator) RegisterColumn(column ...string) {
	for _, c := range column {
		t.ensureColumnExists(c)
	}
}

func (t *Tabulator) ensureColumnExists(column string) {
	if _, added := t.columns.Add(column); added {
		// new column. add data for the new column to each row
		for key, entry := range t.data {
			entry = append(entry, 0)
			t.data[key] = entry
		}

	}
}

// Size returns the number of rows in the table.
func (t *Tabulator) Size() int {
	return len(t.data)
}

// GetTimestamps returns the (sorted) list of timestamps in the table.
func (t *Tabulator) GetTimestamps() []time.Time {
	return t.timestamps.List()
}

// GetColumns returns the (sorted) list of column names.
func (t *Tabulator) GetColumns() []string {
	return t.columns.List()
}

// GetValues returns the value for the specified column for each timestamp in the table. The values are sorted by timestamp.
func (t *Tabulator) GetValues(columnName string) (values []float64, ok bool) {
	var column int
	if column, ok = t.columns.GetIndex(columnName); ok {
		values = make([]float64, len(t.data))
		for index, timestamp := range t.timestamps.List() {
			row, _ := t.timestamps.GetIndex(timestamp)
			values[index] = t.data[row][column]
		}
	}
	return values, ok
}

// Accumulate increments the values in each column
func (t *Tabulator) Accumulate() {
	accumulated := make([]float64, len(t.GetColumns()))

	for _, timestamp := range t.GetTimestamps() {
		row, _ := t.timestamps.GetIndex(timestamp)
		for i, name := range t.columns.List() {
			column, _ := t.columns.GetIndex(name)
			t.data[row][column] += accumulated[i]
			accumulated[i] = t.data[row][column]
		}
	}
}

// Filter removes all rows that do not fall inside the specified time range. Is the specified time is zero, it will be ignored
func (t *Tabulator) Filter(from, to time.Time) {
	timestamps := MakeIndexer[time.Time]()
	d := make([][]float64, 0)

	for _, timestamp := range t.GetTimestamps() {
		if !from.IsZero() && timestamp.Before(from) {
			continue
		}
		if !to.IsZero() && timestamp.After(to) {
			continue
		}
		row, _ := t.timestamps.GetIndex(timestamp)
		timestamps.Add(timestamp)
		d = append(d, t.data[row])
	}

	t.timestamps = timestamps
	t.data = d
}

func (t *Tabulator) Copy() *Tabulator {
	t2 := New(t.GetColumns()...)

	timestamps := t.GetTimestamps()
	for _, column := range t.GetColumns() {
		values, _ := t.GetValues(column)
		for index, value := range values {
			t2.Add(timestamps[index], column, value)
		}
	}

	return t2
}

// MakeTable creates a simplejson Table from a Tabulator
func (t *Tabulator) MakeTable() *data.Table {
	fields := make(grafanaData.Fields, 0, t.Size())
	fields = append(fields, grafanaData.NewField("time", nil, t.GetTimestamps()))
	for _, col := range t.GetColumns() {
		values, _ := t.GetValues(col)
		fields = append(fields, grafanaData.NewField(col, nil, values))
	}
	return &data.Table{Frame: grafanaData.NewFrame("frame", fields...)}
}
