package apihandler_test

import (
	"context"
	grafanaJson "github.com/clambin/grafana-json"
	"github.com/clambin/sciensano/vaccines"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stack := createStack(ctx)
	defer stack.Destroy()

	stack.vaccinesClient.
		On("GetBatches", mock.Anything).
		Return([]*vaccines.Batch{
			{
				Date:   vaccines.Time{Time: time.Now()},
				Amount: 100,
			},
		}, nil)

	annotations, err := stack.apiHandler.Endpoints().Annotations("foo", "bar", args)
	require.NoError(t, err)
	require.Len(t, annotations, 1)
	assert.Equal(t, "Amount: 100", annotations[0].Text)

	stack.vaccinesClient.AssertExpectations(t)
}
