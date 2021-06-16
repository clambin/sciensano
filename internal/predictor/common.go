package predictor

import log "github.com/sirupsen/logrus"

const (
	BatchSize       = 7
	ForecastBatches = 3
	HistoryBatches  = 12
	learnThreshold  = 0.945
	learnRetries    = 3
)

func forecastSamples(series1, series2 []float64, forecastCount int, label string, responses chan []float64) {
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

	responses <- []float64{score}

	buffer := make([]float64, BatchSize)
	copy(buffer, series1[:BatchSize])

	for i := BatchSize; i < len(series1)+forecastCount; i++ {
		prediction, err := p.Predict(buffer)

		if err != nil {
			log.WithError(err).Warning("failed to predict vaccination evolution")
		}

		responses <- prediction

		if i < len(series1) {
			buffer = append(buffer[1:], series1[i])
		} else {
			buffer = append(buffer[1:], prediction[0])
		}
	}

	close(responses)
}

func consolidateSamples(output chan []float64, input []chan []float64, processor func([][]float64) []float64) {
	for values := range input[0] {
		allValues := [][]float64{values}

		for i := 1; i < len(input); i++ {
			allValues = append(allValues, <-input[i])
		}

		output <- processor(allValues)
	}

	close(output)
}
