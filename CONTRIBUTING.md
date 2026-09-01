# Contributing to Repo-lyzer

Thanks for your interest in contributing to Repo-lyzer.

Repo-lyzer is a Go-based CLI that analyzes GitHub repositories and presents insights in the terminal. This guide explains how to set up the project, follow the workflow, and submit changes for review.

## 1. Project Setup

### Prerequisites

- Go 1.24 or newer
- Git
- A GitHub account
- Optional: a GitHub personal access token for higher API limits

### Clone and run locally

```bash
git clone https://github.com/agnivo988/Repo-lyzer.git
cd Repo-lyzer
go mod download
go build -o repo-lyzer main.go
```

### Try it

```bash
./repo-lyzer analyze golang/go
```

If you are on Windows, use `repo-lyzer.exe` instead of `./repo-lyzer`.

## 2. Development Workflow

- Create a branch from `main` before making changes.
- Use clear branch names such as `feature/your-feature-name`, `fix/bug-name`, or `docs/doc-name`.
- Keep commits focused on a single change.
- Use conventional commit prefixes such as `feat:`, `fix:`, `docs:`, `test:`, and `chore:`.

Example:

```bash
git checkout -b docs/add-contributing
git commit -m "docs: add contributing guide"
```

## 3. Code Style

- Run `go fmt ./...` before committing.
- Run `go vet ./...` for static analysis.
- Run `go test ./...` to confirm the project still builds and tests pass.
- Keep changes small, readable, and consistent with existing code.
- Avoid committing secrets, generated files, or build artifacts.

For documentation changes, make sure Markdown renders cleanly on GitHub.

## 4. Submitting a Pull Request

Before opening a PR, please make sure:

- You have self-reviewed the change.
- Tests pass locally.
- Formatting is clean.
- The PR links the related issue.

Helpful PR description items:

- What changed
- Why the change was needed
- Any screenshots or logs if relevant

To reference an issue in your PR, use `Fixes #123` or `Closes #123`.

## 5. GSSoC 2026 Contributor Guide

If you are contributing through GSSoC 2026:

- Comment on the issue to claim it before starting work.
- Wait for a maintainer to confirm that you can work on it.
- Mention the issue number in your PR.
- Use a PR title format like `[GSSoC 2026] feat: short description`.
- Keep your branch and commits focused on the issue you claimed.

Suggested point structure:

- Level 1: Beginner-friendly documentation or small fixes
- Level 2: Moderate feature work or refactors
- Level 3: Advanced feature work or architecture changes

## 6. Issue Labels Guide

Common labels used in this repository:

| Label | Meaning |
| --- | --- |
| `good first issue` | Beginner-friendly task |
| `level1` | Low-complexity contribution |
| `level2` | Medium-complexity contribution |
| `level3` | Advanced contribution |
| `documentation` | Docs-only work |
| `gssoc-2026` | GSSoC 2026 issue |
| `gssoc-ext` | GSSoC Extended Program issue |

If you are unsure whether an issue is a good fit, leave a comment on the issue and ask a maintainer.