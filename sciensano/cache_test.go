package sciensano_test

import (
	"github.com/clambin/sciensano/sciensano"
	"github.com/clambin/sciensano/sciensano/datasets"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type testData struct {
	Value int
}

func (td *testData) Copy() datasets.Copyable {
	return &testData{Value: td.Value}
}

func TestCache(t *testing.T) {
	c := sciensano.NewCache(15 * time.Minute)

	for i := 0; i < 10; i++ {
		e := c.Load("foo")
		e.Once.Do(func() {
			if e.Data == nil {
				e.Data = &datasets.Dataset{
					Timestamps: make([]time.Time, 0),
					Groups:     nil,
				}
			}
			e.Data.Timestamps = append(e.Data.Timestamps, time.Now())
		})
		c.Save("foo", e)
	}

	e := c.Load("foo")
	assert.Len(t, e.Data.Timestamps, 1)
}
