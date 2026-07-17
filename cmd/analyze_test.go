package cmd

import (
	"errors"
	"testing"
)

func TestIsRateLimitError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{name: "primary rate limit", err: errors.New("rate limit exceeded"), want: true},
		{name: "secondary rate limit", err: errors.New("GitHub secondary rate-limit exceeded"), want: true},
		{name: "unrelated error", err: errors.New("repository not found"), want: false},
		{name: "nil", err: nil, want: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := isRateLimitError(test.err); got != test.want {
				t.Fatalf("isRateLimitError(%v) = %v, want %v", test.err, got, test.want)
			}
		})
	}
}
