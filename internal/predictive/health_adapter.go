package predictive

import (
	"fmt"
	"math"
	"sort"

	"github.com/agnivo988/Repo-lyzer/internal/temporal"
)

// ExtractMonthlyHealthSeries returns snapshot health values in chronological order.
func ExtractMonthlyHealthSeries(t *temporal.Timeline) []float64 {
	if t == nil || len(t.Snapshots) == 0 {
		return []float64{}
	}

	snapshots := make([]*temporal.Snapshot, 0, len(t.Snapshots))
	for _, snapshot := range t.Snapshots {
		if snapshot != nil {
			snapshots = append(snapshots, snapshot)
		}
	}

	if len(snapshots) == 0 {
		return []float64{}
	}

	sort.SliceStable(snapshots, func(i, j int) bool {
		return snapshots[i].Timestamp.Before(snapshots[j].Timestamp)
	})

	series := make([]float64, 0, len(snapshots))
	for _, snapshot := range snapshots {
		series = append(series, float64(snapshot.Metrics.AverageHealth))
	}

	return series
}

// ForecastHealthFromTimeline trains a linear regression model from timeline health values and forecasts future health.
func ForecastHealthFromTimeline(predictor *Predictor, timeline *temporal.Timeline, months int) (*ForecastResult, error) {
	if predictor == nil {
		return nil, fmt.Errorf("predictor is nil")
	}
	if timeline == nil {
		return nil, fmt.Errorf("timeline is nil")
	}

	historical := ExtractMonthlyHealthSeries(timeline)
	if len(historical) < 2 {
		return nil, fmt.Errorf("need at least 2 health data points to forecast, got %d", len(historical))
	}

	if months <= 0 {
		months = predictor.ForecastHorizon
	}
	if months <= 0 {
		return nil, fmt.Errorf("invalid forecast horizon: %d", months)
	}

	model := NewLinearRegressionModel("linear_regression")
	if err := model.Train(historical); err != nil {
		return nil, fmt.Errorf("train linear regression model: %w", err)
	}

	predictions, err := model.Forecast(months)
	if err != nil {
		return nil, fmt.Errorf("forecast health: %w", err)
	}

	anchor := timeline.EndTime
	if anchor.IsZero() && len(timeline.Snapshots) > 0 {
		anchor = timeline.Snapshots[len(timeline.Snapshots)-1].Timestamp
	}
	for i := range predictions {
		predictions[i].Timestamp = anchor.AddDate(0, i+1, 0)
	}

	confidenceLevel := predictor.ConfidenceLevel
	if confidenceLevel <= 0 {
		confidenceLevel = 0.95
	}

	lower, upper, err := model.ConfidenceIntervals(months, confidenceLevel)
	if err != nil {
		return nil, fmt.Errorf("confidence intervals: %w", err)
	}

	for i := range predictions {
		predictions[i].LowerBound = lower[i]
		predictions[i].UpperBound = upper[i]
		predictions[i].Confidence = predictor.ConfidenceLevel
	}

	trend := "stable"
	switch {
	case model.Slope > 0.01:
		trend = "improving"
	case model.Slope < -0.01:
		trend = "degrading"
	}

	return &ForecastResult{
		Metric:          "health",
		Predictions:     predictions,
		Trend:           trend,
		RiskLevel:       "medium",
		Recommendations: []string{"Monitor health trend", "Review recent repository activity"},
		ConfidenceScore: confidenceLevel,
		BaselineMean:    mean(historical),
		BaselineStdDev:  stddev(historical),
	}, nil
}

func mean(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	var sum float64
	for _, value := range values {
		sum += value
	}
	return sum / float64(len(values))
}

func stddev(values []float64) float64 {
	if len(values) < 2 {
		return 0
	}
	mu := mean(values)
	var sumSq float64
	for _, value := range values {
		delta := value - mu
		sumSq += delta * delta
	}
	return math.Sqrt(sumSq / float64(len(values)-1))
}
