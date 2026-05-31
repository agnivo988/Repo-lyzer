package analyzer

import (
	"encoding/base64"
	"testing"
)

func TestCalculateComplexityFromGitHubContentDecodesBeforeScanning(t *testing.T) {
	source := `function risky(value) {
  if (value) {
    for (let i = 0; i < 3; i++) {
      if (i && value.ok) {
        return value.ready ? 1 : 0
      }
    }
  }
}`
	encoded := base64.StdEncoding.EncodeToString([]byte(source))
	wrapped := encoded[:24] + "\n" + encoded[24:48] + "\r\n" + encoded[48:]

	got, err := calculateComplexityFromGitHubContent(wrapped, "risk.js")
	if err != nil {
		t.Fatalf("expected encoded GitHub content to decode: %v", err)
	}

	want := calculateComplexity(source, "risk.js")
	if got != want {
		t.Fatalf("expected decoded source complexity %d, got %d", want, got)
	}

	encodedOnly := calculateComplexity(encoded, "risk.js")
	if got <= encodedOnly {
		t.Fatalf("expected decoded complexity %d to exceed encoded-text complexity %d", got, encodedOnly)
	}
}

func TestCalculateComplexityFromGitHubContentRejectsInvalidBase64(t *testing.T) {
	if _, err := calculateComplexityFromGitHubContent("not valid base64", "broken.go"); err == nil {
		t.Fatal("expected invalid base64 content to return an error")
	}
}
