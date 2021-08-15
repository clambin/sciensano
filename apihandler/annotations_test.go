package apihandler_test

import (
	grafanaJson "github.com/clambin/grafana-json"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestHandler_Annotations(t *testing.T) {
	args := &grafanaJson.AnnotationRequestArgs{
		CommonQueryArgs: grafanaJson.CommonQueryArgs{
			Range: grafanaJson.QueryRequestRange{
				To: time.Now(),
			},
		},
	}

	annotations, err := apiHandler.Endpoints().Annotations("foo", "bar", args)
	assert.Nil(t, err)
	assert.Len(t, annotations, 3)
}
