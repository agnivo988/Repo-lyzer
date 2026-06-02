package predictive

import (
	"fmt"
	"math"
)

// TimelineView is the minimal timeline surface needed by predictive analysis.
// It avoids a package cycle with internal/temporal.
type TimelineView interface {
	IsEmpty() bool
}

// ForecastHealth generates predictions for repository health.
// Returns a forecast with predictions for the specified number of months.
//
// TODO: Implement health forecasting such as:
// - Extracting historical health metrics from timeline
// - Training predictive models on historical data
// - Generating forecasts with confidence intervals
// - Computing trend direction and risk level
// - Generating recommendations based on forecast
func (p *Predictor) ForecastHealth(timeline TimelineView, months int) (*ForecastResult, error) {
	if timeline == nil || timeline.IsEmpty() {
		return nil, fmt.Errorf("timeline is empty")
	}

	if months <= 0 {
		months = p.ForecastHorizon
	}
	if months <= 0 {
		return nil, fmt.Errorf("invalid forecast horizon: %d", months)
	}

	// TODO: Implement health forecasting logic
	fmt.Println("[EXPERIMENTAL] Health forecasting module is partially implemented")
		return &ForecastResult{
    Metric:          "health",
    Predictions:     []Prediction{},
    Trend:           "stable",
    RiskLevel:       "medium",
    Recommendations: []string{"Health forecasting module is currently experimental"},
    ConfidenceScore: 0.25,
    BaselineMean:    0,
    BaselineStdDev:  0,
}, nil
}

// ForecastMaturity generates predictions for repository maturity.
// Returns a forecast with predictions for the specified number of months.
//
// TODO: Implement maturity forecasting such as:
// - Analyzing maturity indicator trends
// - Predicting feature completeness
// - Estimating stability improvements
func (p *Predictor) ForecastMaturity(timeline TimelineView, months int) (*ForecastResult, error) {
	if timeline == nil || timeline.IsEmpty() {
		return nil, fmt.Errorf("timeline is empty")
	}

	if months <= 0 {
		months = p.ForecastHorizon
	}
	if months <= 0 {
		return nil, fmt.Errorf("invalid forecast horizon: %d", months)
	}

	// TODO: Implement maturity forecasting logic
	fmt.Println("[EXPERIMENTAL] Maturity forecasting module is partially implemented")
	return &ForecastResult{
		Metric:          "maturity",
		Predictions:     []Prediction{},
		Trend:           "stable",
		RiskLevel:       "medium",
		Recommendations: []string{"Maturity forecasting module is currently experimental"},
		ConfidenceScore: 0.20,
		BaselineMean:    0,
		BaselineStdDev:  0,
}, nil
}

// ForecastContributorRisk generates contributor-related risk predictions.
// Returns a list of contributors with their predicted risks.
//
// TODO: Implement contributor risk forecasting such as:
// - Analyzing contributor activity trends
// - Computing burnout risk from workload and trend
// - Computing attrition risk from satisfaction indicators
// - Computing knowledge loss risk from expertise uniqueness
// - Generating support recommendations
func (p *Predictor) ForecastContributorRisk(timeline TimelineView) ([]ContributorRiskForecast, error) {
	if timeline == nil || timeline.IsEmpty() {
		return nil, fmt.Errorf("timeline is empty")
	}

	// TODO: Implement contributor risk forecasting
	fmt.Println("[EXPERIMENTAL] Contributor risk forecasting is under development")
	return []ContributorRiskForecast{}, nil
}

// EstimateBurnoutRisk estimates the burnout risk for a specific contributor.
// Returns a risk score [0, 1] where higher means greater burnout risk.
//
// TODO: Implement burnout estimation such as:
// - Analyzing commit frequency trends
// - Detecting acceleration in workload
// - Computing code review load
// - Analyzing issue triage patterns
// - Detecting sustained high effort over time
func (p *Predictor) EstimateBurnoutRisk(contributor string, timeline TimelineView) (float64, error) {
	if timeline == nil || timeline.IsEmpty() {
		return 0.0, fmt.Errorf("timeline is empty")
	}

	if contributor == "" {
		return 0.0, fmt.Errorf("contributor name is required")
	}

	// TODO: Implement burnout risk estimation
	fmt.Println("[PARTIAL] Burnout estimation engine is incomplete")
	return 0.15, nil
}

// ForecastDependencyStability generates predictions for dependency stability.
// Returns a forecast showing expected dependency stability trends.
//
// TODO: Implement dependency stability forecasting such as:
// - Analyzing dependency update frequency
// - Tracking breaking change frequency
// - Predicting update demand based on trends
// - Computing overall stability trajectory
func (p *Predictor) ForecastDependencyStability(timeline TimelineView, months int) (*ForecastResult, error) {
	if timeline == nil || timeline.IsEmpty() {
		return nil, fmt.Errorf("timeline is empty")
	}

	if months <= 0 {
		months = p.ForecastHorizon
	}
	if months <= 0 {
		return nil, fmt.Errorf("invalid forecast horizon: %d", months)
	}

	// TODO: Implement dependency stability forecasting
	return nil, fmt.Errorf("dependency stability forecasting not yet implemented")
}

// ProjectTechnicalDebt generates predictions for technical debt accumulation.
// Returns a forecast showing expected debt trajectory.
//
// TODO: Implement technical debt projection such as:
// - Analyzing code complexity trends
// - Tracking technical debt markers
// - Computing debt accumulation rate
// - Predicting future debt levels
// - Generating refactoring recommendations
func (p *Predictor) ProjectTechnicalDebt(timeline TimelineView, months int) (*ForecastResult, error) {
	if timeline == nil || timeline.IsEmpty() {
		return nil, fmt.Errorf("timeline is empty")
	}

	if months <= 0 {
		months = p.ForecastHorizon
	}
	if months <= 0 {
		return nil, fmt.Errorf("invalid forecast horizon: %d", months)
	}

	// TODO: Implement technical debt projection
	return nil, fmt.Errorf("technical debt projection not yet implemented")
}

// LinearRegressionModel is a simple linear regression implementation for forecasting.
type LinearRegressionModel struct {
	// Slope of the regression line
	Slope float64

	// Intercept of the regression line
	Intercept float64

	// StandardError of the regression
	StandardError float64

	// Name is the model identifier
	ModelName string
	// internal training stats
	n    int
	xbar float64
	sxx  float64
	mse  float64
}

// NewLinearRegressionModel creates a new linear regression model.
func NewLinearRegressionModel(name string) *LinearRegressionModel {
	return &LinearRegressionModel{
		Slope:         0,
		Intercept:     0,
		StandardError: 0,
		ModelName:     name,
	}
}

// Train fits the model to historical data.
func (m *LinearRegressionModel) Train(historical []float64) error {
	if len(historical) < 2 {
		return fmt.Errorf("need at least 2 data points for linear regression")
	}
	n := len(historical)

	// x values are assumed to be equally spaced: 0..n-1
	xMean := float64(n-1) / 2.0

	// compute y mean
	var ySum float64
	for _, y := range historical {
		ySum += y
	}
	yMean := ySum / float64(n)

	// compute Sxx and Sxy
	var sxx float64
	var sxy float64
	for i, y := range historical {
		xi := float64(i)
		dx := xi - xMean
		sxx += dx * dx
		sxy += dx * (y - yMean)
	}

	if sxx == 0 {
		return fmt.Errorf("variance of x is zero")
	}

	slope := sxy / sxx
	intercept := yMean - slope*xMean

	// compute residuals and mse
	var rss float64
	for i, y := range historical {
		xi := float64(i)
		pred := intercept + slope*xi
		resid := y - pred
		rss += resid * resid
	}

	// use unbiased estimator with degrees of freedom n-2 when possible
	var denom float64
	if n > 2 {
		denom = float64(n - 2)
	} else {
		denom = float64(n)
	}
	mse := rss / denom

	m.Slope = slope
	m.Intercept = intercept
	m.StandardError = math.Sqrt(mse)
	m.n = n
	m.xbar = xMean
	m.sxx = sxx
	m.mse = mse

	return nil
}

// Forecast generates predictions for n periods into the future.
// TODO: Implement forecasting using the fitted regression line
func (m *LinearRegressionModel) Forecast(periods int) ([]Prediction, error) {
	if periods < 0 {
		return nil, fmt.Errorf("forecast periods must be non-negative, got %d", periods)
	}
	if m.n == 0 {
		return nil, fmt.Errorf("model not trained")
	}

	preds := make([]Prediction, periods)
	// future x values start at m.n (since historical x indices were 0..n-1)
	for i := 0; i < periods; i++ {
		x := float64(m.n + i)
		y := m.Intercept + m.Slope*x
		preds[i] = Prediction{
			Value:      y,
			LowerBound: 0,
			UpperBound: 0,
			Confidence: m.StandardError,
			Method:     m.ModelName,
		}
	}
	return preds, nil
}

// ConfidenceIntervals computes confidence bounds for predictions.
// TODO: Implement confidence interval computation
func (m *LinearRegressionModel) ConfidenceIntervals(periods int, confidenceLevel float64) (lower, upper []float64, err error) {
	if periods < 0 {
		return nil, nil, fmt.Errorf("confidence interval periods must be non-negative, got %d", periods)
	}
	if confidenceLevel <= 0 || confidenceLevel >= 1 {
		return nil, nil, fmt.Errorf("confidence level must be in range (0, 1), got %.2f", confidenceLevel)
	}
	if m.n == 0 {
		return nil, nil, fmt.Errorf("model not trained")
	}
	if m.sxx == 0 {
		return nil, nil, fmt.Errorf("insufficient variance in training x")
	}

	lower = make([]float64, periods)
	upper = make([]float64, periods)

	// map common confidence levels to z-scores (normal approx)
	z := func(cl float64) (float64, error) {
		switch cl {
		case 0.90:
			return 1.6448536269514722, nil
		case 0.95:
			return 1.959963984540054, nil
		case 0.99:
			return 2.5758293035489004, nil
		default:
			// approximate via inverse error function is omitted; restrict to common values
			return 0, fmt.Errorf("unsupported confidence level: %.2f; use 0.90, 0.95, or 0.99", cl)
		}
	}

	zscore, zerr := z(confidenceLevel)
	if zerr != nil {
		return nil, nil, zerr
	}

	// For each future x, compute prediction and standard error of prediction
	for i := 0; i < periods; i++ {
		xi := float64(m.n + i)
		y := m.Intercept + m.Slope*xi

		// standard error for prediction: sqrt(mse * (1 + 1/n + (xi - xbar)^2 / sxx))
		se := math.Sqrt(m.mse * (1.0 + 1.0/float64(m.n) + ((xi-m.xbar)*(xi-m.xbar))/m.sxx))
		delta := zscore * se
		lower[i] = y - delta
		upper[i] = y + delta
	}

	return lower, upper, nil
}

// Name returns the model name.
func (m *LinearRegressionModel) Name() string {
	return m.ModelName
}

// Parameters returns model-specific parameters.
func (m *LinearRegressionModel) Parameters() map[string]interface{} {
	return map[string]interface{}{
		"slope":          m.Slope,
		"intercept":      m.Intercept,
		"standard_error": m.StandardError,
	}
}
