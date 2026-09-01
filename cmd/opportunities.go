package cmd

import (
    "context"
    "fmt"
    "github.com/spf13/cobra"
    "github.com/agnivo988/Repo-lyzer/internal/opportunity"
)

var opportunitiesCmd = &cobra.Command{
    Use:   "opportunities owner/repo",
    Short: "Analyze contribution opportunities for a repository",
    Args:  cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        ownerRepo := args[0]
        // Validate format owner/repo
        var owner, repo string
        parts := strings.SplitN(ownerRepo, "/", 2)
        if len(parts) != 2 {
            return fmt.Errorf("invalid repository format, expected owner/repo")
        }
        owner, repo = parts[0], parts[1]
        ctx := context.Background()
        eng := opportunity.NewEngine()
        report, err := eng.Run(ctx, owner, repo)
        if err != nil {
            return err
        }
        fmt.Printf("Contribution Opportunity Score for %s/%s: %d\n", owner, repo, report.Score)
        return nil
    },
}

func init() {
    rootCmd.AddCommand(opportunitiesCmd)
}
