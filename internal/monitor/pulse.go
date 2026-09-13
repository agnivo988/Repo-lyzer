package monitor

import (
	"time"

	"github.com/agnivo988/Repo-lyzer/internal/github"
)

// PulseStatus represents the current live activity level of the repository
type PulseStatus string

const (
	PulseSpiking PulseStatus = "SPIKING" // High activity in last hour
	PulseActive  PulseStatus = "ACTIVE"  // Normal activity in last 24h
	PulseQuiet   PulseStatus = "QUIET"   // Low activity
	PulseDormant PulseStatus = "DORMANT" // No activity recently
)

// CalculateLivePulse computes the repository's pulse based on recent events
func CalculateLivePulse(events []github.Event) PulseStatus {
	if len(events) == 0 {
		return PulseDormant
	}

	now := time.Now()

	eventsInLastHour := 0
	eventsInLastDay := 0

	for _, e := range events {
		age := now.Sub(e.CreatedAt)

		if age <= time.Hour {
			eventsInLastHour++
		}
		if age <= 24*time.Hour {
			eventsInLastDay++
		}
	}

	if eventsInLastHour > 10 {
		return PulseSpiking
	}
	if eventsInLastDay > 5 {
		return PulseActive
	}
	if eventsInLastDay > 0 {
		return PulseQuiet
	}

	return PulseDormant
}
