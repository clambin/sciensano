package fetcher

import (
	"context"
	"github.com/clambin/sciensano/apiclient"
	"time"
)

// Fetcher retrieves data from an API
//go:generate mockery --name Fetcher
type Fetcher interface {
	Fetch(ctx context.Context, dataType int) ([]apiclient.APIResponse, error)
	GetLastUpdates(ctx context.Context, dataType int) (time.Time, error)
	DataTypes() map[int]string
}
