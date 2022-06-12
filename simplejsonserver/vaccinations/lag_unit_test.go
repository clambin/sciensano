package vaccinations

import (
	"github.com/clambin/simplejson/v3/data"
	grafanaData "github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestBuildLag(t *testing.T) {
	vaccinations := data.Table{Frame: grafanaData.NewFrame("frame",
		grafanaData.NewField("time", nil, []time.Time{
			time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC),
			time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC),
			time.Date(2021, 1, 4, 0, 0, 0, 0, time.UTC),
			time.Date(2021, 1, 5, 0, 0, 0, 0, time.UTC),
			time.Date(2021, 1, 6, 0, 0, 0, 0, time.UTC),
			time.Date(2021, 1, 7, 0, 0, 0, 0, time.UTC),
		}),
		grafanaData.NewField("partial", nil, []float64{0, 1, 2, 3, 4, 5, 6}),
		grafanaData.NewField("full", nil, []float64{0, 0, 1, 2, 3, 4, 5}),
	)}

	lag := buildLag(&vaccinations)
	values, ok := lag.GetFloatValues("lag")
	require.True(t, ok)
	assert.Equal(t, []float64{1.0, 1.0, 1.0, 1.0, 1.0}, values)

	vaccinations = data.Table{Frame: grafanaData.NewFrame("frame",
		grafanaData.NewField("time", nil, []time.Time{
			time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC),
			time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC),
			time.Date(2021, 1, 4, 0, 0, 0, 0, time.UTC),
			time.Date(2021, 1, 5, 0, 0, 0, 0, time.UTC),
			time.Date(2021, 1, 6, 0, 0, 0, 0, time.UTC),
			time.Date(2021, 1, 7, 0, 0, 0, 0, time.UTC),
		}),
		grafanaData.NewField("partial", nil, []float64{1, 1, 2, 3, 4, 5, 6}),
		grafanaData.NewField("full", nil, []float64{1, 1, 1, 2, 3, 4, 5}),
	)}

	lag = buildLag(&vaccinations)
	values, ok = lag.GetFloatValues("lag")
	require.True(t, ok)
	assert.Equal(t, []float64{0.0, 1.0, 1.0, 1.0, 1.0}, values)
}
