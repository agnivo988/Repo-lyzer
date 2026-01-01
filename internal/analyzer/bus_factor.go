package analyzer

import "Repo-lyzer/internal/github"

func BusFactor(contributors []github.Contributor) (int, string) {
	if len(contributors) == 0 {
		return 0, "Unknown"
	}

	total := 0
	for _, c := range contributors {
		total += c.Commits
	}

	top := contributors[0].Commits
	ratio := float64(top) / float64(total)

	switch {
	case ratio > 0.7:
		return 1, "High Risk"
	case ratio > 0.4:
		return 2, "Medium Risk"
	default:
		return 3, "Low Risk"
	}
}
