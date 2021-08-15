package apihandler_test

import (
	"context"
	"github.com/clambin/grafana-json"
	"github.com/clambin/sciensano/apihandler"
	"github.com/clambin/sciensano/demographics"
	mockDemographics "github.com/clambin/sciensano/demographics/mock"
	"github.com/clambin/sciensano/sciensano"
	mockSciensano "github.com/clambin/sciensano/sciensano/mock"
	mockVaccines "github.com/clambin/sciensano/vaccines/mock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var (
	sciensanoServer    mockSciensano.Handler
	sciensanoAPIServer *httptest.Server
	vaccinesServer     mockVaccines.Server
	vaccinesAPIServer  *httptest.Server
	demoServer         *mockDemographics.Server
	demoAPIServer      *demographics.Server
	apiHandler         *apihandler.Handler
)

func TestMain(m *testing.M) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sciensanoAPIServer = httptest.NewServer(http.HandlerFunc(sciensanoServer.Handle))
	defer sciensanoAPIServer.Close()

	vaccinesAPIServer = httptest.NewServer(http.HandlerFunc(vaccinesServer.Handler))
	defer vaccinesAPIServer.Close()

	demoServer = mockDemographics.New("../data/demographics.zip")
	defer demoServer.Close()
	demoAPIServer = demographics.New()
	demoAPIServer.URL = demoServer.URL()
	go func() {
		_ = demoAPIServer.Run(ctx, 24*time.Hour)
	}()

	client := sciensano.NewClient(time.Hour)
	client.SetURL(sciensanoAPIServer.URL)

	apiHandler = apihandler.Create(demoAPIServer)
	apiHandler.Sciensano = client
	apiHandler.Vaccines.URL = vaccinesAPIServer.URL

	m.Run()
}

var realTargets = map[string]bool{
	"tests":                    false,
	"vaccinations":             false,
	"vacc-age-partial":         false,
	"vacc-age-full":            false,
	"vacc-age-rate-partial":    false,
	"vacc-age-rate-full":       false,
	"vacc-region-partial":      false,
	"vacc-region-full":         false,
	"vacc-region-rate-partial": false,
	"vacc-region-rate-full":    false,
	"vaccination-lag":          false,
	"vaccines":                 false,
	"vaccines-stats":           false,
	"vaccines-time":            false,
}

func TestAPIHandler_Search(t *testing.T) {
	targets := apiHandler.Endpoints().Search()

	for _, target := range targets {
		_, ok := realTargets[target]
		if assert.True(t, ok, "unexpected target: "+target) {
			realTargets[target] = true
		}
	}

	for target, found := range realTargets {
		assert.True(t, found, "missing target:"+target)
	}

}

func TestAPIHandler_Invalid(t *testing.T) {
	endDate := time.Date(2021, 01, 06, 0, 0, 0, 0, time.UTC)
	request := &grafana_json.TableQueryArgs{
		CommonQueryArgs: grafana_json.CommonQueryArgs{
			Range: grafana_json.QueryRequestRange{To: endDate},
		},
	}

	var err error

	// Unknown target should return an error
	_, err = apiHandler.TableQuery(context.Background(), "invalid", request)
	assert.NotNil(t, err)

}

func BenchmarkHandler_QueryTable(b *testing.B) {
	endDate := time.Date(2021, 01, 06, 0, 0, 0, 0, time.UTC)
	request := &grafana_json.TableQueryArgs{
		CommonQueryArgs: grafana_json.CommonQueryArgs{
			Range: grafana_json.QueryRequestRange{
				To: endDate,
			},
		},
	}

	for target := range realTargets {
		for i := 0; i < 100; i++ {
			_, _ = apiHandler.Endpoints().TableQuery(context.Background(), target, request)
		}
	}
}
