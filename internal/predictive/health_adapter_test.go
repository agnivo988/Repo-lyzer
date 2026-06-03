package predictive

import (
	"testing"
	"time"

	"github.com/agnivo988/Repo-lyzer/internal/temporal"
)

func newTestTimeline(healthValues ...int) *temporal.Timeline {
	timeline := temporal.NewTimeline("owner", "repo")
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	for i, health := range healthValues {
		snapshot := temporal.NewSnapshot(base.AddDate(0, i, 0), nil)
		snapshot.Metrics.AverageHealth = health
		_ = timeline.AddSnapshot(snapshot)
	}
	return timeline
}

func TestExtractMonthlyHealthSeriesNilTimeline(t *testing.T) {
	series := ExtractMonthlyHealthSeries(nil)
	if len(series) != 0 {
		t.Fatalf("expected empty series for nil timeline, got %v", series)
	}
}

func TestExtractMonthlyHealthSeriesEmptyTimeline(t *testing.T) {
	timeline := temporal.NewTimeline("owner", "repo")
	series := ExtractMonthlyHealthSeries(timeline)
	if len(series) != 0 {
		t.Fatalf("expected empty series for empty timeline, got %v", series)
	}
}

func TestExtractMonthlyHealthSeriesSingleSnapshot(t *testing.T) {
	timeline := newTestTimeline(42)
	series := ExtractMonthlyHealthSeries(timeline)
	if len(series) != 1 {
		t.Fatalf("expected 1 health value, got %d", len(series))
	}
	if series[0] != 42 {
		t.Fatalf("expected health value 42, got %v", series[0])
	}
}

func TestExtractMonthlyHealthSeriesMultipleSnapshots(t *testing.T) {
	timeline := newTestTimeline(40, 50, 60)
	series := ExtractMonthlyHealthSeries(timeline)
	if len(series) != 3 {
		t.Fatalf("expected 3 health values, got %d", len(series))
	}
	want := []float64{40, 50, 60}
	for i := range want {
		if series[i] != want[i] {
			t.Fatalf("series[%d] = %v, want %v", i, series[i], want[i])
		}
	}
}

func TestForecastHealthFromTimelineNilTimeline(t *testing.T) {
	predictor := NewPredictor()
	_, err := ForecastHealthFromTimeline(predictor, nil, 3)
	if err == nil {
		t.Fatal("expected error for nil timeline")
	}
}

func TestForecastHealthFromTimelineEmptyTimeline(t *testing.T) {
	predictor := NewPredictor()
	timeline := temporal.NewTimeline("owner", "repo")
	_, err := ForecastHealthFromTimeline(predictor, timeline, 3)
	if err == nil {
		t.Fatal("expected error for empty timeline")
	}
}

func TestForecastHealthFromTimelineSingleSnapshot(t *testing.T) {
	predictor := NewPredictor()
	timeline := newTestTimeline(55)
	_, err := ForecastHealthFromTimeline(predictor, timeline, 3)
	if err == nil {
		t.Fatal("expected error for insufficient data")
	}
}

func TestForecastHealthFromTimelineSuccess(t *testing.T) {
	predictor := NewPredictor()
	predictor.ForecastHorizon = 4
	predictor.ConfidenceLevel = 0.95
	timeline := newTestTimeline(40, 50, 60, 70)

	result, err := ForecastHealthFromTimeline(predictor, timeline, 2)
	if err != nil {
		t.Fatalf("ForecastHealthFromTimeline failed: %v", err)
	}
	if result == nil {
		t.Fatal("expected forecast result")
	}
	if result.Metric != "health" {
		t.Fatalf("unexpected metric: %s", result.Metric)
	}
	if len(result.Predictions) != 2 {
		t.Fatalf("expected 2 predictions, got %d", len(result.Predictions))
	}
	if result.Trend != "improving" {
		t.Fatalf("expected improving trend, got %s", result.Trend)
	}
	if result.BaselineMean <= 0 {
		t.Fatalf("expected positive baseline mean, got %v", result.BaselineMean)
	}
	if result.BaselineStdDev <= 0 {
		t.Fatalf("expected positive baseline stddev, got %v", result.BaselineStdDev)
	}
	for i, prediction := range result.Predictions {
		if prediction.Method != "linear_regression" {
			t.Fatalf("prediction %d method = %s, want linear_regression", i, prediction.Method)
		}
		if prediction.UpperBound < prediction.LowerBound {
			t.Fatalf("prediction %d has invalid bounds: [%v, %v]", i, prediction.LowerBound, prediction.UpperBound)
		}
	}
}

func TestForecastHealthFromTimelineAnchorsTimestampsToLatestSnapshot(t *testing.T) {
	predictor := NewPredictor()
	predictor.ForecastHorizon = 2
	timeline := newTestTimeline(40, 50, 60, 70)

	result, err := ForecastHealthFromTimeline(predictor, timeline, 2)
	if err != nil {
		t.Fatalf("ForecastHealthFromTimeline failed: %v", err)
	}

	want := []time.Time{
		time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
	}

	if len(result.Predictions) != len(want) {
		t.Fatalf("unexpected prediction count: got %d want %d", len(result.Predictions), len(want))
	}

	for i := range want {
		if !result.Predictions[i].Timestamp.Equal(want[i]) {
			t.Fatalf("prediction %d timestamp = %v, want %v", i, result.Predictions[i].Timestamp, want[i])
		}
	}
}
