package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/agnivo988/Repo-lyzer/internal/config"
	"github.com/agnivo988/Repo-lyzer/internal/gitpush"
	"github.com/agnivo988/Repo-lyzer/internal/pathutil"
	"github.com/agnivo988/Repo-lyzer/internal/progress"
	"github.com/spf13/cobra"
)

var (
	pushToken  string
	pushMsg    string
	pushBranch string
)

var pushCmd = &cobra.Command{
	Use:   "push <local-folder> <owner/repo>",
	Short: "Push a local folder or repository to GitHub",
	Long: `Initialize, commit, and push a local folder or existing repository 
to a remote GitHub repository.

If the local folder is not already a Git repository, it will be initialized 
automatically. Any untracked/modified files will be staged and committed 
before pushing.

Authentication:
Uses your configured GitHub Personal Access Token. Alternatively, you can 
provide a token via the --token (-t) flag or set the GITHUB_TOKEN environment 
variable.

Examples:
  # Push local folder to remote repo
  repo-lyzer push /path/to/project octocat/my-new-repo

  # With custom commit message and branch
  repo-lyzer push ./my-app octocat/my-app -m "Deploy v1.0" -b production`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		localPath := args[0]
		repoName := args[1]

		parts := strings.Split(repoName, "/")
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return fmt.Errorf("invalid repository: must be in owner/repo format")
		}
		owner := parts[0]
		repo := parts[1]

		// Resolve and clean local folder path
		resolvedPath, err := pathutil.CleanAndResolvePath(localPath)
		if err != nil {
			return err
		}

		// Check if local folder exists
		if stat, err := os.Stat(resolvedPath); os.IsNotExist(err) {
			return fmt.Errorf("local folder does not exist: %s", resolvedPath)
		} else if !stat.IsDir() {
			return fmt.Errorf("path is not a directory: %s", resolvedPath)
		}
		localPath = resolvedPath

		spinner := progress.NewSpinner()
		spinner.Start(fmt.Sprintf("Preparing and pushing to GitHub (%s)...", pushBranch))

		settings, err := config.LoadSettings()
		token := pushToken
		if token == "" && err == nil {
			token = settings.GitHubToken
		}
		if token == "" {
			token = os.Getenv("GITHUB_TOKEN")
		}

		errPush := gitpush.PushRepo(context.Background(), gitpush.PushOptions{
			LocalPath: localPath,
			RepoOwner: owner,
			RepoName:  repo,
			CommitMsg: pushMsg,
			Token:     token,
			Branch:    pushBranch,
		})

		spinner.Stop()
		if errPush != nil {
			return errPush
		}

		fmt.Printf("✓ Pushed code successfully to %s/%s branch %s!\n", owner, repo, pushBranch)
		return nil
	},
}

func init() {
	pushCmd.Flags().StringVarP(&pushToken, "token", "t", "", "GitHub Personal Access Token")
	pushCmd.Flags().StringVarP(&pushMsg, "message", "m", "Push via Repo-lyzer", "Commit message")
	pushCmd.Flags().StringVarP(&pushBranch, "branch", "b", "main", "Target branch")

	rootCmd.AddCommand(pushCmd)
}
