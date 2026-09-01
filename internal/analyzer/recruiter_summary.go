package analyzer

type RecruiterSummary struct {
	RepoName        string
	Stars           int
	Forks           int
	CommitsLastYear int
	Contributors    int

	MaturityScore int
	MaturityLevel string

	BusFactor int
	BusRisk   string

	IssueHealth   string
	PRHealth      string
	ActivityLevel string
}

func BuildRecruiterSummary(
	repoName string,
	stars, forks int,
	commits, contributors int,
	maturityScore int,
	maturityLevel string,
	busFactor int,
	busRisk string,
	openIssues int,
	prMergeRate float64,
) RecruiterSummary {

	activity := "Low"
	if commits > 300 {
		activity = "High"
	} else if commits > 100 {
		activity = "Moderate"
	}

	issueHealth := "Critical"
	if openIssues < 20 {
		issueHealth = "Healthy"
	} else if openIssues < 50 {
		issueHealth = "Backlogged"
	}

	prHealth := "Stalled"
	if prMergeRate > 60 {
		prHealth = "Active"
	} else if prMergeRate >= 30 {
		prHealth = "Slow"
	}

	return RecruiterSummary{
		RepoName:        repoName,
		Stars:           stars,
		Forks:           forks,
		CommitsLastYear: commits,
		Contributors:    contributors,
		MaturityScore:   maturityScore,
		MaturityLevel:   maturityLevel,
		BusFactor:       busFactor,
		BusRisk:         busRisk,
		IssueHealth:     issueHealth,
		PRHealth:        prHealth,
		ActivityLevel:   activity,
	}
}
