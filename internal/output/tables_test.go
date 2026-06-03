package output

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/agnivo988/Repo-lyzer/internal/github"
)

func TestPrintRepo(t *testing.T) {
	repo := &github.Repo{
		FullName:   "owner/example",
		Stars:      120,
		Forks:      18,
		OpenIssues: 7,
	}

	// Regression guard: ensure table rendering path executes without panicking
	// and preserves the expected repo header and row values.
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe() failed: %v", err)
	}
	var out bytes.Buffer

	os.Stdout = w
	defer func() {
		w.Close()
		os.Stdout = oldStdout

		if r := recover(); r != nil {
			t.Fatalf("PrintRepo panicked: %v", r)
		}
	}()

	PrintRepo(repo)

	w.Close()
	if _, err := io.Copy(&out, r); err != nil {
		t.Fatalf("failed to read captured output: %v", err)
	}

	output := out.String()
	outputUpper := strings.ToUpper(output)
	for _, expected := range []string{
		"REPOSITORY",
		"STARS",
		"FORKS",
		"OPEN ISSUES",
		"OWNER/EXAMPLE",
		"120",
		"18",
		"7",
	} {
		if !strings.Contains(outputUpper, expected) {
			t.Fatalf("expected output to contain %q, got:\n%s", expected, output)
		}
	}
}
