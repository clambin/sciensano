package vaccines_test

import (
	"context"
	"github.com/clambin/sciensano/vaccines"
	"github.com/clambin/sciensano/vaccines/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestCache_GetBatches(t *testing.T) {
	client := &mocks.APIClient{}
	cache := vaccines.Cache{
		APIClient: client,
		Retention: time.Hour,
	}
	ctx := context.Background()

	// Cache should only call the client once.
	client.
		On("GetBatches", mock.Anything).
		Return([]*vaccines.Batch{{
			Date:   vaccines.Time{Time: time.Now()},
			Amount: 100,
		}}, nil).
		Once()

	results, err := cache.GetBatches(ctx)
	require.NoError(t, err)
	require.Len(t, results, 1)

	results, err = cache.GetBatches(ctx)
	require.NoError(t, err)
	require.Len(t, results, 1)

	mock.AssertExpectationsForObjects(t, client)
}
