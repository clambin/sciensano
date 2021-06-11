package predictor

const (
	BatchSize       = 21
	ForecastBatches = 3
	learnRetries    = 1
)

type forecastFigures struct {
	figures []float64
	score   float64
	err     error
}
