package vaccines

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/clambin/httpclient"
	"github.com/clambin/sciensano/apiclient"
	"github.com/clambin/sciensano/apiclient/fetcher"
	"github.com/go-http-utils/headers"
	"io"
	"net/http"
	"sort"
	"time"
)

const (
	TypeBatches = iota
)

// Client calls the different Vaccines APIs
type Client struct {
	httpclient.Caller
	URL string
}

var _ fetcher.Fetcher = &Client{}

func (c Client) Fetch(ctx context.Context, dataType int) ([]apiclient.APIResponse, error) {
	resp, err := c.call(ctx, dataType, http.MethodGet)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	return c.parseResponse(dataType, resp.Body)
}

func (c Client) GetLastUpdated(ctx context.Context, dataType int) (time.Time, error) {
	var lastModified time.Time
	resp, err := c.call(ctx, dataType, http.MethodHead)
	if err == nil {
		_ = resp.Body.Close()
		if lastModified, err = time.Parse(time.RFC1123, resp.Header.Get(headers.LastModified)); err != nil {
			err = fmt.Errorf("invalid timestamp %s: %w", headers.LastModified, err)
		}
	}
	return lastModified, err
}

func (c Client) DataTypes() map[int]string {
	return map[int]string{
		TypeBatches: "Batches",
	}
}

var endpoints = map[int]string{
	TypeBatches: "delivered.json",
}

func (c Client) getURL(dataType int) (string, error) {
	endpoint, found := endpoints[dataType]
	if !found {
		return "", fmt.Errorf("invalid data type: %d", dataType)
	}
	target := "https://covid-vaccinatie.be"
	if c.URL != "" {
		target = c.URL
	}
	return target + "/api/v1/" + endpoint, nil
}

func (c Client) call(ctx context.Context, dataType int, method string) (resp *http.Response, err error) {
	var target string
	if target, err = c.getURL(dataType); err != nil {
		return
	}

	req, _ := http.NewRequestWithContext(ctx, method, target, nil)

	if resp, err = c.Caller.Do(req); err == nil && resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("call failed: %s", resp.Status)
	}
	return
}

func (c Client) parseResponse(dataType int, body io.Reader) (results []apiclient.APIResponse, err error) {
	switch dataType {
	case TypeBatches:
		results, err = c.parseBatches(body)
	}
	return
}

type apiBatchesResponse struct {
	Result struct {
		Delivered []*APIBatchResponse `json:"delivered"`
	} `json:"result"`
}

func (c Client) parseBatches(body io.Reader) (results []apiclient.APIResponse, err error) {
	var response apiBatchesResponse
	if err = json.NewDecoder(body).Decode(&response); err != nil {
		err = fmt.Errorf("decode failed: %w", err)
	}

	results = make([]apiclient.APIResponse, len(response.Result.Delivered))
	var mustSort bool
	var timestamp time.Time
	for idx, entry := range response.Result.Delivered {
		results[idx] = entry
		if !mustSort && entry.GetTimestamp().Before(timestamp) {
			mustSort = true
		}
		timestamp = entry.GetTimestamp()
	}

	if mustSort {
		sort.Slice(results, func(i, j int) bool { return results[i].GetTimestamp().Before(results[j].GetTimestamp()) })
	}
	return
}
