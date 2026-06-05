package predictive

import (
	"math"
	"testing"
)

func approxEqual(a, b, tol float64) bool {
	return math.Abs(a-b) <= tol
}

func TestLinearRegressionTrainForecast(t *testing.T) {
	// y = 2*x + 5 for x=0..4 -> [5,7,9,11,13]
	hist := []float64{5, 7, 9, 11, 13}
	m := NewLinearRegressionModel("linear_test")
	if err := m.Train(hist); err != nil {
		t.Fatalf("Train failed: %v", err)
	}

	if !approxEqual(m.Slope, 2.0, 1e-9) {
		t.Fatalf("unexpected slope: got %v want ~2.0", m.Slope)
	}
	if !approxEqual(m.Intercept, 5.0, 1e-9) {
		t.Fatalf("unexpected intercept: got %v want ~5.0", m.Intercept)
	}

	preds, err := m.Forecast(3)
	if err != nil {
		t.Fatalf("Forecast failed: %v", err)
	}
	if len(preds) != 3 {
		t.Fatalf("expected 3 predictions, got %d", len(preds))
	}

	// expected values for x=5,6,7 -> 15,17,19
	want := []float64{15, 17, 19}
	for i := 0; i < 3; i++ {
		if !approxEqual(preds[i].Value, want[i], 1e-6) {
			t.Fatalf("pred[%d] = %v; want %v", i, preds[i].Value, want[i])
		}
	}
}

func TestConfidenceIntervals(t *testing.T) {
	// Slight noise avoids a perfectly degenerate fit, which would produce a zero-width interval.
	hist := []float64{5, 7.1, 8.9, 11.2, 13.05}
	m := NewLinearRegressionModel("linear_ci")
	if err := m.Train(hist); err != nil {
		t.Fatalf("Train failed: %v", err)
	}

	periods := 2
	lower, upper, err := m.ConfidenceIntervals(periods, 0.95)
	if err != nil {
		t.Fatalf("ConfidenceIntervals failed: %v", err)
	}
	if len(lower) != periods || len(upper) != periods {
		t.Fatalf("unexpected CI lengths: got %d/%d", len(lower), len(upper))
	}

	preds, err := m.Forecast(periods)
	if err != nil {
		t.Fatalf("Forecast failed: %v", err)
	}

	for i := 0; i < periods; i++ {
		if !(lower[i] <= preds[i].Value && preds[i].Value <= upper[i]) {
			t.Fatalf("prediction %v not inside CI [%v, %v]", preds[i].Value, lower[i], upper[i])
		}
	}

	if !(upper[0] > lower[0]) {
		t.Fatalf("expected a non-zero width confidence interval, got [%v, %v]", lower[0], upper[0])
	}
}
