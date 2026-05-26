# Plan: Repository Analysis Gracefulness for Invalid/Empty/Private URLs

This plan improves the repository analysis error handling and authentication token routing for invalid, empty, or private URLs across both CLI and TUI (Bubble Tea) interfaces.

## User Review Required

> [!NOTE]
> All CLI commands (including summary analysis, comparison, and monitoring) will now gracefully prompt the user for their GitHub access token if they attempt to access a private repository and do not have a token configured via environment variables.

> [!IMPORTANT]
> The Bubble Tea TUI will now seamlessly route the GitHub token configured in the settings screen to all background analysis operations, comparisons, the file viewer, and the real-time monitoring/dashboard screens. Previously, the settings token was not propagated to these client initializations, preventing private repository analysis and resulting in low rate limits.

## Proposed Changes

### CLI Component

#### [MODIFY] [cmd/analyze.go](file:///d:/rayzer/Repo-lyzer/cmd/analyze.go)
- Add GitHub access token prompting to the `runSummary` function when a repository is private and no token is configured.
- Clean up any manual string checks by utilizing the robust `github.ParseGitHubURL` parsing throughout analysis commands.

#### [MODIFY] [cmd/compare.go](file:///d:/rayzer/Repo-lyzer/cmd/compare.go)
- Add GitHub access token prompting to both first and second repository fetches when they are private and no token is configured.

#### [MODIFY] [cmd/monitor.go](file:///d:/rayzer/Repo-lyzer/cmd/monitor.go)
- Before beginning the monitor loop in CLI mode, check if the repository is private/accessible.
- Prompt for a GitHub token if needed and store it in the current process's environment via `os.Setenv` to propagate it to the monitor client.

---

### Core & Monitor Component

#### [MODIFY] [internal/monitor/monitor.go](file:///d:/rayzer/Repo-lyzer/internal/monitor/monitor.go)
- Add a `SetToken(token string)` method to the `Monitor` struct so that the token can be set on its internal GitHub client.

---

### TUI / Bubble Tea Component

#### [MODIFY] [internal/ui/app.go](file:///d:/rayzer/Repo-lyzer/internal/ui/app.go)
- In `NewMainModel`, initialize the `monitorDashboard` field with the token from settings if present.
- In `SetStateMonitorDashboard`, pass the saved settings token to `NewMonitorDashboardModel`.
- In `analyzeRepo`, use `github.ParseGitHubURL` to validate and extract the owner and repository name, instead of raw `strings.Split`. Set the client's token from `m.appConfig.GitHubToken`.
- In `compareRepos`, use `github.ParseGitHubURL` for both repository URLs and set the client's token from settings.
- In `checkOwnership`, set the client's token from settings to allow ownership checks for private repositories.
- In `cloneRepo`, use `github.ParseGitHubURL` to allow cloning via full GitHub URLs, not just `owner/repo` formats.
- In the `MainModel.View` method for `stateCloneInput`, assign `m.cloneInput.err = m.err` so cloning errors are actually displayed on the clone screen.
- In the `stateLoading` and `CachedAnalysisResult` handlers, pass the configured settings token to the `dashboard` model.

#### [MODIFY] [internal/ui/clone_input.go](file:///d:/rayzer/Repo-lyzer/internal/ui/clone_input.go)
- Replace raw string checks under `tea.KeyEnter` with `github.ParseGitHubURL` validation to allow full GitHub URLs to be cloned.

#### [MODIFY] [internal/ui/dashboard.go](file:///d:/rayzer/Repo-lyzer/internal/ui/dashboard.go)
- Add a `token` field to `DashboardModel` along with a `SetToken` method.
- Pass the token to the client in the `apiStatusView` method so that rate limit displays reflect token capabilities.

#### [MODIFY] [internal/ui/monitor_dashboard.go](file:///d:/rayzer/Repo-lyzer/internal/ui/monitor_dashboard.go)
- Add `token` field and update `NewMonitorDashboardModel` to accept it.
- Apply the token to the monitor in `startMonitoring`.

## Verification Plan

### Automated Tests
- Since `go` test runner is not configured in the host's system path, we will verify the code builds and is syntactically correct by compiling it.
- Build the binary using `go build` (or similar command) if available, or visually inspect the changes closely.

### Manual Verification
- We will verify using dry-run modes, running CLI commands, and checking the code compiles.
