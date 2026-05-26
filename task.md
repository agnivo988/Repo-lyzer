# Task Checklist

- [x] Core & Monitor Changes
  - [x] Implement `SetToken` in `internal/monitor/monitor.go`
- [x] TUI (Bubble Tea) Changes
  - [x] Update `DashboardModel` in `internal/ui/dashboard.go` (add `token` field, `SetToken` method, set token in `apiStatusView`)
  - [x] Update `MonitorDashboardModel` in `internal/ui/monitor_dashboard.go` (add `token` field, update `NewMonitorDashboardModel`, update `startMonitoring`)
  - [x] Update `CloneInputModel` in `internal/ui/clone_input.go` (use `github.ParseGitHubURL`)
  - [x] Update `MainModel` in `internal/ui/app.go` (propagate error to `cloneInput`, parse URLs using `github.ParseGitHubURL`, propagate token from settings to clients, pass token to dashboard and monitor dashboard)
- [x] CLI / Cmd Changes
  - [x] Update `runSummary` in `cmd/analyze.go` (add private repo token prompting)
  - [x] Update `cmd/compare.go` (add private repo token prompting for both repos)
  - [x] Update `cmd/monitor.go` (pre-validate repo and prompt/set GITHUB_TOKEN environment variable)
- [x] Verification
  - [x] Run dry-run commands and compile checks

