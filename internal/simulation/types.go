// Package simulation provides repository evolution simulation and scenario testing.
package simulation

import "time"

// SimulationScenario defines a what-if scenario for repository evolution simulation.
type SimulationScenario struct {
	// Name is a descriptive name for the scenario
	Name string

	// Description explains what this scenario models
	Description string

	// ScenarioType: "contributor_departure", "subsystem_growth", "dependency_update", etc.
	ScenarioType string

	// Parameters are scenario-specific configuration values
	Parameters map[string]interface{}

	// Duration is how long to simulate
	Duration time.Duration

	// StartTime is when the simulation begins (optional, defaults to now)
	StartTime time.Time
}

// SimulationResult contains the outcomes of a simulation run.
type SimulationResult struct {
	// Scenario that was simulated
	Scenario SimulationScenario

	// InitialState is the repository state at simulation start
	InitialState map[string]float64

	// FinalState is the repository state at simulation end
	FinalState map[string]float64

	// HealthTrajectory shows repository health over time during simulation
	HealthTrajectory []float64

	// RiskTrajectory shows risk levels over time during simulation
	RiskTrajectory []float64

	// ComplexityTrajectory shows complexity over time during simulation
	ComplexityTrajectory []float64

	// Timestamps for each trajectory point
	Timestamps []time.Time

	// KeyFindings summarizes important outcomes
	KeyFindings []string

	// Recommendations lists suggested actions based on simulation outcome
	Recommendations []string

	// HealthChange: absolute change in health score
	HealthChange float64

	// RiskChange: absolute change in risk level
	RiskChange float64

	// Success: true if simulation ran without errors
	Success bool
}

// ScenarioRunner executes simulations based on repository history.
type ScenarioRunner struct {
	// Repository name and owner
	RepoName string
	Owner    string

	// Simulation timestep in days
	TimestepDays int

	// Random seed for reproducibility
	RandomSeed int64

	// Results storage
	Results []SimulationResult
}

// NewScenarioRunner creates a new scenario runner with default settings.
func NewScenarioRunner(owner, repoName string) *ScenarioRunner {
	return &ScenarioRunner{
		RepoName:     repoName,
		Owner:        owner,
		TimestepDays: 1,
		RandomSeed:   42,
		Results:      make([]SimulationResult, 0),
	}
}

// NewScenario creates a new simulation scenario.
func NewScenario(name, scenarioType string, duration time.Duration) *SimulationScenario {
	return &SimulationScenario{
		Name:         name,
		ScenarioType: scenarioType,
		Duration:     duration,
		Parameters:   make(map[string]interface{}),
		StartTime:    time.Now(),
	}
}

// PredefinedScenarios contains scenario templates.
var PredefinedScenarios = map[string]*SimulationScenario{
	"key_contributor_departure": {
		Name:         "Key Contributor Departure",
		Description:  "Simulates the effect of losing a key contributor",
		ScenarioType: "contributor_departure",
		Duration:     90 * 24 * time.Hour, // 90 days
		Parameters: map[string]interface{}{
			"contributor_id":     "unknown",
			"replacement_months": 6,
		},
	},
	"rapid_subsystem_growth": {
		Name:         "Rapid Subsystem Growth",
		Description:  "Simulates rapid growth in a subsystem",
		ScenarioType: "subsystem_growth",
		Duration:     180 * 24 * time.Hour, // 6 months
		Parameters: map[string]interface{}{
			"subsystem_id": "unknown",
			"growth_rate":  0.1, // 10% growth per month
			"team_size":    5,
		},
	},
	"major_dependency_upgrade": {
		Name:         "Major Dependency Upgrade",
		Description:  "Simulates upgrading a major dependency with breaking changes",
		ScenarioType: "dependency_update",
		Duration:     60 * 24 * time.Hour, // 2 months
		Parameters: map[string]interface{}{
			"dependency_name": "unknown",
			"breaking_change": true,
			"effort_hours":    200,
		},
	},
	"large_refactoring": {
		Name:         "Large Refactoring Project",
		Description:  "Simulates effects of a large refactoring initiative",
		ScenarioType: "refactoring",
		Duration:     120 * 24 * time.Hour, // 4 months
		Parameters: map[string]interface{}{
			"subsystem_id":         "unknown",
			"team_size":            5,
			"complexity_reduction": 0.2, // 20% complexity reduction
			"effort_hours":         1000,
		},
	},
}
