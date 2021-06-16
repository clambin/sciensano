package predictor_test

import (
	"fmt"
	"github.com/clambin/sciensano/internal/predictor"
	"github.com/stretchr/testify/assert"
	"math"
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

// TODO: validating accuracy of models during testing is very expensive
// should move this to separate tuning testing rather than automated unit testing

func TestPredictor_Single(t *testing.T) {
	const (
		batchSize = 14
		dataSize  = 40 * batchSize
	)
	p := predictor.New(batchSize, 1000)

	values := make([]float64, dataSize)
	for i := 0; i < dataSize; i++ {
		values[i] = 50.0 + (360.0/float64(i+1))*math.Cos(float64(i)*10*math.Pi*2/360) - float64(i)/10.0
		// values[i] = float64(i)
	}

	const targetScore = 0.98
	score := 0.0
	attempts := 10

	for score < targetScore && attempts > 0 {
		score = p.Learn([][]float64{values[:dataSize-batchSize]})
		attempts--
	}
	assert.Greater(t, score, targetScore)

	buffer := make([]float64, batchSize)
	copy(buffer, values[dataSize-2*batchSize:dataSize-batchSize])

	for i := 0; i < batchSize; i++ {
		predicted, err := p.Predict(buffer)
		if assert.NoError(t, err) {
			target := values[dataSize-batchSize+i]
			difference := math.Abs(predicted[0] - target)
			assert.Less(t, difference, 5.0, fmt.Sprintf("%d: %.1f <-> %.1f", i, target, predicted[0]))
			buffer = append(buffer[1:], predicted[0])

		} else {
			break
		}
	}
}

func TestPredictor_Multi(t *testing.T) {
	const (
		batchSize = 14
		dataSize  = 40 * batchSize
	)
	p := predictor.New(batchSize, 1000)

	series1 := make([]float64, dataSize)
	for i := 0; i < dataSize; i++ {
		series1[i] = float64(i)
	}
	series2 := make([]float64, dataSize)
	for i := 0; i < dataSize; i++ {
		series2[i] = 50.0 + (360.0/float64(i+1))*math.Cos(float64(i)*10*math.Pi*2/360) - float64(i)/10.0
	}

	input := [][]float64{series1[:dataSize-batchSize], series2[:dataSize-batchSize]}

	const targetScore = 0.98
	score := 0.0
	attempts := 10

	for score < targetScore && attempts > 0 {
		score = p.Learn(input)
		attempts--
	}
	assert.Greater(t, score, targetScore)
}
