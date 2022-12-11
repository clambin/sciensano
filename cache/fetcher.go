package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/clambin/sciensano/cache/sciensano"
	"github.com/go-http-utils/headers"
	"github.com/mailru/easyjson"
	"io"
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
	client *http.Client
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
	return unmarshal[T](resp.Body)
}

func (f *fetcher[T]) GetTarget() string {
	return f.target
}

func unmarshal[T any](r io.Reader) (v T, err error) {
	switch interface{}(v).(type) {
	case sciensano.Vaccinations:
		var v2 sciensano.Vaccinations
		if err = easyjson.UnmarshalFromReader(r, &v2); err == nil {
			v = interface{}(v2).(T)
		}
	case sciensano.Cases:
		var v2 sciensano.Cases
		if err = easyjson.UnmarshalFromReader(r, &v2); err == nil {
			v = interface{}(v2).(T)
		}
	default:
		err = json.NewDecoder(r).Decode(&v)
	}
	return v, err
}
