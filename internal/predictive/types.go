// Package predictive provides predictive modeling and repository state forecasting.
package predictive

import "time"

// Prediction represents a single forecasted value with uncertainty bounds.
type Prediction struct {
	// Timestamp is when this prediction is for
	Timestamp time.Time

	// Value is the predicted value
	Value float64

	// LowerBound is the lower bound of the confidence interval
	LowerBound float64

	// UpperBound is the upper bound of the confidence interval
	UpperBound float64

	// Confidence is the confidence level [0, 1]
	Confidence float64

	// Method is the prediction method used (e.g., "linear_regression", "exponential_smoothing")
	Method string
}

// ForecastResult contains complete forecasting output for a metric.
type ForecastResult struct {
	// Metric name
	Metric string

	// Predictions contains the forecasted values
	Predictions []Prediction

	// Trend: "improving", "stable", "degrading"
	Trend string

	// RiskLevel: "low", "medium", "high"
	RiskLevel string

	// Recommendations lists suggested actions based on the forecast
	Recommendations []string

	// ConfidenceScore [0, 1] indicating overall forecast reliability
	ConfidenceScore float64

	// BaselineMean is the mean of historical data
	BaselineMean float64

	// BaselineStdDev is the standard deviation of historical data
	BaselineStdDev float64
}

// ContributorRiskForecast predicts contributor-related risks.
type ContributorRiskForecast struct {
	// ContributorID identifies the contributor
	ContributorID string

	// BurnoutRisk [0, 1] indicates likelihood of burnout
	BurnoutRisk float64

	// AttritionRisk [0, 1] indicates likelihood of departure
	AttritionRisk float64

	// KnowledgeLossRisk [0, 1] indicates project impact if contributor leaves
	KnowledgeLossRisk float64

	// Trajectory: "improving", "stable", "worsening"
	Trajectory string

	// Recommendations lists ways to support this contributor
	Recommendations []string
}

// PredictiveModel defines the interface for forecasting models.
type PredictiveModel interface {
	// Train fits the model to historical data
	Train(historical []float64) error

	// Forecast generates predictions for n periods into the future
	Forecast(periods int) ([]Prediction, error)

	// ConfidenceIntervals computes confidence bounds for predictions
	ConfidenceIntervals(periods int, confidenceLevel float64) (lower, upper []float64, error error)

	// Name returns the model name
	Name() string

	// Parameters returns model-specific parameters
	Parameters() map[string]interface{}
}

// Predictor provides forecasting capabilities for repository metrics.
type Predictor struct {
	// Models maps metric names to their prediction models
	Models map[string]PredictiveModel

	// HistoricalData stores historical metric values for training
	HistoricalData map[string][]float64

	// ForecastHorizon is the number of periods to forecast (default: 12 months)
	ForecastHorizon int

	// ConfidenceLevel [0.8, 0.99] for confidence intervals (default: 0.95)
	ConfidenceLevel float64
}

// NewPredictor creates a new predictor with default settings.
func NewPredictor() *Predictor {
	return &Predictor{
		Models:          make(map[string]PredictiveModel),
		HistoricalData:  make(map[string][]float64),
		ForecastHorizon: 12,
		ConfidenceLevel: 0.95,
	}
}
