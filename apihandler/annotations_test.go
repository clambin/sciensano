package apihandler_test

import (
	grafanaJson "github.com/clambin/grafana-json"
	"github.com/clambin/sciensano/apihandler"
	"github.com/clambin/sciensano/vaccines/mock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestHandler_Annotations(t *testing.T) {
	handler, _ := apihandler.Create()
	handler.Vaccines.HTTPClient = mock.GetServer()

	args := &grafanaJson.AnnotationRequestArgs{
		CommonQueryArgs: grafanaJson.CommonQueryArgs{
			Range: grafanaJson.QueryRequestRange{
				To: time.Now(),
			},
		},
	}

	annotations, err := handler.Endpoints().Annotations("foo", "bar", args)
	assert.Nil(t, err)
	assert.Len(t, annotations, 3)
}
