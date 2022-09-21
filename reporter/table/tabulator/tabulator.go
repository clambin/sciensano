package tabulator

import (
	"time"
)

// Tabulator tabulates a set of entries in rows by timestamp and columns by label.  For performance reasons, data must
// be added sequentially.
type Tabulator struct {
	timestamps       []time.Time
	currentTimestamp time.Time
	currentRow       int
	columns          *Indexer[string]
	data             [][]float64
}

// New creates a new Tabulator
func New(columns ...string) *Tabulator {
	t := &Tabulator{
		columns:    MakeIndexer[string](),
		currentRow: -1,
	}
	t.RegisterColumn(columns...)
	return t
}

// Add adds a value for a specified timestamp and column to the table.  If there is already a value for that
// timestamp and column, the specified value is added to the existing value.
//
// Returns false if the column does not exist. Use RegisterColumn to add it first.
func (d *Tabulator) Add(timestamp time.Time, column string, value float64) bool {
	col, found := d.columns.GetIndex(column)
	if !found {
		return false
	}

	if !timestamp.Equal(d.currentTimestamp) {
		d.timestamps = append(d.timestamps, timestamp)
		d.data = append(d.data, make([]float64, d.columns.Count()))
		d.currentTimestamp = timestamp
		d.currentRow++
	}
	d.data[d.currentRow][col] += value
	return true
}

func (d *Tabulator) RegisterColumn(column ...string) {
	for _, c := range column {
		d.ensureColumnExists(c)
	}
}

func (d *Tabulator) ensureColumnExists(column string) {
	_, added := d.columns.Add(column)
	if !added {
		return
	}

	// new column. add data for the new column to each row
	for key, entry := range d.data {
		entry = append(entry, 0)
		d.data[key] = entry
	}
}

// Size returns the number of rows in the table.
func (d *Tabulator) Size() int {
	return len(d.data)
}

// GetTimestamps returns the (sorted) list of timestamps in the table.
func (d *Tabulator) GetTimestamps() (timestamps []time.Time) {
	return d.timestamps
}

// GetColumns returns the (sorted) list of column names.
func (d *Tabulator) GetColumns() (columns []string) {
	return d.columns.List()
}

// GetValues returns the value for the specified column for each timestamp in the table. The values are sorted by timestamp.
func (d *Tabulator) GetValues(column string) (values []float64, ok bool) {
	var index int
	index, ok = d.columns.GetIndex(column)

	if !ok {
		return
	}

	values = make([]float64, len(d.data))
	for i, row := range d.data {
		values[i] = row[index]
	}
	return
}
