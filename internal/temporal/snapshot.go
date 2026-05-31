package temporal

import (
	"fmt"
	"time"

	"github.com/agnivo988/Repo-lyzer/internal/graph"
)

// Snapshot represents the state of a repository at a specific point in time.
type Snapshot struct {
	// Timestamp is when this snapshot was taken
	Timestamp time.Time

	// CommitHash is the commit hash at the time of this snapshot
	CommitHash string

	// Graph is the repository ecosystem graph at this snapshot
	Graph *graph.Graph

	// Metrics are computed metrics for this snapshot
	Metrics RepositoryMetrics

	// Contributors is a list of active contributors at this snapshot
	Contributors []string

	// Files is the total number of files at this snapshot
	Files int
}

// NewSnapshot creates a new snapshot with the given timestamp and graph.
func NewSnapshot(timestamp time.Time, g *graph.Graph) *Snapshot {
	return &Snapshot{
		Timestamp:    timestamp,
		Graph:        g,
		Metrics:      NewRepositoryMetrics(),
		Contributors: []string{},
		Files:        0,
	}
}

// Timeline represents the evolution of a repository over time.
// It is a sequence of snapshots showing how the repository changed.
type Timeline struct {
	// RepoName is the name of the repository
	RepoName string

	// Owner is the repository owner/organization
	Owner string

	// Snapshots is a chronologically ordered list of repository states
	Snapshots []*Snapshot

	// StartTime is the time of the first snapshot
	StartTime time.Time

	// EndTime is the time of the last snapshot
	EndTime time.Time

	// IntervalDays is the number of days between consecutive snapshots
	IntervalDays int

	// Events is a chronologically ordered list of all temporal events
	Events []TemporalEvent
}

// NewTimeline creates a new empty timeline.
func NewTimeline(owner, repoName string) *Timeline {
	return &Timeline{
		RepoName:     repoName,
		Owner:        owner,
		Snapshots:    make([]*Snapshot, 0),
		Events:       make([]TemporalEvent, 0),
		IntervalDays: 1,
	}
}

// AddSnapshot adds a snapshot to the timeline.
// Snapshots must be added in chronological order.
func (t *Timeline) AddSnapshot(snapshot *Snapshot) error {
	if snapshot == nil {
		return fmt.Errorf("cannot add nil snapshot")
	}

	if len(t.Snapshots) > 0 {
		lastSnapshot := t.Snapshots[len(t.Snapshots)-1]
		if snapshot.Timestamp.Before(lastSnapshot.Timestamp) {
			return fmt.Errorf("snapshots must be added in chronological order")
		}
	}

	t.Snapshots = append(t.Snapshots, snapshot)

	if len(t.Snapshots) == 1 {
		t.StartTime = snapshot.Timestamp
	}
	t.EndTime = snapshot.Timestamp

	return nil
}

// AddEvent adds a temporal event to the timeline.
func (t *Timeline) AddEvent(event TemporalEvent) error {
	if event.Timestamp.IsZero() {
		return fmt.Errorf("event timestamp cannot be zero")
	}

	if len(t.Events) > 0 {
		lastEvent := t.Events[len(t.Events)-1]
		if event.Timestamp.Before(lastEvent.Timestamp) {
			return fmt.Errorf("events must be added in chronological order")
		}
	}

	t.Events = append(t.Events, event)
	return nil
}

// GetSnapshot retrieves the snapshot closest to the given timestamp.
// If exact is true, returns an error if no exact match is found.
func (t *Timeline) GetSnapshot(timestamp time.Time, exact bool) (*Snapshot, error) {
	if len(t.Snapshots) == 0 {
		return nil, fmt.Errorf("timeline has no snapshots")
	}

	if exact {
		for _, snapshot := range t.Snapshots {
			if snapshot.Timestamp.Equal(timestamp) {
				return snapshot, nil
			}
		}
		return nil, fmt.Errorf("no exact snapshot found at %v", timestamp)
	}

	// Find closest snapshot
	closest := t.Snapshots[0]
	minDiff := timestamp.Sub(closest.Timestamp)
	if minDiff < 0 {
		minDiff = -minDiff
	}

	for _, snapshot := range t.Snapshots[1:] {
		diff := timestamp.Sub(snapshot.Timestamp)
		if diff < 0 {
			diff = -diff
		}
		if diff < minDiff {
			minDiff = diff
			closest = snapshot
		}
	}

	return closest, nil
}

// SnapshotsBetween returns all snapshots within the given time range (inclusive).
func (t *Timeline) SnapshotsBetween(startTime, endTime time.Time) []*Snapshot {
	var result []*Snapshot
	for _, snapshot := range t.Snapshots {
		if (snapshot.Timestamp.Equal(startTime) || snapshot.Timestamp.After(startTime)) &&
			(snapshot.Timestamp.Equal(endTime) || snapshot.Timestamp.Before(endTime)) {
			result = append(result, snapshot)
		}
	}
	return result
}

// EventsBetween returns all events within the given time range (inclusive).
func (t *Timeline) EventsBetween(startTime, endTime time.Time) []TemporalEvent {
	var result []TemporalEvent
	for _, event := range t.Events {
		if (event.Timestamp.Equal(startTime) || event.Timestamp.After(startTime)) &&
			(event.Timestamp.Equal(endTime) || event.Timestamp.Before(endTime)) {
			result = append(result, event)
		}
	}
	return result
}

// SnapshotCount returns the number of snapshots in the timeline.
func (t *Timeline) SnapshotCount() int {
	return len(t.Snapshots)
}

// EventCount returns the number of events in the timeline.
func (t *Timeline) EventCount() int {
	return len(t.Events)
}

// Duration returns the time span covered by this timeline.
func (t *Timeline) Duration() time.Duration {
	if len(t.Snapshots) == 0 {
		return 0
	}
	return t.EndTime.Sub(t.StartTime)
}

// IsEmpty returns true if the timeline has no snapshots or events.
func (t *Timeline) IsEmpty() bool {
	return len(t.Snapshots) == 0 && len(t.Events) == 0
}

// LatestSnapshot returns the most recent snapshot, or nil if none exist.
func (t *Timeline) LatestSnapshot() *Snapshot {
	if len(t.Snapshots) == 0 {
		return nil
	}
	return t.Snapshots[len(t.Snapshots)-1]
}

// EarliestSnapshot returns the earliest snapshot, or nil if none exist.
func (t *Timeline) EarliestSnapshot() *Snapshot {
	if len(t.Snapshots) == 0 {
		return nil
	}
	return t.Snapshots[0]
}

// SnapshotAtIndex returns the snapshot at the given index.
func (t *Timeline) SnapshotAtIndex(index int) (*Snapshot, error) {
	if index < 0 || index >= len(t.Snapshots) {
		return nil, fmt.Errorf("index %d out of range [0, %d)", index, len(t.Snapshots))
	}
	return t.Snapshots[index], nil
}

// WindowedSnapshots returns snapshots grouped by time windows.
// For example, if windowDays=7, returns snapshots grouped by week.
func (t *Timeline) WindowedSnapshots(windowDays int) [][]*Snapshot {
	if len(t.Snapshots) == 0 {
		return [][]*Snapshot{}
	}

	if windowDays <= 0 {
		windowDays = 1
	}

	var windows [][]*Snapshot
	var currentWindow []*Snapshot
	var currentWindowStart time.Time

	for i, snapshot := range t.Snapshots {
		if i == 0 {
			currentWindowStart = snapshot.Timestamp
		}

		windowEnd := currentWindowStart.AddDate(0, 0, windowDays)
		if snapshot.Timestamp.Before(windowEnd) {
			currentWindow = append(currentWindow, snapshot)
		} else {
			if len(currentWindow) > 0 {
				windows = append(windows, currentWindow)
			}
			currentWindow = []*Snapshot{snapshot}
			currentWindowStart = snapshot.Timestamp
		}
	}

	if len(currentWindow) > 0 {
		windows = append(windows, currentWindow)
	}

	return windows
}
