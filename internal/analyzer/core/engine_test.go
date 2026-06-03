package core

import (
	"testing"
)

func TestWeightedScoreEngine_CalculateScore(t *testing.T) {
	thresholds := Thresholds{
		Warning:   40,
		Healthy:   70,
		Excellent: 90,
	}
	engine := NewWeightedScoreEngine(100.0, thresholds)

	tests := []struct {
		name             string
		metrics          []Metric
		expectedScore    float64
		expectedCategory ScoreCategory
	}{
		{
			name: "Perfect Score",
			metrics: []Metric{
				{Score: 100.0, Weight: 1.0},
				{Score: 100.0, Weight: 1.0},
			},
			expectedScore:    100.0,
			expectedCategory: Excellent,
		},
		{
			name: "Healthy Score",
			metrics: []Metric{
				{Score: 80.0, Weight: 1.0},
				{Score: 70.0, Weight: 1.0},
			},
			expectedScore:    75.0,
			expectedCategory: Healthy,
		},
		{
			name: "Warning Score",
			metrics: []Metric{
				{Score: 50.0, Weight: 1.0},
				{Score: 60.0, Weight: 1.0},
			},
			expectedScore:    55.0,
			expectedCategory: Warning,
		},
		{
			name: "Critical Score",
			metrics: []Metric{
				{Score: 20.0, Weight: 1.0},
				{Score: 10.0, Weight: 1.0},
			},
			expectedScore:    15.0,
			expectedCategory: Critical,
		},
		{
			name: "Zero Weight Fallback",
			metrics: []Metric{
				{Score: 100.0, Weight: 0.0},
			},
			expectedScore:    0.0,
			expectedCategory: Critical,
		},
		{
			name:             "Empty Metrics Fallback",
			metrics:          []Metric{},
			expectedScore:    0.0,
			expectedCategory: Critical,
		},
		{
			name: "Weighted Bias",
			metrics: []Metric{
				{Score: 100.0, Weight: 3.0}, // 300
				{Score: 0.0, Weight: 1.0},   // 0
			},
			// Total: 300 / 4 = 75.0
			expectedScore:    75.0,
			expectedCategory: Healthy,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score, category := engine.CalculateScore(tt.metrics)
			if score != tt.expectedScore {
				t.Errorf("expected score %.2f, got %.2f", tt.expectedScore, score)
			}
			if category != tt.expectedCategory {
				t.Errorf("expected category %s, got %s", tt.expectedCategory, category)
			}
		})
	}
}

func TestGetMaxScore(t *testing.T) {
	engine := NewWeightedScoreEngine(150.0, Thresholds{})
	if engine.GetMaxScore() != 150.0 {
		t.Errorf("expected max score 150.0, got %.2f", engine.GetMaxScore())
	}
}
