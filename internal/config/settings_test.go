package config

import (
	"testing"
)

func TestGetScheduledJobByID(t *testing.T) {
	settings := &AppSettings{
		ScheduledJobs: []ScheduledJob{
			{ID: "job-1", Owner: "owner1", Repo: "repo1"},
			{ID: "job-2", Owner: "owner2", Repo: "repo2"},
			{ID: "job-3", Owner: "owner3", Repo: "repo3"},
		},
	}

	tests := []struct {
		jobID string
		want  string
	}{
		{"job-1", "owner1"},
		{"job-2", "owner2"},
		{"job-3", "owner3"},
		{"non-existent", ""},
	}

	for _, tt := range tests {
		job := settings.GetScheduledJobByID(tt.jobID)
		if tt.want == "" {
			if job != nil {
				t.Errorf("GetScheduledJobByID(%s) = %v, want nil", tt.jobID, job)
			}
		} else {
			if job == nil {
				t.Errorf("GetScheduledJobByID(%s) = nil, want %s", tt.jobID, tt.want)
			} else if job.Owner != tt.want {
				t.Errorf("GetScheduledJobByID(%s) owner = %s, want %s", tt.jobID, job.Owner, tt.want)
			}
		}
	}
}
