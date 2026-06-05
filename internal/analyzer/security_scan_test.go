package analyzer

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestQueryOSV_Errors(t *testing.T) {
	tests := []struct {
		name           string
		handler        http.HandlerFunc
		expectedErrMsg string
	}{
		{
			name: "HTTP error status code",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			expectedErrMsg: "OSV API returned status: 500 Internal Server Error",
		},
		{
			name: "Invalid JSON response",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{invalid json`))
			},
			expectedErrMsg: "failed to decode OSV response",
		},
		{
			name: "Successful response",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"vulns": [{"id": "TEST-VULN"}]}`))
			},
			expectedErrMsg: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()

			ctx := context.Background()
			client := &http.Client{}
			vulns, err := queryOSV(ctx, client, "test-pkg", "1.0.0", "npm", server.URL)

			if tt.expectedErrMsg != "" {
				if err == nil {
					t.Errorf("expected error containing %q, got nil", tt.expectedErrMsg)
				} else if !strings.Contains(err.Error(), tt.expectedErrMsg) {
					t.Errorf("expected error containing %q, got %q", tt.expectedErrMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if len(vulns) != 1 || vulns[0].ID != "TEST-VULN" {
					t.Errorf("unexpected vulns: %+v", vulns)
				}
			}
		})
	}
}

func TestScanDependencies_ErrorPropagation(t *testing.T) {
	// Mock OSV API to return an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	deps := &DependencyAnalysis{
		Files: []DependencyFile{
			{
				FileType: "npm",
				Dependencies: []Dependency{
					{Name: "test-pkg", Version: "1.0.0"},
				},
			},
		},
	}

	ctx := context.Background()
	result, err := ScanDependencies(ctx, deps, server.URL)
	if err == nil {
		t.Errorf("ScanDependencies should have failed but got nil error")
	} else if !strings.Contains(err.Error(), "failed to query OSV") {
		t.Errorf("expected error containing %q, got %q", "failed to query OSV", err.Error())
	}
	
	if result != nil {
		t.Errorf("result should be nil on error")
	}
}

func TestScanDependencies_Cancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate slow response
		<-r.Context().Done()
	}))
	defer server.Close()

	deps := &DependencyAnalysis{
		Files: []DependencyFile{
			{
				FileType: "npm",
				Dependencies: []Dependency{
					{Name: "test-pkg", Version: "1.0.0"},
				},
			},
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	result, err := ScanDependencies(ctx, deps, server.URL)
	if err == nil {
		t.Errorf("expected error from cancelled context, got nil")
	} else if !strings.Contains(err.Error(), context.Canceled.Error()) {
		t.Errorf("expected error containing %q, got %q", context.Canceled.Error(), err.Error())
	}
	if result != nil {
		t.Errorf("result should be nil on cancellation")
	}
}
