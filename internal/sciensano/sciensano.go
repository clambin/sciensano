package sciensano

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

// Client queries different Sciensano APIs
type Client struct {
	apiClient http.Client

	VaccinationsCacheDuration time.Duration
	vaccinationsCacheExpiry   time.Time
	vaccinationsCache         []apiVaccinationsResponse
}

const baseURL = "https://epistat.sciensano.be/Data/"

// API for covid tests

type apiTestResponse struct {
	TimeStamp string `json:"DATE"`
	Province  string `json:"PROVINCE"`
	Region    string `json:"REGION"`
	Total     int    `json:"TESTS_ALL"`
	Positive  int    `json:"TESTS_ALL_POS"`
}

func (client *Client) getTests() (stats []apiTestResponse, err error) {
	req, _ := http.NewRequest("GET", baseURL+"COVID19BE_tests.json", nil)

	var resp *http.Response
	if resp, err = client.apiClient.Do(req); err == nil {
		defer resp.Body.Close()
		if resp.StatusCode == 200 {
			var (
				body []byte
			)
			if body, err = ioutil.ReadAll(resp.Body); err == nil {
				err = json.Unmarshal(body, &stats)
			}
		} else {
			err = errors.New(resp.Status)
		}
	}
	return
}

// API for vaccinations

type apiVaccinationsResponse struct {
	TimeStamp string `json:"DATE"`
	Region    string `json:"REGION"`
	AgeGroup  string `json:"AGEGROUP"`
	Gender    string `json:"SEX"`
	Dose      string `json:"DOSE"`
	Count     int    `json:"Count"`
}

func (client *Client) getVaccinations() ([]apiVaccinationsResponse, error) {
	var err error

	if client.vaccinationsCache == nil || time.Now().After(client.vaccinationsCacheExpiry) {

		req, _ := http.NewRequest("GET", baseURL+"COVID19BE_VACC.json", nil)

		var resp *http.Response
		var stats []apiVaccinationsResponse

		if resp, err = client.apiClient.Do(req); err == nil {
			defer resp.Body.Close()
			if resp.StatusCode == 200 {
				var (
					body []byte
				)
				if body, err = ioutil.ReadAll(resp.Body); err == nil {
					if err = json.Unmarshal(body, &stats); err == nil {
						client.vaccinationsCache = stats
						client.vaccinationsCacheExpiry = time.Now().Add(client.VaccinationsCacheDuration)
					}
				}
			} else {
				err = errors.New(resp.Status)
			}
		}
	}

	return client.vaccinationsCache, err
}
