package apihandler_test

import (
	"context"
	"github.com/clambin/grafana-json"
	"github.com/clambin/sciensano/apihandler"
	mockDemographics "github.com/clambin/sciensano/demographics/mocks"
	"github.com/clambin/sciensano/sciensano"
	"github.com/clambin/sciensano/sciensano/apiclient"
	mockSciensano "github.com/clambin/sciensano/sciensano/apiclient/mocks"
	"github.com/clambin/sciensano/vaccines"
	mockVaccines "github.com/clambin/sciensano/vaccines/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

type Stack struct {
	sciensanoClient mockSciensano.APIClient
	vaccinesClient  mockVaccines.APIClient
	demoClient      mockDemographics.Demographics
	apiHandler      *apihandler.Handler
}

func createStack(_ context.Context) *Stack {
	stack := &Stack{
		apiHandler: apihandler.Create(),
	}
	stack.apiHandler.Sciensano = &sciensano.Client{APIClient: &stack.sciensanoClient}
	stack.apiHandler.Vaccines = &stack.vaccinesClient
	stack.apiHandler.Demographics = &stack.demoClient

	return stack
}

func (stack *Stack) Destroy() {
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stack := createStack(ctx)
	defer stack.Destroy()

	targets := stack.apiHandler.Endpoints().Search()

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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stack := createStack(ctx)
	defer stack.Destroy()

	// Unknown target should return an error
	_, err := stack.apiHandler.TableQuery(context.Background(), "invalid", request)
	require.Error(t, err)
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stack := createStack(ctx)
	defer stack.Destroy()

	stack.sciensanoClient.
		On("GetVaccinations", mock.Anything).
		Return([]*apiclient.APIVaccinationsResponse{}, nil)
	stack.sciensanoClient.
		On("GetTestResults", mock.Anything).
		Return([]*apiclient.APITestResultsResponse{}, nil)
	stack.vaccinesClient.
		On("GetBatches", mock.Anything).
		Return([]*vaccines.Batch{}, nil)

	for target := range realTargets {
		for i := 0; i < 100; i++ {
			_, _ = stack.apiHandler.Endpoints().TableQuery(context.Background(), target, request)
		}
	}

	stack.sciensanoClient.AssertExpectations(b)
}
