package booster

import (
	"github.com/clambin/sciensano/reporter"
	"github.com/clambin/simplejson/v3/data"
	grafanaData "github.com/grafana/grafana-plugin-sdk-go/data"
)

type Handler struct {
	Reporter *reporter.Client
}

func (h *Handler) Fetch() (results *data.Table, err error) {
	results, err = h.Reporter.Vaccinations.Get()
	if err != nil {
		return
	}

	// FIXME: would be easier if data.Table had a "DeleteColumn" method
	timestamps := results.GetTimestamps()
	booster, _ := results.GetFloatValues("booster")
	booster2, _ := results.GetFloatValues("booster2")
	booster3, _ := results.GetFloatValues("booster3")

	var fields []*grafanaData.Field
	fields = append(fields, grafanaData.NewField("timestamps", grafanaData.Labels{}, timestamps))
	fields = append(fields, grafanaData.NewField("booster", grafanaData.Labels{}, booster))
	fields = append(fields, grafanaData.NewField("booster2", grafanaData.Labels{}, booster2))
	fields = append(fields, grafanaData.NewField("booster3", grafanaData.Labels{}, booster3))

	return &data.Table{Frame: grafanaData.NewFrame("test", fields...)}, nil
}
