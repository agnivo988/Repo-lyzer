package predictive

import (
	"fmt"
	"math"
	"time"

	"github.com/agnivo988/Repo-lyzer/internal/github"
)

// ForecastHealth generates predictions for repository health.
func (p *Predictor) ForecastHealth(commits []github.Commit, months int) (*ForecastResult, error) {
	if len(commits) == 0 {
		return nil, fmt.Errorf("commits list is empty")
	}

	if months <= 0 {
		months = p.ForecastHorizon
	}

	// Group commits by month
	monthlyCounts := make(map[string]float64)
	var monthsOrder []string
	for _, c := range commits {
		monthStr := c.Commit.Author.Date.Format("2006-01")
		if _, exists := monthlyCounts[monthStr]; !exists {
			monthsOrder = append(monthsOrder, monthStr)
		}
		monthlyCounts[monthStr]++
	}

	// Sort months
	// Just use a simple approach for historical array
	var historical []float64
	for _, m := range monthsOrder {
		historical = append(historical, monthlyCounts[m])
	}

	model := NewLinearRegressionModel("commit_velocity")
	err := model.Train(historical)
	if err != nil {
		return nil, err
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
func (p *Predictor) ForecastMaturity(timeline interface{}, months int) (*ForecastResult, error) {
	if timeline == nil {
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
func (p *Predictor) ForecastContributorRisk(commits []github.Commit) ([]ContributorRiskForecast, error) {
	if len(commits) == 0 {
		return nil, fmt.Errorf("commits list is empty")
	}

	authorCounts := make(map[string]int)
	for _, c := range commits {
		if c.Author != nil {
			authorCounts[c.Author.Login]++
		}
	}

	var risks []ContributorRiskForecast
	for author, count := range authorCounts {
		if count > 50 { // Core maintainers
			riskScore, _ := p.EstimateBurnoutRisk(author, commits)
			trajectory := "stable"
			if riskScore > 0.7 {
				trajectory = "worsening"
			}
			risks = append(risks, ContributorRiskForecast{
				ContributorID:     author,
				BurnoutRisk:       riskScore,
				AttritionRisk:     riskScore * 0.8,
				KnowledgeLossRisk: 0.9,
				Trajectory:        trajectory,
				Recommendations:   []string{"Suggest a vacation", "Delegate code reviews"},
			})
		}
	}

	// TODO: Implement contributor risk forecasting
	fmt.Println("[EXPERIMENTAL] Contributor risk forecasting is under development")
	return []ContributorRiskForecast{}, nil
}

// EstimateBurnoutRisk estimates the burnout risk for a specific contributor.
func (p *Predictor) EstimateBurnoutRisk(contributor string, commits []github.Commit) (float64, error) {
	if len(commits) == 0 {
		return 0.0, fmt.Errorf("commits list is empty")
	}

	recentCount := 0
	oldCount := 0
	now := time.Now()

	for _, c := range commits {
		if c.Author != nil && c.Author.Login == contributor {
			age := now.Sub(c.Commit.Author.Date).Hours() / 24.0
			if age < 30 {
				recentCount++
			} else if age < 90 {
				oldCount++
			}
		}
	}

	// Acceleration = recent (1 month) vs average of previous 2 months
	oldMonthlyAvg := float64(oldCount) / 2.0
	if oldMonthlyAvg == 0 {
		return 0.2, nil
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
func (p *Predictor) ForecastDependencyStability(timeline interface{}, months int) (*ForecastResult, error) {
	if timeline == nil {
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
func (p *Predictor) ProjectTechnicalDebt(timeline interface{}, months int) (*ForecastResult, error) {
	if timeline == nil {
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

	// DataCount stores the number of historical points trained on
	DataCount int

	// Name is the model identifier
	ModelName string
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

// Train fits the model to historical data using Least Squares fitting.
func (m *LinearRegressionModel) Train(historical []float64) error {
	n := len(historical)
	if n < 2 {
		return fmt.Errorf("need at least 2 data points for linear regression")
	}

	var sumX, sumY, sumXY, sumX2 float64
	for i, y := range historical {
		x := float64(i)
		sumX += x
		sumY += y
		sumXY += x * y
		sumX2 += x * x
	}

	nF := float64(n)
	denominator := (nF * sumX2) - (sumX * sumX)
	if denominator == 0 {
		m.Slope = 0
		m.Intercept = sumY / nF
	} else {
		m.Slope = ((nF * sumXY) - (sumX * sumY)) / denominator
		m.Intercept = (sumY - (m.Slope * sumX)) / nF
	}

	// Calculate Standard Error
	var sse float64
	for i, y := range historical {
		pred := m.Intercept + m.Slope*float64(i)
		diff := y - pred
		sse += diff * diff
	}
	m.StandardError = math.Sqrt(sse / nF)
	m.DataCount = n

	return nil
}

// Forecast generates predictions for n periods into the future.
func (m *LinearRegressionModel) Forecast(periods int) ([]Prediction, error) {
	if periods < 0 {
		return nil, fmt.Errorf("forecast periods must be non-negative, got %d", periods)
	}

	predictions := make([]Prediction, periods)
	for i := 0; i < periods; i++ {
		futureX := float64(m.DataCount + i)
		val := m.Intercept + (m.Slope * futureX)
		// Ensure non-negative if dealing with counts, but we'll leave it raw here
		predictions[i] = Prediction{
			Value:  val,
			Method: "linear_regression",
		}
	}

	return predictions, nil
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

	lower = make([]float64, periods)
	upper = make([]float64, periods)

	// TODO: Implement confidence interval computation

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

// ForecastIssueAccumulation forecasts the number of open issues based on creation velocity.
func (p *Predictor) ForecastIssueAccumulation(issues []github.Issue, days int) (*ForecastResult, error) {
	if len(issues) == 0 {
		return nil, fmt.Errorf("issues list is empty")
	}

	// Group issues by creation date (last 30 days)
	now := time.Now()
	dailyCounts := make([]float64, 30)

	for _, issue := range issues {
		age := int(now.Sub(issue.CreatedAt).Hours() / 24.0)
		if age >= 0 && age < 30 {
			// Older index first: index 0 = 30 days ago, index 29 = today
			index := 29 - age
			dailyCounts[index]++
		}
	}

	model := NewLinearRegressionModel("issue_accumulation")
	err := model.Train(dailyCounts)
	if err != nil {
		return nil, err
	}

	preds, err := model.Forecast(days)
	if err != nil {
		return nil, err
	}

	trend := "stable"
	if model.Slope > 0.2 {
		trend = "degrading" // issues increasing quickly
	} else if model.Slope < -0.2 {
		trend = "improving" // issue creation slowing down
	}

	return &ForecastResult{
		Metric:          "Issue Accumulation",
		Predictions:     preds,
		Trend:           trend,
		RiskLevel:       "low",
		Recommendations: []string{"Address older issues first", "Prioritize bug triage"},
	}, nil
}
