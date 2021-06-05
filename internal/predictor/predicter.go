package predictor

import (
	"errors"
	nn "github.com/pa-m/sklearn/neural_network"
	"gonum.org/v1/gonum/mat"
)

type Predictor struct {
	regressor *nn.MLPRegressor
	batchSize int
}

func New(batchSize int, maxIter int) *Predictor {
	var (
		hiddenLayerSizes []int
		Alpha            float64 = 0
	)

	r := nn.NewMLPRegressor(hiddenLayerSizes, "relu", "adam", Alpha)
	r.BatchSize = batchSize
	r.MaxIter = maxIter

	return &Predictor{
		regressor: r,
		batchSize: batchSize,
	}
}

func (r *Predictor) Learn(values []float64) (score float64) {
	// create the trainX and trainY matrices
	//
	// trainX holds rows of length 'window' of training data
	// trainY holds the next expected value for that window
	rows := 0
	trainXData := make([]float64, 0)
	trainYData := make([]float64, 0)
	for i := 0; i < len(values)-r.batchSize; i++ {
		trainXData = append(trainXData, values[i:i+r.batchSize]...)
		trainYData = append(trainYData, values[i+r.batchSize])
		rows++
	}

	trainX, trainY := mat.NewDense(rows, r.batchSize, trainXData), mat.NewDense(rows, 1, trainYData)

	r.regressor.Fit(trainX, trainY)

	return r.regressor.Score(trainX, trainY)
}

func (r *Predictor) Predict(values []float64) (value float64, err error) {
	if len(values) != r.batchSize {
		return 0, errors.New("input must be a full batchSize")
	}

	input := mat.NewDense(1, r.regressor.BatchSize, values)

	var predictions mat.Dense
	forecast := r.regressor.Predict(input, &predictions)

	return forecast.At(0, 0), nil
}

func (r *Predictor) PredictN(input []float64, count int) (output []float64, err error) {
	if len(input) != r.batchSize {
		return nil, errors.New("input must be a full batchSize")
	}

	buffer := make([]float64, r.batchSize)
	copy(buffer, input)

	var value float64
	for i := 0; i < count && err == nil; i++ {
		if value, err = r.Predict(buffer); err == nil {
			output = append(output, value)
			buffer = append(buffer[1:], value)
		}
	}
	return
}
