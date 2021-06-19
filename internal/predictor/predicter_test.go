package predictor_test

import (
	"github.com/clambin/sciensano/internal/predictor"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPredictor_Failures(t *testing.T) {
	const (
		batchSize = 7
		dataSize  = 3 * batchSize
	)
	p := predictor.New(batchSize, 1000)

	values := make([]float64, dataSize)
	for i := 0; i < dataSize; i++ {
		values[i] = float64(i)
	}

	_ = p.Learn([][]float64{values})
	_, err := p.Predict([]float64{})

	assert.Error(t, err)
}

func TestForecastSamples_Single(t *testing.T) {
	const (
		batchSize   = 14
		dataBatches = 4
	)

	values := make([]float64, dataBatches*batchSize)
	for i := 0; i < len(values); i++ {
		values[i] = float64(i)
	}

	forecast := predictor.ForecastSamples(batchSize*2, batchSize, "samples", values)

	result := make([]float64, 0)
	for x := range forecast {
		result = append(result, x...)
	}

	assert.Len(t, result, (dataBatches+1)*batchSize)
}

func TestForecastSamples_Double(t *testing.T) {
	const (
		batchSize   = 14
		dataBatches = 4
	)

	series1 := make([]float64, dataBatches*batchSize)
	series2 := make([]float64, dataBatches*batchSize)
	for i := 0; i < len(series1); i++ {
		series1[i] = float64(i)
		series2[i] = float64(-i)
	}

	fc1 := predictor.ForecastSamples(batchSize*2, batchSize, "samples", series1)
	fc2 := predictor.ForecastSamples(batchSize*2, batchSize, "samples", series2)
	output := predictor.ConsolidateSamples(predictor.SingleConsolidator, fc1, fc2)

	result := make([]float64, 0)
	for x := range output {
		result = append(result, x...)
	}

	assert.Len(t, result, 2*(dataBatches+1)*batchSize)
}

func BenchmarkForecastSamples(b *testing.B) {
	const (
		batchSize   = 14
		dataBatches = 40
	)

	series1 := make([]float64, dataBatches*batchSize)
	series2 := make([]float64, dataBatches*batchSize)
	for i := 0; i < len(series1); i++ {
		series1[i] = float64(i)
		series2[i] = float64(-i)
	}

	fc1 := predictor.ForecastSamples(batchSize*2, batchSize, "samples", series1)
	fc2 := predictor.ForecastSamples(batchSize*2, batchSize, "samples", series2)
	output := predictor.ConsolidateSamples(predictor.SingleConsolidator, fc1, fc2)

	result := make([]float64, 0)
	for x := range output {
		result = append(result, x...)
	}

	assert.Len(b, result, 2*(dataBatches+1)*batchSize)

}
