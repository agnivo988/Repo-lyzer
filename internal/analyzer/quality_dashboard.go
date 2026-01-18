package analyzer

import (
	"fmt"
	"sort"

	"github.com/agnivo988/Repo-lyzer/internal/github"
)

// QualityDashboard represents the high-level summary
type QualityDashboard struct {
	OverallScore    int              `json:"overall_score"`
	RiskLevel       string           `json:"risk_level"`
	QualityGrade    string           `json:"quality_grade"`
	ProblemHotspots []ProblemHotspot `json:"problem_hotspots"`
	Recommendations []string         `json:"recommendations"`
	KeyMetrics      DashboardMetrics `json:"key_metrics"`
}

// ProblemHotspot identifies high-risk areas
type ProblemHotspot struct {
	Area        string `json:"area"`
	Severity    string `json:"severity"`
	Description string `json:"description"`
	Impact      string `json:"impact"`
}

// DashboardMetrics contains key metrics for quick assessment
type DashboardMetrics struct {
	HealthScore      int    `json:"health_score"`
	SecurityScore    int    `json:"security_score"`
	MaturityLevel    string `json:"maturity_level"`
	BusFactor        int    `json:"bus_factor"`
	ActivityLevel    string `json:"activity_level"`
	ContributorCount int    `json:"contributor_count"`
}

// GenerateQualityDashboard creates a comprehensive quality and risk summary
func GenerateQualityDashboard(
	repo *github.Repo,
	commits []github.Commit,
	contributors []github.Contributor,
	healthScore int,
	busFactor int,
	maturityLevel string,
	maturityScore int,
	security *SecurityScanResult,
	codeQuality *CodeQualityMetrics,
	dependencies *DependencyAnalysis,
) *QualityDashboard {

	dashboard := &QualityDashboard{
		ProblemHotspots: []ProblemHotspot{},
		Recommendations: []string{},
	}

	// Calculate overall score (weighted average)
	securityScore := 100
	if security != nil {
		securityScore = security.SecurityScore
	}

	dashboard.OverallScore = calculateOverallScore(healthScore, securityScore, maturityScore, busFactor)
	dashboard.RiskLevel = determineRiskLevel(dashboard.OverallScore, busFactor, securityScore)
	dashboard.QualityGrade = getQualityGrade(dashboard.OverallScore)

	// Populate key metrics
	dashboard.KeyMetrics = DashboardMetrics{
		HealthScore:      healthScore,
		SecurityScore:    securityScore,
		MaturityLevel:    maturityLevel,
		BusFactor:        busFactor,
		ActivityLevel:    getActivityLevel(commits),
		ContributorCount: len(contributors),
	}

	// Identify problem hotspots
	dashboard.ProblemHotspots = identifyProblemHotspots(
		healthScore, securityScore, busFactor, commits, security,
	)

	// Generate actionable recommendations
	dashboard.Recommendations = generateDashboardRecommendations(
		healthScore, securityScore, busFactor, commits, contributors, security, dependencies,
	)

	return dashboard
}

func calculateOverallScore(health, security, maturity, busFactor int) int {
	// Weighted scoring: Health(30%), Security(30%), Maturity(25%), Bus Factor(15%)
	busFactorScore := normalizeBusFactor(busFactor)

	score := int(
		float64(health)*0.30 +
			float64(security)*0.30 +
			float64(maturity)*0.25 +
			float64(busFactorScore)*0.15,
	)

	if score > 100 {
		score = 100
	}
	if score < 0 {
		score = 0
	}

	return score
}

func normalizeBusFactor(busFactor int) int {
	// Convert bus factor to 0-100 scale (higher bus factor = better score)
	if busFactor >= 10 {
		return 100
	}
	if busFactor <= 1 {
		return 20
	}
	// Linear scaling between 1 and 10
	return 20 + ((busFactor-1)*80)/9
}

func determineRiskLevel(overallScore, busFactor, securityScore int) string {
	if overallScore >= 80 && busFactor >= 3 && securityScore >= 80 {
		return "Low"
	}
	if overallScore >= 60 && busFactor >= 2 && securityScore >= 60 {
		return "Medium"
	}
	return "High"
}

func getQualityGrade(score int) string {
	if score >= 90 {
		return "A+"
	}
	if score >= 80 {
		return "A"
	}
	if score >= 70 {
		return "B"
	}
	if score >= 60 {
		return "C"
	}
	if score >= 50 {
		return "D"
	}
	return "F"
}

func getActivityLevel(commits []github.Commit) string {
	commitCount := len(commits)
	if commitCount >= 500 {
		return "Very High"
	}
	if commitCount >= 200 {
		return "High"
	}
	if commitCount >= 50 {
		return "Medium"
	}
	if commitCount >= 10 {
		return "Low"
	}
	return "Very Low"
}

func identifyProblemHotspots(
	health, security, busFactor int,
	commits []github.Commit,
	securityResult *SecurityScanResult,
) []ProblemHotspot {

	var hotspots []ProblemHotspot

	// Security hotspots
	if security < 60 {
		severity := "High"
		if security < 30 {
			severity = "Critical"
		}
		hotspots = append(hotspots, ProblemHotspot{
			Area:        "Security",
			Severity:    severity,
			Description: fmt.Sprintf("Security score is %d/100", security),
			Impact:      "Potential vulnerabilities may expose the project to security risks",
		})
	}

	// Bus factor hotspots
	if busFactor <= 1 {
		hotspots = append(hotspots, ProblemHotspot{
			Area:        "Bus Factor",
			Severity:    "Critical",
			Description: "Single contributor dependency",
			Impact:      "Project is at risk if the main contributor becomes unavailable",
		})
	} else if busFactor <= 2 {
		hotspots = append(hotspots, ProblemHotspot{
			Area:        "Bus Factor",
			Severity:    "High",
			Description: "Very low contributor diversity",
			Impact:      "Limited knowledge distribution across the team",
		})
	}

	// Activity hotspots
	recentCommits := countRecentCommits(commits)
	if recentCommits == 0 {
		hotspots = append(hotspots, ProblemHotspot{
			Area:        "Activity",
			Severity:    "High",
			Description: "No commits in the last 90 days",
			Impact:      "Project may be abandoned or inactive",
		})
	}

	// Health hotspots
	if health < 40 {
		hotspots = append(hotspots, ProblemHotspot{
			Area:        "Repository Health",
			Severity:    "High",
			Description: fmt.Sprintf("Health score is %d/100", health),
			Impact:      "Poor repository maintenance and documentation",
		})
	}

	// Dependency hotspots
	if securityResult != nil && securityResult.CriticalCount > 0 {
		hotspots = append(hotspots, ProblemHotspot{
			Area:        "Dependencies",
			Severity:    "Critical",
			Description: fmt.Sprintf("%d critical vulnerabilities found", securityResult.CriticalCount),
			Impact:      "Critical security vulnerabilities in dependencies",
		})
	}

	// Sort by severity (Critical > High > Medium > Low)
	sort.Slice(hotspots, func(i, j int) bool {
		severityOrder := map[string]int{"Critical": 4, "High": 3, "Medium": 2, "Low": 1}
		return severityOrder[hotspots[i].Severity] > severityOrder[hotspots[j].Severity]
	})

	return hotspots
}

func generateDashboardRecommendations(
	health, security, busFactor int,
	commits []github.Commit,
	contributors []github.Contributor,
	securityResult *SecurityScanResult,
	deps *DependencyAnalysis,
) []string {

	var recommendations []string

	// Security recommendations
	if security < 70 {
		recommendations = append(recommendations, "🔒 Update dependencies to fix security vulnerabilities")
		if securityResult != nil && securityResult.CriticalCount > 0 {
			recommendations = append(recommendations, "🚨 Immediately address critical security vulnerabilities")
		}
	}

	// Bus factor recommendations
	if busFactor <= 2 {
		recommendations = append(recommendations, "👥 Encourage more contributors to reduce bus factor risk")
		recommendations = append(recommendations, "📚 Improve documentation to enable easier onboarding")
	}

	// Activity recommendations
	recentCommits := countRecentCommits(commits)
	if recentCommits == 0 {
		recommendations = append(recommendations, "🔄 Resume development activity or archive if project is complete")
	} else if recentCommits < 10 {
		recommendations = append(recommendations, "📈 Increase development activity and regular maintenance")
	}

	// Health recommendations
	if health < 60 {
		recommendations = append(recommendations, "📝 Add comprehensive README and project description")
		recommendations = append(recommendations, "🐛 Address open issues to improve repository health")
	}

	// Dependency recommendations
	if deps != nil && !deps.HasLockFile {
		recommendations = append(recommendations, "🔒 Add dependency lock files for reproducible builds")
	}

	// General recommendations
	if len(contributors) < 3 {
		recommendations = append(recommendations, "🌟 Promote the project to attract more contributors")
	}

	// Limit to top 5 recommendations
	if len(recommendations) > 5 {
		recommendations = recommendations[:5]
	}

	return recommendations
}

func countRecentCommits(commits []github.Commit) int {
	// This is a simplified version - in practice you'd check commit dates
	// For now, we'll assume all commits in the slice are recent
	return len(commits)
}

// GetRiskLevelColor returns color styling for risk level
func (d *QualityDashboard) GetRiskLevelColor() string {
	switch d.RiskLevel {
	case "Low":
		return "🟢"
	case "Medium":
		return "🟡"
	case "High":
		return "🔴"
	default:
		return "⚪"
	}
}

// GetGradeColor returns color styling for quality grade
func (d *QualityDashboard) GetGradeColor() string {
	switch d.QualityGrade {
	case "A+", "A":
		return "🟢"
	case "B":
		return "🟡"
	case "C", "D":
		return "🟠"
	case "F":
		return "🔴"
	default:
		return "⚪"
	}
}

// FormatSummary returns a formatted summary string
func (d *QualityDashboard) FormatSummary() string {
	return fmt.Sprintf(
		"%s Overall Score: %d/100 (Grade: %s)\n%s Risk Level: %s",
		d.GetGradeColor(), d.OverallScore, d.QualityGrade,
		d.GetRiskLevelColor(), d.RiskLevel,
	)
}
