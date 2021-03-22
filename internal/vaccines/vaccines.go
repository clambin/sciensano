package vaccines

import (
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"sort"
	"time"
)

type Handler struct {
	CacheTime time.Duration
	Request   chan ResponseChannel

	HTTPClient *http.Client
	cache      []Batch
	expiry     time.Time
}

type ResponseChannel chan []Batch

type Batch struct {
	Date         Time
	Manufacturer string
	Amount       int64
}

func Create() (handler *Handler) {
	handler = &Handler{
		CacheTime:  6 * time.Hour,
		Request:    make(chan ResponseChannel),
		HTTPClient: &http.Client{},
	}
	go handler.Run()
	return
}

func (handler *Handler) Run() {
	for {
		select {
		case channel := <-handler.Request:
			// TODO: is this a race condition? does channel get its own copy of the batches?
			batches, _ := handler.getBatches()
			channel <- batches
		}
	}
}

type apiResult struct {
	Result struct {
		Delivered []Batch `json:"delivered"`
	} `json:"result"`
}

type Time time.Time

func (date *Time) UnmarshalJSON(b []byte) (err error) {
	var timestamp time.Time
	if timestamp, err = time.Parse(`"2006-01-02"`, string(b)); err == nil {
		*date = Time(timestamp)
	}
	return
}

func (handler *Handler) getBatches() (batches []Batch, err error) {
	if handler.cache == nil || time.Now().After(handler.expiry) {
		var resp *http.Response
		var stats apiResult

		if resp, err = handler.HTTPClient.Get("https://covid-vaccinatie.be/api/v1/delivered.json"); err == nil {
			defer resp.Body.Close()
			if resp.StatusCode == 200 {
				var body []byte

				if body, err = ioutil.ReadAll(resp.Body); err == nil {
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
			handler.cache = batches
			handler.expiry = time.Now().Add(handler.CacheTime)
		}

	}
	batches = handler.cache
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
