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
) RecruiterSummary {

	activity := "Low"
	if commits > 300 {
		activity = "High"
	} else if commits > 100 {
		activity = "Moderate"
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
		ActivityLevel:   activity,
	}
}
