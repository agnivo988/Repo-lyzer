package analyzer

import (
	"fmt"

	"github.com/agnivo988/Repo-lyzer/internal/analyzer/quality"
	"github.com/agnivo988/Repo-lyzer/internal/github"
)

// CertificateData holds all the data needed to generate a repository certificate
type CertificateData struct {
	RepoName    string
	Owner       string
	Description string
	Stars       int
	Forks       int
	OpenIssues  int
	CreatedAt   string
	UpdatedAt   string

	// Scores
	HealthScore   int
	MaturityScore int
	MaturityLevel string
	BusFactor     int
	BusRisk       string

	// Activity
	CommitsLastYear int
	Contributors    int
	ActivityLevel   string

	// Languages
	PrimaryLanguage string
	LanguageCount   int

	// Calculated overall score
	OverallScore int
	Grade        string
	Uses         []string
}

// CalculateOverallScore computes a weighted overall score from all metrics
func CalculateOverallScore(health, maturity, busFactor, commits, contributors int) int {
	// Activity score based on commits and contributors
	activityScore := 0
	if commits > 1000 {
		activityScore = 100
	} else if commits > 500 {
		activityScore = 80
	} else if commits > 100 {
		activityScore = 60
	} else if commits > 50 {
		activityScore = 40
	} else if commits > 10 {
		activityScore = 20
	}

	if contributors > 50 {
		activityScore += 20
	} else if contributors > 20 {
		activityScore += 15
	} else if contributors > 10 {
		activityScore += 10
	} else if contributors > 5 {
		activityScore += 5
	}

	if activityScore > 100 {
		activityScore = 100
	}

	// Use shared scoring function with security=100 (not calculated in certificate)
	return quality.CalculateOverallScore(health, 100, maturity, busFactor, activityScore)
}

// GetGrade returns a letter grade based on overall score
func GetGrade(score int) string {
	return quality.GetGrade(score)
}

// GetPotentialUses returns a list of potential uses based on the repository's characteristics
func GetPotentialUses(data CertificateData) []string {
	uses := []string{}

	// Based on maturity and activity
	if data.MaturityLevel == "Mature" && data.ActivityLevel == "High" {
		uses = append(uses, "Production-ready software for enterprise use")
		uses = append(uses, "Open-source project for community contributions")
	}

	if data.Stars > 1000 {
		uses = append(uses, "Popular library or framework")
	}

	if data.HealthScore > 80 {
		uses = append(uses, "Well-maintained project suitable for dependencies")
	}

	if data.BusRisk == "Low" {
		uses = append(uses, "Reliable project with diverse contributor base")
	}

	if data.CommitsLastYear > 500 {
		uses = append(uses, "Actively developed project")
	}

	// Default uses if none match
	if len(uses) == 0 {
		uses = append(uses, "Learning resource for developers")
		uses = append(uses, "Personal project with potential for growth")
	}

	return uses
}

// GenerateCertificate analyzes a repository and returns certificate data
func GenerateCertificate(owner, repo string, client *github.Client) (*CertificateData, error) {
	// Fetch all necessary data
	repoInfo, err := client.GetRepo(owner, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to get repo info: %w", err)
	}

	langs, err := client.GetLanguages(owner, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to get languages: %w", err)
	}

	commits, err := client.GetCommits(owner, repo, 365)
	if err != nil {
		return nil, fmt.Errorf("failed to get commits: %w", err)
	}

	contributors, err := client.GetContributorsWithAvatars(owner, repo, 100)
	if err != nil {
		return nil, fmt.Errorf("failed to get contributors: %w", err)
	}

	// Check releases
	hasReleases, _ := client.HasReleases(owner, repo)

	// Calculate scores
	healthScore := CalculateHealth(repoInfo, commits)
	maturityScore, maturityLevel := RepoMaturityScore(repoInfo, len(commits), len(contributors), hasReleases)
	busFactor, busRisk := BusFactor(contributors)

	// Determine primary language
	primaryLang := "Unknown"
	maxSize := 0
	langCount := len(langs)
	for lang, size := range langs {
		if size > maxSize {
			maxSize = size
			primaryLang = lang
		}
	}

	// Activity level
	activityLevel := "Low"
	if len(commits) > 300 {
		activityLevel = "High"
	} else if len(commits) > 100 {
		activityLevel = "Moderate"
	}

	// Calculate overall score
	overallScore := CalculateOverallScore(healthScore, maturityScore, busFactor, len(commits), len(contributors))
	grade := GetGrade(overallScore)

	// Create certificate data
	cert := &CertificateData{
		RepoName:        repoInfo.Name,
		Owner:           owner,
		Description:     repoInfo.Description,
		Stars:           repoInfo.Stars,
		Forks:           repoInfo.Forks,
		OpenIssues:      repoInfo.OpenIssues,
		CreatedAt:       repoInfo.CreatedAt.Format("2006-01-02"),
		UpdatedAt:       repoInfo.UpdatedAt.Format("2006-01-02"),
		HealthScore:     healthScore,
		MaturityScore:   maturityScore,
		MaturityLevel:   maturityLevel,
		BusFactor:       busFactor,
		BusRisk:         busRisk,
		CommitsLastYear: len(commits),
		Contributors:    len(contributors),
		ActivityLevel:   activityLevel,
		PrimaryLanguage: primaryLang,
		LanguageCount:   langCount,
		OverallScore:    overallScore,
		Grade:           grade,
	}

	// Generate potential uses
	cert.Uses = GetPotentialUses(*cert)

	return cert, nil
}
