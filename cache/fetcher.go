package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/clambin/httpclient"
	"github.com/go-http-utils/headers"
	"net/http"
	"time"
)

//go:generate mockery --name Fetcher
type Fetcher[T any] interface {
	GetLastModified(ctx context.Context) (time.Time, error)
	Fetch(ctx context.Context) (T, error)
	GetTarget() string
}

type fetcher[T any] struct {
	client httpclient.Caller
	target string
}

func (f *fetcher[T]) GetLastModified(ctx context.Context) (time.Time, error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodHead, f.target, nil)
	resp, err := f.client.Do(req)
	if err != nil {
		return time.Time{}, fmt.Errorf("%s: HEAD failed: %w", f.target, err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return time.Time{}, fmt.Errorf("%s: GET failed: %s", f.target, resp.Status)
	}
	return time.Parse(time.RFC1123, resp.Header.Get(headers.LastModified))
}

func (f *fetcher[T]) Fetch(ctx context.Context) (T, error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, f.target, nil)
	resp, err := f.client.Do(req)
	var records T
	if err != nil {
		return records, fmt.Errorf("%s: GET failed: %w", f.target, err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		return records, fmt.Errorf("%s: GET failed: %s", f.target, resp.Status)
	}
	err = json.NewDecoder(resp.Body).Decode(&records)
	return records, err
}

func (f *fetcher[T]) GetTarget() string {
	return f.target
}
