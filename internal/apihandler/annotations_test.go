package apihandler_test

import (
	grafana_json "github.com/clambin/grafana-json"
	"github.com/clambin/sciensano/internal/apihandler"
	"github.com/clambin/sciensano/internal/vaccines/mock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestHandler_Annotations(t *testing.T) {
	handler, _ := apihandler.Create()
	handler.Vaccines.HTTPClient = mock.GetServer()

	args := &grafana_json.AnnotationRequestArgs{
		CommonQueryArgs: grafana_json.CommonQueryArgs{
			Range: grafana_json.QueryRequestRange{
				To: time.Now(),
			},
		},
	}

	annotations, err := handler.Endpoints().Annotations("foo", "bar", args)
	assert.Nil(t, err)
	assert.Len(t, annotations, 3)
}
