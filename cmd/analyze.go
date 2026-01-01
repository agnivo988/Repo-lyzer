package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"Repo-lyzer/internal/analyzer"
	"Repo-lyzer/internal/github"
	"Repo-lyzer/internal/output"
)

func RunAnalyze(owner, repo string) error {
	args := []string{owner + "/" + repo}
	analyzeCmd.SetArgs(args)
	return analyzeCmd.Execute()
}


var analyzeCmd = &cobra.Command{
	Use:   "analyze owner/repo",
	Short: "Analyze a GitHub repository",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		parts := strings.Split(args[0], "/")
		if len(parts) != 2 {
			return fmt.Errorf("repository must be in owner/repo format")
		}

		client := github.NewClient()
		repo, err := client.GetRepo(parts[0], parts[1])
		if err != nil {
			return err
		}

		langs, _ := client.GetLanguages(parts[0], parts[1])
		commits, _ := client.GetCommits(parts[0], parts[1], 365)
         
		
		score := analyzer.CalculateHealth(repo, commits)
		activity := analyzer.CommitsPerDay(commits)
		contributors, err := client.GetContributors(parts[0], parts[1])
            if err != nil {
	              return err
                     }

					 busFactor, busRisk := analyzer.BusFactor(contributors)

		maturityScore, maturityLevel :=
			analyzer.RepoMaturityScore(
				repo,
				len(commits),
				len(contributors),
				false,
			)

		summary := analyzer.BuildRecruiterSummary(
			repo.FullName,
			repo.Forks,
		    repo.Stars,
			len(commits),
			len(contributors),
			maturityScore,
			maturityLevel,
			busFactor,
			busRisk,
		)

		output.PrintRepo(repo)
		output.PrintLanguages(langs)
		output.PrintCommitActivity(activity,14)
		output.PrintHealth(score)
		output.PrintRecruiterSummary(summary)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(analyzeCmd)
}
