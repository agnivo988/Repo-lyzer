package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/agnivo988/Repo-lyzer/internal/analyzer"
	"github.com/agnivo988/Repo-lyzer/internal/github"
	"github.com/agnivo988/Repo-lyzer/internal/output"
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

		langs, err := client.GetLanguages(parts[0], parts[1])
		if err != nil {
			return fmt.Errorf("failed to get languages: %w", err)
		}

		commits, err := client.GetCommits(parts[0], parts[1], 365)
		if err != nil {
			return fmt.Errorf("failed to get commits: %w", err)
		}

		_, err = client.GetFileTree(parts[0], parts[1], repo.DefaultBranch)
		if err != nil {
			return fmt.Errorf("failed to get file tree: %w", err)
		}
         
		
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
		output.PrintGitHubAPIStatus(client)
		output.PrintRecruiterSummary(summary)

		return nil
	},
}
