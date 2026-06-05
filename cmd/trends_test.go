package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/agnivo988/Repo-lyzer/internal/github"
	"github.com/agnivo988/Repo-lyzer/internal/predictive"
	"github.com/agnivo988/Repo-lyzer/internal/temporal"
	"github.com/spf13/cobra"
)

func TestRunTrendsForecastOutputsJSON(t *testing.T) {
	originalBuildTimeline := buildTimelineFromGitHub
	originalForecast := forecastHealthFromTimeline
	originalPredictor := newPredictor
	defer func() {
		buildTimelineFromGitHub = originalBuildTimeline
		forecastHealthFromTimeline = originalForecast
		newPredictor = originalPredictor
	}()

	timeline := temporal.NewTimeline("owner", "repo")
	first := temporal.NewSnapshot(time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC), nil)
	first.Metrics.AverageHealth = 70
	second := temporal.NewSnapshot(time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC), nil)
	second.Metrics.AverageHealth = 78
	_ = timeline.AddSnapshot(first)
	_ = timeline.AddSnapshot(second)

	buildTimelineFromGitHub = func(_ *github.Client, _, _ string, _ int) (*temporal.Timeline, error) {
		return timeline, nil
	}
	forecastHealthFromTimeline = func(_ *predictive.Predictor, _ *temporal.Timeline, _ int) (*predictive.ForecastResult, error) {
		return &predictive.ForecastResult{
			Metric: "health",
			Predictions: []predictive.Prediction{
				{
					Timestamp:  time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC),
					Value:      81,
					LowerBound: 79,
					UpperBound: 83,
					Confidence: 0.95,
					Method:     "linear_regression",
				},
			},
		}, nil
	}

	cmd := &cobra.Command{}
	cmd.Flags().Bool("forecast", false, "")
	cmd.Flags().Bool("json", false, "")
	cmd.Flags().String("model", "linear", "")
	cmd.Flags().Int("months", 6, "")
	if err := cmd.Flags().Set("forecast", "true"); err != nil {
		t.Fatalf("set forecast flag: %v", err)
	}
	if err := cmd.Flags().Set("json", "true"); err != nil {
		t.Fatalf("set json flag: %v", err)
	}
	if err := cmd.Flags().Set("model", "linear"); err != nil {
		t.Fatalf("set model flag: %v", err)
	}
	if err := cmd.Flags().Set("months", "2"); err != nil {
		t.Fatalf("set months flag: %v", err)
	}

	var buf bytes.Buffer
	originalStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe stdout: %v", err)
	}
	os.Stdout = w

	resultCh := make(chan error, 1)
	go func() {
		resultCh <- runTrends("owner/repo", cmd)
		_ = w.Close()
	}()

	_, _ = buf.ReadFrom(r)
	os.Stdout = originalStdout

	if err := <-resultCh; err != nil {
		t.Fatalf("runTrends returned error: %v", err)
	}

	var payload struct {
		CurrentHealth int    `json:"current_health"`
		ForecastModel string `json:"forecast_model"`
		Predictions   []struct {
			Method string `json:"method"`
		} `json:"predictions"`
	}
	if err := json.Unmarshal(buf.Bytes(), &payload); err != nil {
		t.Fatalf("unmarshal forecast JSON: %v\noutput: %s", err, buf.String())
	}
	if payload.CurrentHealth != 78 {
		t.Fatalf("current_health = %d, want 78", payload.CurrentHealth)
	}
	if payload.ForecastModel != "linear" {
		t.Fatalf("forecast_model = %q, want linear", payload.ForecastModel)
	}
	if len(payload.Predictions) != 1 {
		t.Fatalf("predictions length = %d, want 1", len(payload.Predictions))
	}
	if payload.Predictions[0].Method != "linear_regression" {
		t.Fatalf("prediction method = %q, want linear_regression", payload.Predictions[0].Method)
	}
}
