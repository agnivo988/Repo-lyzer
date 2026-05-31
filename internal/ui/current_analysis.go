package ui

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const currentAnalysisFile = "exports/current_analysis.json"

func SaveCurrentAnalysis(result AnalysisResult) error {
	if err := os.MkdirAll(filepath.Dir(currentAnalysisFile), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(currentAnalysisFile, data, 0o644)
}

func LoadCurrentAnalysis() (*AnalysisResult, error) {
	data, err := os.ReadFile(currentAnalysisFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var result AnalysisResult
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func ClearCurrentAnalysis() error {
	if err := os.MkdirAll(filepath.Dir(currentAnalysisFile), 0o755); err != nil {
		return err
	}
	return os.WriteFile(currentAnalysisFile, []byte("null"), 0o644)
}
