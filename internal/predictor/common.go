package predictor

import log "github.com/sirupsen/logrus"

const (
	BatchSize       = 7
	ForecastBatches = 3
	HistoryBatches  = 12
	learnThreshold  = 0.98
	learnRetries    = 3
)

func forecastSamples(series1, series2 []float64, samples int, label string, responses chan float64) {
	p := New(BatchSize, 1000)

	totalSeries := make([][]float64, 0)

	totalSeries = append(totalSeries, series1)
	if len(series2) != 0 {
		totalSeries = append(totalSeries, series2)
	}

	var score float64
	retries := learnRetries
	for score < learnThreshold && retries > 0 {
		score = p.Learn(totalSeries)
		retries--
	}

	log.WithFields(log.Fields{
		"score":    score,
		"attempts": learnRetries - retries,
	}).Debugf("learned from %d %s samples", len(series1), label)

	responses <- score

	buffer := make([]float64, BatchSize)
	copy(buffer, series1[len(series1)-BatchSize:])

	for i := 0; i < samples; i++ {
		prediction, err := p.Predict(buffer)

		if err != nil {
			log.WithError(err).Warning("failed to predict vaccination evolution")
		}

		for _, value := range prediction {
			responses <- value
		}

		buffer = append(buffer[1:], prediction[0])
	}
}
