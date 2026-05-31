// Package temporal provides temporal repository data management and reconstruction.
// It models repositories as evolving systems with snapshots at different points in time,
// enabling analysis of how repositories change over time.
package temporal

import "time"

// TemporalEvent represents a change event in the repository timeline.
type TemporalEvent struct {
	// Timestamp indicates when the event occurred
	Timestamp time.Time

	// EventType describes what kind of event this is (e.g., "commit", "contributor_added")
	EventType string

	// Contributor is the person associated with this event (if applicable)
	Contributor string

	// Files affected by this event
	Files []string

	// Details contains additional event-specific information
	Details map[string]interface{}
}

// RepositoryMetrics captures computed metrics about a repository at a point in time.
type RepositoryMetrics struct {
	// CommitCount is the total number of commits up to this time
	CommitCount int

	// ContributorCount is the total number of unique contributors up to this time
	ContributorCount int

	// ActiveContributors is the number of contributors active in the recent window
	ActiveContributors int

	// AverageBusFactor is the average bus factor score
	AverageBusFactor float64

	// AverageHealth is the average repository health score
	AverageHealth int

	// FilesChanged is the number of files changed up to this time
	FilesChanged int

	// LinesAdded is the total lines added up to this time
	LinesAdded int

	// LinesRemoved is the total lines removed up to this time
	LinesRemoved int

	// AverageCommitFrequency is the average commits per day
	AverageCommitFrequency float64

	// DependencyCount is the number of dependencies
	DependencyCount int

	// IssuesOpen is the number of open issues
	IssuesOpen int

	// PullRequestsOpen is the number of open pull requests
	PullRequestsOpen int
}

// NewRepositoryMetrics creates a new RepositoryMetrics struct.
func NewRepositoryMetrics() RepositoryMetrics {
	return RepositoryMetrics{}
}

// AggregatedMetrics represents metrics aggregated over a time window.
type AggregatedMetrics struct {
	// StartTime is the beginning of the aggregation window
	StartTime time.Time

	// EndTime is the end of the aggregation window
	EndTime time.Time

	// CommitCount is the number of commits in this window
	CommitCount int

	// Contributors is the set of contributors active in this window
	Contributors []string

	// FilesModified is the number of unique files modified in this window
	FilesModified int

	// LinesAdded is the total lines added in this window
	LinesAdded int

	// LinesRemoved is the total lines removed in this window
	LinesRemoved int
}

// MetricTrend represents how a metric changes over time.
type MetricTrend struct {
	// Metric name
	Name string

	// Values at each time point
	Values []float64

	// Timestamps for each value
	Timestamps []time.Time

	// Trend direction: "increasing", "decreasing", or "stable"
	Direction string

	// Velocity: rate of change per day
	Velocity float64
}
