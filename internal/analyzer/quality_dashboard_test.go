package analyzer

import (
	"testing"
	"time"

	"github.com/agnivo988/Repo-lyzer/internal/github"
)

func TestGenerateQualityDashboard(t *testing.T) {
	// Mock repository data
	repo := &github.Repo{
		FullName:    "test/repo",
		Description: "Test repository",
		Stars:       100,
		Forks:       20,
		OpenIssues:  5,
	}

	// Mock commits
	commits := []github.Commit{
		{SHA: "abc123"},
		{SHA: "def456"},
	}
	// Set commit dates
	commits[0].Commit.Author.Date = time.Now()
	commits[1].Commit.Author.Date = time.Now().AddDate(0, 0, -1)

	// Mock contributors
	contributors := []github.Contributor{
		{Login: "user1", Commits: 50},
		{Login: "user2", Commits: 30},
		{Login: "user3", Commits: 20},
	}

	// Mock security result
	security := &SecurityScanResult{
		SecurityScore:   80,
		TotalCount:      2,
		CriticalCount:   0,
		HighCount:       1,
		MediumCount:     1,
		LowCount:        0,
		ScannedPackages: 10,
		Vulnerabilities: []Vulnerability{},
	}

	dashboard := GenerateQualityDashboard(
		repo,
		commits,
		contributors,
		85, // health score
		3,  // bus factor
		"Mature",
		75, // maturity score
		security,
		nil, // code quality
		nil, // dependencies
	)

	// Test overall score calculation
	if dashboard.OverallScore < 0 || dashboard.OverallScore > 100 {
		t.Errorf("Overall score should be between 0-100, got %d", dashboard.OverallScore)
	}

	// Test risk level determination
	expectedRiskLevels := []string{"Low", "Medium", "High"}
	found := false
	for _, level := range expectedRiskLevels {
		if dashboard.RiskLevel == level {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Risk level should be Low, Medium, or High, got %s", dashboard.RiskLevel)
	}

	// Test quality grade
	expectedGrades := []string{"A+", "A", "B", "C", "D", "F"}
	found = false
	for _, grade := range expectedGrades {
		if dashboard.QualityGrade == grade {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Quality grade should be A+ through F, got %s", dashboard.QualityGrade)
	}

	// Test key metrics
	if dashboard.KeyMetrics.HealthScore != 85 {
		t.Errorf("Expected health score 85, got %d", dashboard.KeyMetrics.HealthScore)
	}

	if dashboard.KeyMetrics.BusFactor != 3 {
		t.Errorf("Expected bus factor 3, got %d", dashboard.KeyMetrics.BusFactor)
	}

	if dashboard.KeyMetrics.ContributorCount != 3 {
		t.Errorf("Expected 3 contributors, got %d", dashboard.KeyMetrics.ContributorCount)
	}

	// Test that recommendations are generated
	if len(dashboard.Recommendations) == 0 {
		t.Error("Expected at least one recommendation")
	}
}

func TestCalculateOverallScore(t *testing.T) {
	tests := []struct {
		name      string
		health    int
		security  int
		maturity  int
		busFactor int
		expectMin int
		expectMax int
	}{
		{"Perfect scores", 100, 100, 100, 10, 95, 100},
		{"Good scores", 80, 80, 80, 5, 70, 85},
		{"Poor scores", 30, 30, 30, 1, 20, 35},
		{"Mixed scores", 90, 50, 70, 3, 60, 80},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := calculateOverallScore(tt.health, tt.security, tt.maturity, tt.busFactor)
			if score < tt.expectMin || score > tt.expectMax {
				t.Errorf("Expected score between %d-%d, got %d", tt.expectMin, tt.expectMax, score)
			}
		})
	}
}

func TestDetermineRiskLevel(t *testing.T) {
	tests := []struct {
		name          string
		overallScore  int
		busFactor     int
		securityScore int
		expected      string
	}{
		{"Low risk", 90, 5, 90, "Low"},
		{"Medium risk", 70, 3, 70, "Medium"},
		{"High risk - low overall", 40, 5, 80, "High"},
		{"High risk - low bus factor", 80, 1, 80, "High"},
		{"High risk - low security", 80, 5, 40, "High"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := determineRiskLevel(tt.overallScore, tt.busFactor, tt.securityScore)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestGetQualityGrade(t *testing.T) {
	tests := []struct {
		score    int
		expected string
	}{
		{95, "A+"},
		{85, "A"},
		{75, "B"},
		{65, "C"},
		{55, "D"},
		{45, "F"},
		{25, "F"},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := getQualityGrade(tt.score)
			if result != tt.expected {
				t.Errorf("Score %d: expected %s, got %s", tt.score, tt.expected, result)
			}
		})
	}
}

func TestIdentifyProblemHotspots(t *testing.T) {
	commits := []github.Commit{
		{SHA: "abc123"},
	}
	commits[0].Commit.Author.Date = time.Now()

	// Test critical security hotspot
	security := &SecurityScanResult{
		SecurityScore: 20,
		CriticalCount: 2,
	}

	hotspots := identifyProblemHotspots(50, 20, 1, commits, security)

	// Should identify security and bus factor hotspots
	if len(hotspots) < 2 {
		t.Errorf("Expected at least 2 hotspots, got %d", len(hotspots))
	}

	// Check for critical severity
	foundCritical := false
	for _, hotspot := range hotspots {
		if hotspot.Severity == "Critical" {
			foundCritical = true
			break
		}
	}
	if !foundCritical {
		t.Error("Expected at least one critical hotspot")
	}
}

func TestGenerateRecommendations(t *testing.T) {
	commits := []github.Commit{
		{SHA: "abc123"},
	}
	commits[0].Commit.Author.Date = time.Now()

	contributors := []github.Contributor{
		{Login: "user1", Commits: 90},
		{Login: "user2", Commits: 10},
	}

	security := &SecurityScanResult{
		SecurityScore: 40,
		CriticalCount: 1,
	}

	recommendations := generateDashboardRecommendations(
		40, // low health
		40, // low security
		1,  // low bus factor
		commits,
		contributors,
		security,
		nil,
	)

	if len(recommendations) == 0 {
		t.Error("Expected recommendations for poor metrics")
	}

	// Should recommend security improvements
	foundSecurityRec := false
	for _, rec := range recommendations {
		if rec == "🔒 Update dependencies to fix security vulnerabilities" {
			foundSecurityRec = true
			break
		}
	}
	if !foundSecurityRec {
		t.Error("Expected security recommendation for low security score")
	}
}

func TestNormalizeBusFactor(t *testing.T) {
	tests := []struct {
		busFactor int
		expected  int
	}{
		{1, 20},   // minimum
		{10, 100}, // maximum
		{5, 55},   // middle range
		{15, 100}, // above maximum
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := normalizeBusFactor(tt.busFactor)
			if result != tt.expected {
				t.Errorf("Bus factor %d: expected %d, got %d", tt.busFactor, tt.expected, result)
			}
		})
	}
}

func TestGetActivityLevel(t *testing.T) {
	tests := []struct {
		commitCount int
		expected    string
	}{
		{600, "Very High"},
		{300, "High"},
		{100, "Medium"},
		{20, "Low"},
		{5, "Very Low"},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			commits := make([]github.Commit, tt.commitCount)
			result := getActivityLevel(commits)
			if result != tt.expected {
				t.Errorf("Commit count %d: expected %s, got %s", tt.commitCount, tt.expected, result)
			}
		})
	}
}
