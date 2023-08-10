package sciensano

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mailru/easyjson"
	"io"
	"net/http"
	"time"
)

type Fetcher[T any] struct {
	Target string
	Client *http.Client
}

func (f *Fetcher[T]) GetLastModified(ctx context.Context) (time.Time, error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodHead, f.Target, nil)
	resp, err := f.Client.Do(req)
	if err != nil {
		return time.Time{}, fmt.Errorf("%s: HEAD failed: %w", f.Target, err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return time.Time{}, fmt.Errorf("%s: GET failed: %s", f.Target, resp.Status)
	}
	return time.Parse(time.RFC1123, resp.Header.Get("Last-Modified"))
}

func (f *Fetcher[T]) Fetch(ctx context.Context) (T, error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, f.Target, nil)
	resp, err := f.Client.Do(req)
	var records T
	if err != nil {
		return records, fmt.Errorf("%s: GET failed: %w", f.Target, err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		return records, fmt.Errorf("%s: GET failed: %s", f.Target, resp.Status)
	}
	return unmarshal[T](resp.Body)
}

func (f *Fetcher[T]) GetTarget() string {
	return f.Target
}

func unmarshal[T any](r io.Reader) (v T, err error) {
	switch interface{}(v).(type) {
	case Vaccinations:
		var v2 Vaccinations
		if err = easyjson.UnmarshalFromReader(r, &v2); err == nil {
			v = any(v2).(T)
		}
	case Cases:
		var v2 Cases
		if err = easyjson.UnmarshalFromReader(r, &v2); err == nil {
			v = any(v2).(T)
		}
	default:
		err = json.NewDecoder(r).Decode(&v)
	}
	return v, err
}
