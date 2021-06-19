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
		hiddenLayerSizes []int   = []int{100}
		Alpha            float64 // = 0.0001
	)

	r := nn.NewMLPRegressor(hiddenLayerSizes, "relu", "adam", Alpha)
	r.BatchSize = batchSize
	r.MaxIter = maxIter

	return &Predictor{
		regressor: r,
		batchSize: batchSize,
	}
}

func (r *Predictor) Learn(values [][]float64) (score float64) {
	// create the trainX and trainY matrices
	//
	// trainX holds rows of length 'window' of training data
	// trainY holds the next expected values for that window
	trainXData := make([]float64, 0)
	trainYData := make([]float64, 0)
	rows := 0
	for i := 0; i+1+r.batchSize < len(values[0]); i++ {
		trainXData = append(trainXData, values[0][i:i+r.batchSize]...)
		for j := 0; j < len(values); j++ {
			trainYData = append(trainYData, values[j][i+r.batchSize])
		}
		rows++
	}

	trainX, trainY := mat.NewDense(rows, r.batchSize, trainXData), mat.NewDense(rows, len(values), trainYData)

	r.regressor.Fit(trainX, trainY)

	return r.regressor.Score(trainX, trainY)
}

func (r *Predictor) Predict(input []float64) (output []float64, err error) {
	if len(input) != r.batchSize {
		return nil, errors.New("input must be a full BatchSize")
	}

	X := mat.NewDense(1, r.regressor.BatchSize, input)

	var predictions mat.Dense
	fc := r.regressor.Predict(X, &predictions)

	_, cols := fc.Dims()

	for j := 0; j < cols; j++ {
		output = append(output, fc.At(0, j))
	}
	return output, nil
}
