package vaccinations

import (
	"github.com/clambin/sciensano/reporter/datasets"
	"github.com/clambin/simplejson"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestVaccinationLag(t *testing.T) {
	vaccinations := &datasets.Dataset{
		Timestamps: []time.Time{
			time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC),
			time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC),
			time.Date(2021, 1, 4, 0, 0, 0, 0, time.UTC),
			time.Date(2021, 1, 5, 0, 0, 0, 0, time.UTC),
			time.Date(2021, 1, 6, 0, 0, 0, 0, time.UTC),
			time.Date(2021, 1, 7, 0, 0, 0, 0, time.UTC),
		},
		Groups: []datasets.DatasetGroup{
			{Name: "partial", Values: []float64{0, 1, 2, 3, 4, 5, 6}},
			{Name: "full", Values: []float64{0, 0, 1, 2, 3, 4, 5}},
		},
	}
	_, lag := buildLag(vaccinations)
	assert.Equal(t, simplejson.TableQueryResponseNumberColumn{1.0, 1.0, 1.0, 1.0, 1.0}, lag)

	vaccinations.Groups = []datasets.DatasetGroup{
		{Name: "partial", Values: []float64{1, 1, 2, 3, 4, 4, 6}},
		{Name: "full", Values: []float64{1, 1, 1, 2, 3, 4, 5}},
	}
	_, lag = buildLag(vaccinations)
	assert.Equal(t, simplejson.TableQueryResponseNumberColumn{0.0, 1.0, 1.0, 1.0, 0.0}, lag)
}
