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
	URL           string
	cacheDuration time.Duration
	cache         []Batch
	expiry        time.Time
	lock          sync.Mutex
}

type Batch struct {
	Date Time
	// Manufacturer string
	Amount int
}

func New() (server *Server) {
	server = &Server{
		HTTPClient:    &http.Client{},
		cacheDuration: 1 * time.Hour,
		URL:           "https://covid-vaccinatie.be",
	}
	return
}

type Time time.Time

func (date *Time) UnmarshalJSON(b []byte) (err error) {
	var timestamp time.Time
	if timestamp, err = time.Parse(`"2006-01-02"`, string(b)); err == nil {
		*date = Time(timestamp)
	}
	return
}

func (server *Server) GetBatches() (batches []Batch, err error) {
	server.lock.Lock()
	defer server.lock.Unlock()

	if server.cache == nil || time.Now().After(server.expiry) {
		var resp *http.Response
		resp, err = server.HTTPClient.Get(server.URL + "/api/v1/delivered.json")

		if err == nil {
			if resp.StatusCode == http.StatusOK {
				var body []byte
				body, err = io.ReadAll(resp.Body)

				if err == nil {
					var stats struct {
						Result struct {
							Delivered []Batch `json:"delivered"`
						} `json:"result"`
					}
					err = json.Unmarshal(body, &stats)

					if err == nil {
						batches = stats.Result.Delivered
					} else {
						log.Error(err)
					}
				}
			} else {
				err = errors.New(resp.Status)
			}
			_ = resp.Body.Close()
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
	var total int
	for _, batch := range batches {
		total += batch.Amount
		accumulated = append(accumulated, Batch{
			Date:   batch.Date,
			Amount: total,
		})
	}
	return
}
