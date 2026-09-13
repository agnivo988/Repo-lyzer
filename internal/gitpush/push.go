package gitpush

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// PushOptions holds options for the Git push operation.
type PushOptions struct {
	LocalPath string
	RepoOwner string
	RepoName  string
	CommitMsg string
	Token     string
	Branch    string
	Username  string
}

// PushRepo encapsulates staging, committing, remote URL configuration, and pushing to GitHub.
func PushRepo(parentCtx context.Context, opts PushOptions) error {
	ctx, cancel := context.WithTimeout(parentCtx, 5*time.Minute)
	defer cancel()

	// 1. Initialize git if not already inside a git worktree
	cmdIsRepo := exec.CommandContext(ctx, "git", "rev-parse", "--is-inside-work-tree")
	cmdIsRepo.Dir = opts.LocalPath
	out, err := cmdIsRepo.Output()
	if err != nil || strings.TrimSpace(string(out)) != "true" {
		cmdInit := exec.CommandContext(ctx, "git", "init")
		cmdInit.Dir = opts.LocalPath
		if err := cmdInit.Run(); err != nil {
			if ctx.Err() != nil {
				return fmt.Errorf("failed to initialize git: %w: %v", ctx.Err(), err)
			}
			return fmt.Errorf("failed to initialize git: %w", err)
		}
	}

	// 2. Check current branch if none provided, or default to main
	branch := opts.Branch
	if branch == "" {
		cmdBranch := exec.CommandContext(ctx, "git", "symbolic-ref", "--quiet", "--short", "HEAD")
		cmdBranch.Dir = opts.LocalPath
		branchBytes, err := cmdBranch.Output()
		if err == nil {
			b := strings.TrimSpace(string(branchBytes))
			if b != "" && b != "HEAD" {
				branch = b
			}
		}

		if branch == "" {
			cmdBranch = exec.CommandContext(ctx, "git", "rev-parse", "--abbrev-ref", "HEAD")
			cmdBranch.Dir = opts.LocalPath
			branchBytes, err = cmdBranch.Output()
			if err == nil {
				b := strings.TrimSpace(string(branchBytes))
				if b != "" && b != "HEAD" {
					branch = b
				}
			}
		}

		if branch == "" {
			branch = "main"
		}
	}

	// 3. Stage and commit files
	cmdStatus := exec.CommandContext(ctx, "git", "status", "--porcelain")
	cmdStatus.Dir = opts.LocalPath
	statusBytes, errStatus := cmdStatus.Output()
	if errStatus != nil {
		if ctx.Err() != nil {
			return fmt.Errorf("failed to inspect working tree: %w: %v", errStatus, ctx.Err())
		}
		return fmt.Errorf("failed to inspect working tree: %w", errStatus)
	}
	if len(strings.TrimSpace(string(statusBytes))) > 0 {
		cmdAdd := exec.CommandContext(ctx, "git", "add", ".")
		cmdAdd.Dir = opts.LocalPath
		if err := cmdAdd.Run(); err != nil {
			if ctx.Err() != nil {
				return fmt.Errorf("failed to stage files: %w: %v", ctx.Err(), err)
			}
			return fmt.Errorf("failed to stage files: %w", err)
		}

		commitMsg := opts.CommitMsg
		if commitMsg == "" {
			commitMsg = "Push via Repo-lyzer"
		}
		cmdCommit := exec.CommandContext(ctx, "git", "commit", "-m", commitMsg)
		cmdCommit.Dir = opts.LocalPath
		if output, err := cmdCommit.CombinedOutput(); err != nil {
			outStr := string(output)
			if !strings.Contains(strings.ToLower(outStr), "nothing to commit") && !strings.Contains(strings.ToLower(outStr), "working tree clean") {
				if ctx.Err() != nil {
					return fmt.Errorf("failed to commit changes: %s: %w: %v", strings.TrimSpace(outStr), err, ctx.Err())
				}
				return fmt.Errorf("failed to commit changes: %s: %w", strings.TrimSpace(outStr), err)
			}
		}
	}

	// 4. Remote configuration
	cmdGetRemote := exec.CommandContext(ctx, "git", "remote", "get-url", "origin")
	cmdGetRemote.Dir = opts.LocalPath
	origBytes, errGet := cmdGetRemote.Output()
	hasOrigin := errGet == nil
	var originalURL string
	hasCustomRemote := false

	repoURL := fmt.Sprintf("https://github.com/%s/%s.git", opts.RepoOwner, opts.RepoName)

	if hasOrigin {
		originalURL = strings.TrimSpace(string(origBytes))
		cmdSetRemote := exec.CommandContext(ctx, "git", "remote", "set-url", "origin", repoURL)
		cmdSetRemote.Dir = opts.LocalPath
		if err := cmdSetRemote.Run(); err != nil {
			if ctx.Err() != nil {
				return fmt.Errorf("failed to update remote origin URL to %s in %s: %w: %v", repoURL, opts.LocalPath, err, ctx.Err())
			}
			return fmt.Errorf("failed to update remote origin URL to %s in %s: %w", repoURL, opts.LocalPath, err)
		}
		hasCustomRemote = true
	} else {
		cmdAddRemote := exec.CommandContext(ctx, "git", "remote", "add", "origin", repoURL)
		cmdAddRemote.Dir = opts.LocalPath
		if err := cmdAddRemote.Run(); err != nil {
			if ctx.Err() != nil {
				return fmt.Errorf("failed to add remote origin: %w: %v", ctx.Err(), err)
			}
			return fmt.Errorf("failed to add remote origin: %w", err)
		}
	}

	// Clean up remote configuration on function exit
	success := false
	defer func() {
		_ = success
		if !hasOrigin {
			cmdRemove := exec.CommandContext(context.Background(), "git", "remote", "remove", "origin")
			cmdRemove.Dir = opts.LocalPath
			_ = cmdRemove.Run()
		} else if hasCustomRemote && originalURL != "" {
			cmdRestore := exec.CommandContext(context.Background(), "git", "remote", "set-url", "origin", originalURL)
			cmdRestore.Dir = opts.LocalPath
			_ = cmdRestore.Run()
		}
	}()

	// 5. Push to GitHub
	var cmdPush *exec.Cmd
	username := opts.Username
	if username == "" {
		username = os.Getenv("REPO_LYZER_USERNAME")
	}
	if username == "" {
		username = "oauth2"
	}

	if opts.Token != "" {
		cmdPush = exec.CommandContext(ctx, "git", "-c", "credential.helper=!f() { echo username=$REPO_LYZER_USERNAME; echo password=$REPO_LYZER_TOKEN; }; f", "push", "-u", "origin", branch)
		cmdPush.Env = append(os.Environ(),
			"REPO_LYZER_USERNAME="+username,
			"REPO_LYZER_TOKEN="+opts.Token,
		)
	} else {
		cmdPush = exec.CommandContext(ctx, "git", "push", "-u", "origin", branch)
	}
	cmdPush.Dir = opts.LocalPath
	if err := cmdPush.Run(); err != nil {
		if ctx.Err() != nil {
			return fmt.Errorf("failed to push: %w: %v", ctx.Err(), err)
		}
		return fmt.Errorf("failed to push: %w", err)
	}
	success = true

	return nil
}

