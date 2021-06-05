package predictor_test

import (
	"github.com/clambin/sciensano/internal/predictor"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestPredictor(t *testing.T) {
	const (
		batchSize = 10
		dataSize  = 40 * batchSize
	)
	p := predictor.New(batchSize, 1000)

	values := make([]float64, dataSize)
	for i := 0; i < dataSize; i++ {
		values[i] = 50.0 + (360.0/float64(i+1))*math.Cos(float64(i)*10*math.Pi*2/360) - float64(i)/10.0
	}

	score := p.Learn(values[:dataSize-batchSize])
	assert.Greater(t, score, 0.9)

	predicted, err := p.PredictN(values[dataSize-2*batchSize:dataSize-batchSize], batchSize)
	if assert.NoError(t, err) {
		for i := 0; i < batchSize; i++ {
			assert.Less(t, math.Abs(predicted[i]-values[dataSize-2*batchSize+1]), 10.0)
		}
	}
}

func BenchmarkPredictor(b *testing.B) {
	const (
		batchSize = 10
		dataSize  = 40 * batchSize
	)
	p := predictor.New(batchSize, 1000)

	values := make([]float64, dataSize)
	for i := 0; i < dataSize; i++ {
		values[i] = 50.0 + (360.0/float64(i+1))*math.Cos(float64(i)*10*math.Pi*2/360) - float64(i)/10.0
	}

	score := p.Learn(values[:dataSize-batchSize])
	assert.Greater(b, score, 0.9)

	for i := 0; i+batchSize+1 < dataSize; i++ {
		_, err := p.Predict(values[i : i+batchSize])
		assert.NoError(b, err)
	}

}
