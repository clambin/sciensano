package demographics_test

import (
	"context"
	"github.com/clambin/sciensano/demographics"
	"github.com/clambin/sciensano/demographics/mock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestServer_GetAgeGroupFigures(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	testServer := mock.New("")
	defer testServer.Close()

	server := demographics.New()
	server.URL = testServer.URL()

	go func() {
		_ = server.Run(ctx, 24*time.Hour)
	}()

	assert.Eventually(t, server.AvailableData, 1000*time.Millisecond, 10*time.Millisecond)

	figures := server.GetAgeGroupFigures()
	assert.Len(t, figures, 11)

}

func TestServer_GetRegionFigures(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	testServer := mock.New("")
	defer testServer.Close()

	server := demographics.New()
	server.URL = testServer.URL()

	go func() {
		_ = server.Run(ctx, 24*time.Hour)
	}()

	assert.Eventually(t, server.AvailableData, 1000*time.Millisecond, 10*time.Millisecond)

	figures := server.GetRegionFigures()
	assert.Len(t, figures, 4)
	assert.Contains(t, figures, "Flanders")
	assert.Contains(t, figures, "Wallonia")
	assert.Contains(t, figures, "Brussels")
	assert.Contains(t, figures, "Ostbelgien")
}
