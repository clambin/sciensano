package sciensano

import (
	"fmt"
)

const baseURL = "https://epistat.sciensano.be"

type Endpoint int

const (
	CasesEndpoint Endpoint = iota
	HospitalisationsEndpoint
	MortalitiesEndpoint
	TestResultsEndpoint
	VaccinationsEndpoint
)

var routes = map[Endpoint]string{
	CasesEndpoint:            "/Data/COVID19BE_CASES_AGESEX.json",
	HospitalisationsEndpoint: "/Data/COVID19BE_HOSP.json",
	MortalitiesEndpoint:      "/Data/COVID19BE_MORT.json",
	TestResultsEndpoint:      "/Data/COVID19BE_tests.json",
	VaccinationsEndpoint:     "/Data/COVID19BE_VACC.json",
}

func MustGetURL(base string, endpoint Endpoint) string {
	url, err := GetURL(base, endpoint)
	if err != nil {
		panic(err)
	}
	return url
}

func GetURL(base string, endpoint Endpoint) (string, error) {
	if base == "" {
		base = baseURL
	}
	route, ok := routes[endpoint]
	if !ok {
		return "", fmt.Errorf("invalid endpoint %d", endpoint)
	}
	return base + route, nil
}
