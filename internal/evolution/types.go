// Package evolution provides repository evolution analysis, including pattern detection,
// architectural drift detection, and risk indicator computation.
package evolution

import "time"

// EvolutionPattern describes a detected evolution pattern in the repository.
type EvolutionPattern struct {
	// Name of the pattern (e.g., "Increasing Complexity", "Contributor Consolidation")
	Name string

	// Description explains what this pattern means
	Description string

	// StartTime is when this pattern began
	StartTime time.Time

	// EndTime is when this pattern ended (or time.Time{} if ongoing)
	EndTime time.Time

	// Indicators contains metric names and their values
	Indicators map[string]float64

	// Severity describes the pattern severity: "low", "medium", "high"
	Severity string

	// Confidence is the confidence score [0, 1]
	Confidence float64

	// Affected lists items affected by this pattern (subsystem IDs, contributor IDs, etc.)
	Affected []string
}

// DriftIndicator represents detected architectural drift in a subsystem.
type DriftIndicator struct {
	// SubsystemID identifies the subsystem showing drift
	SubsystemID string

	// MetricName is the metric showing drift (e.g., "complexity", "coupling")
	MetricName string

	// Direction indicates change direction: "increasing", "decreasing", or "stable"
	Direction string

	// Magnitude is the magnitude of change (0-1)
	Magnitude float64

	// StartValue is the value at the start of the period
	StartValue float64

	// EndValue is the value at the end of the period
	EndValue float64

	// TimeSpan is the duration over which drift was observed
	TimeSpan time.Duration

	// Threshold is the threshold that triggered this indicator
	Threshold float64

	// Severity: "low", "medium", "high"
	Severity string
}

// RiskIndicator represents a detected risk in the repository.
type RiskIndicator struct {
	// Category: "complexity", "contributor", "dependency", "sustainability"
	Category string

	// Name is a descriptive name for this risk
	Name string

	// Severity: "low", "medium", "high", "critical"
	Severity string

	// Affected lists affected items (subsystem IDs, contributors, etc.)
	Affected []string

	// Current is the current value of the risk metric
	Current float64

	// Threshold is the threshold that triggered this risk
	Threshold float64

	// Trajectory: "improving", "stable", "worsening"
	Trajectory string

	// Recommendations lists suggested actions to mitigate this risk
	Recommendations []string
}

// ComplexityReport summarizes complexity analysis for a repository.
type ComplexityReport struct {
	// AverageComplexity is the average complexity across all subsystems
	AverageComplexity float64

	// MaxComplexity is the maximum complexity among all subsystems
	MaxComplexity float64

	// ComplexityGrowthRate is the rate of complexity growth per month
	ComplexityGrowthRate float64

	// HighComplexitySubsystems lists subsystems with high complexity
	HighComplexitySubsystems []string

	// CriticalSubsystems lists subsystems with high complexity and low knowledge distribution
	CriticalSubsystems []string

	// OverallHealthScore is a score [0, 100] reflecting complexity health
	OverallHealthScore int
}

// ContributorRole describes a contributor's role and evolution in the project.
type ContributorRole struct {
	// ContributorID identifies the contributor
	ContributorID string

	// Role: "core", "active", "occasional", "inactive"
	Role string

	// CommitCount is the total number of commits by this contributor
	CommitCount int

	// FilesTouched is the number of unique files touched
	FilesTouched int

	// Expertise describes the contributor's areas of expertise
	Expertise []string

	// KnowledgeConcentration: 0-1, higher = more knowledge concentrated in this contributor
	KnowledgeConcentration float64

	// ActivityTrend: "increasing", "stable", "decreasing"
	ActivityTrend string

	// RiskScore: 0-1, higher = contributor loss would be more damaging
	RiskScore float64
}

// Bottleneck represents a knowledge bottleneck in the repository.
type Bottleneck struct {
	// ContributorID is the bottleneck person
	ContributorID string

	// CriticalAreas lists areas of expertise unique or near-unique to this person
	CriticalAreas []string

	// RiskLevel: "low", "medium", "high", "critical"
	RiskLevel string

	// Recommendations lists actions to distribute knowledge
	Recommendations []string

	// ReplacementCost estimates the cost/effort to replace this person's knowledge
	ReplacementCost string // "low", "medium", "high"
}

// Detector performs evolution analysis on temporal repository data.
type Detector struct {
	// Configuration parameters for detection
	ComplexityThreshold float64
	DriftThreshold      float64
	RiskThreshold       float64

	// Results storage
	Patterns     []EvolutionPattern
	DriftIndices []DriftIndicator
	RiskIndices  []RiskIndicator
}

// NewDetector creates a new evolution detector with default parameters.
func NewDetector() *Detector {
	return &Detector{
		ComplexityThreshold: 0.7,
		DriftThreshold:      0.5,
		RiskThreshold:       0.6,
		Patterns:            make([]EvolutionPattern, 0),
		DriftIndices:        make([]DriftIndicator, 0),
		RiskIndices:         make([]RiskIndicator, 0),
	}
}
