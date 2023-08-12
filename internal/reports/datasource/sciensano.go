package datasource

import (
	"github.com/clambin/go-common/taskmanager"
	"github.com/clambin/sciensano/internal/sciensano"
	"golang.org/x/exp/slog"
	"net/http"
	"time"
)

type SciensanoSources struct {
	taskmanager.Manager
	Cases            DataSource[sciensano.Cases]
	Hospitalisations DataSource[sciensano.Hospitalisations]
	Mortalities      DataSource[sciensano.Mortalities]
	TestResults      DataSource[sciensano.TestResults]
	Vaccinations     DataSource[sciensano.Vaccinations]
}

func NewSciensanoDatastore(url string, pollingInterval time.Duration, httpClient *http.Client, logger *slog.Logger) *SciensanoSources {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	store := SciensanoSources{
		Cases: DataSource[sciensano.Cases]{
			Fetcher:         &sciensano.Fetcher[sciensano.Cases]{Target: sciensano.MustGetURL(url, sciensano.CasesEndpoint), Client: httpClient},
			PollingInterval: pollingInterval,
			Logger:          logger.With("datasource", "cases"),
		},
		Hospitalisations: DataSource[sciensano.Hospitalisations]{
			Fetcher:         &sciensano.Fetcher[sciensano.Hospitalisations]{Target: sciensano.MustGetURL(url, sciensano.HospitalisationsEndpoint), Client: httpClient},
			PollingInterval: pollingInterval,
			Logger:          logger.With("datasource", "hospitalisations"),
		},
		Mortalities: DataSource[sciensano.Mortalities]{
			Fetcher:         &sciensano.Fetcher[sciensano.Mortalities]{Target: sciensano.MustGetURL(url, sciensano.MortalitiesEndpoint), Client: httpClient},
			PollingInterval: pollingInterval,
			Logger:          logger.With("datasource", "mortalities"),
		},
		TestResults: DataSource[sciensano.TestResults]{
			Fetcher:         &sciensano.Fetcher[sciensano.TestResults]{Target: sciensano.MustGetURL(url, sciensano.TestResultsEndpoint), Client: httpClient},
			PollingInterval: pollingInterval,
			Logger:          logger.With("datasource", "testResults"),
		},
		Vaccinations: DataSource[sciensano.Vaccinations]{
			Fetcher:         &sciensano.Fetcher[sciensano.Vaccinations]{Target: sciensano.MustGetURL(url, sciensano.VaccinationsEndpoint), Client: httpClient},
			PollingInterval: pollingInterval,
			Logger:          logger.With("datasource", "vaccinations"),
		},
	}
	_ = store.Add(&store.Cases)
	_ = store.Add(&store.Hospitalisations)
	_ = store.Add(&store.Mortalities)
	_ = store.Add(&store.TestResults)
	_ = store.Add(&store.Vaccinations)

	return &store
}
