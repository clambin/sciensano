package vaccines

import (
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"sort"
	"sync"
	"time"
)

type Server struct {
	HTTPClient    *http.Client
	cacheDuration time.Duration
	cache         []Batch
	expiry        time.Time
	lock          sync.Mutex
}

type Batch struct {
	Date         Time
	Manufacturer string
	Amount       int64
}

func New() (server *Server) {
	server = &Server{
		HTTPClient:    &http.Client{},
		cacheDuration: 6 * time.Hour,
	}
	return
}

func (server *Server) GetBatches() (batches []Batch, err error) {
	server.lock.Lock()
	defer server.lock.Unlock()

	return server.getBatches()
}

type Time time.Time

func (date *Time) UnmarshalJSON(b []byte) (err error) {
	var timestamp time.Time
	if timestamp, err = time.Parse(`"2006-01-02"`, string(b)); err == nil {
		*date = Time(timestamp)
	}
	return
}

func (server *Server) getBatches() (batches []Batch, err error) {
	if server.cache == nil || time.Now().After(server.expiry) {
		var resp *http.Response
		var stats struct {
			Result struct {
				Delivered []Batch `json:"delivered"`
			} `json:"result"`
		}

		if resp, err = server.HTTPClient.Get("https://covid-vaccinatie.be/api/v1/delivered.json"); err == nil {
			defer resp.Body.Close()
			if resp.StatusCode == 200 {
				var body []byte

				if body, err = io.ReadAll(resp.Body); err == nil {
					if err = json.Unmarshal(body, &stats); err == nil {
						batches = stats.Result.Delivered
					} else {
						log.Error(err)
					}
				}
			} else {
				err = errors.New(resp.Status)
			}
		}

		if err == nil {
			sort.Slice(batches, func(i, j int) bool { return time.Time(batches[i].Date).Before(time.Time(batches[j].Date)) })
			server.cache = batches
			server.expiry = time.Now().Add(server.cacheDuration)
		}
	}
	batches = server.cache
	return
}

func AccumulateBatches(batches []Batch) (accumulated []Batch) {
	var total int64
	for _, batch := range batches {
		total += batch.Amount
		accumulated = append(accumulated, Batch{
			Date:   batch.Date,
			Amount: total,
		})
	}
	return
}
