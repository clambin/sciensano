package vaccinations

import (
	grafanajson "github.com/clambin/grafana-json"
	"github.com/clambin/sciensano/sciensano"
	"github.com/clambin/sciensano/sciensano/datasets"
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
		Groups: []datasets.GroupedDatasetEntry{
			{
				Name: "",
				Values: []datasets.Copyable{
					&sciensano.VaccinationsEntry{Partial: 0, Full: 0},
					&sciensano.VaccinationsEntry{Partial: 1, Full: 0},
					&sciensano.VaccinationsEntry{Partial: 2, Full: 1},
					&sciensano.VaccinationsEntry{Partial: 3, Full: 2},
					&sciensano.VaccinationsEntry{Partial: 4, Full: 3},
					&sciensano.VaccinationsEntry{Partial: 5, Full: 4},
					&sciensano.VaccinationsEntry{Partial: 6, Full: 5},
				}},
		},
	}

	_, lag := buildLag(vaccinations)

	assert.Equal(t, grafanajson.TableQueryResponseNumberColumn{1.0, 1.0, 1.0, 1.0, 1.0}, lag)

	vaccinations.Groups[0].Values = []datasets.Copyable{
		&sciensano.VaccinationsEntry{Partial: 1, Full: 1}, // 0
		&sciensano.VaccinationsEntry{Partial: 1, Full: 1}, // -
		&sciensano.VaccinationsEntry{Partial: 2, Full: 1}, // -
		&sciensano.VaccinationsEntry{Partial: 3, Full: 2}, // 1
		&sciensano.VaccinationsEntry{Partial: 4, Full: 3}, // -
		&sciensano.VaccinationsEntry{Partial: 4, Full: 4}, // -
		&sciensano.VaccinationsEntry{Partial: 6, Full: 5}, // 0
	}

	_, lag = buildLag(vaccinations)

	assert.Equal(t, grafanajson.TableQueryResponseNumberColumn{0.0, 1.0, 1.0, 1.0, 0.0}, lag)
}
