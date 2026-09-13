# Contributing to Repo-lyzer

Thank you for your interest in contributing to **Repo-lyzer**! We welcome contributions from everyone. By participating in this project, you agree to abide by the [Code of Conduct](https://github.com/agnivo988/Repo-lyzer/blob/main/CODE_OF_CONDUCT.md).

## Table of Contents

- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Fork and Clone](#fork-and-clone)
  - [Building from Source](#building-from-source)
- [Development Workflow](#development-workflow)
  - [Branching](#branching)
  - [Coding Standards](#coding-standards)
  - [Running Tests](#running-tests)
- [Submitting a Pull Request](#submitting-a-pull-request)
- [Reporting Issues](#reporting-issues)

## Getting Started

### Prerequisites

- **Go 1.24+** — install via [gvm](https://github.com/moovweb/gvm) or [official installer](https://go.dev/dl/)
- **Git** — version 2.30 or later
- A **GitHub Personal Access Token** (classic) with `repo` and `read:org` scopes

### Fork and Clone

1. Fork the repository on GitHub.
2. Clone your fork locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/Repo-lyzer.git
   cd Repo-lyzer
   ```
3. Add the upstream remote:
   ```bash
   git remote add upstream https://github.com/agnivo988/Repo-lyzer.git
   ```

### Building from Source

```bash
go build -o repo-lyzer ./cmd
```

The binary will be placed at `./repo-lyzer` in the project root. Run it with:

```bash
./repo-lyzer analyze --token YOUR_GITHUB_TOKEN --owner agnivo988 --repo Repo-lyzer
```

## Development Workflow

### Branching

- Create a feature branch off `main`: `git checkout -b feat/my-feature`
- Use a descriptive branch name: `fix/login-error`, `feat/dashboard-widget`, `docs/contributing-guide`

### Coding Standards

- Format all Go code with `gofmt` before committing:
  ```bash
  gofmt -s -w .
  ```
- Follow idiomatic Go conventions as described in [Effective Go](https://go.dev/doc/effective_go).
- Use meaningful variable names and avoid abbreviations where clarity matters.
- Keep functions focused and under 50 lines where practical.
- Add comments for exported functions, types, and package-level declarations.

### Running Tests

Run the full test suite:

```bash
go test ./...
```

Run tests for a specific package:

```bash
go test ./cmd/...
```

### Pre-commit Checklist

Before committing, ensure:

- [ ] `go build ./...` passes with no errors
- [ ] `go test ./...` passes with no failures
- [ ] `gofmt -s -w .` has been run
- [ ] New functionality includes tests
- [ ] Changes are accompanied by updated documentation

## Submitting a Pull Request

1. Push your branch to your fork:
   ```bash
   git push origin feat/my-feature
   ```
2. Open a pull request on GitHub against the `main` branch.
3. In the PR description, explain:
   - What the change does
   - Why it is needed
   - How you tested it
4. Reference any related issues (e.g., `Fixes #123`).
5. A maintainer will review your PR. Please respond to feedback promptly.

## Reporting Issues

- Use the [GitHub issue tracker](https://github.com/agnivo988/Repo-lyzer/issues) to report bugs or request features.
- Search existing issues before opening a new one.
- Include as much detail as possible: steps to reproduce, expected behavior, actual behavior, environment details.

### GSSoC 2026 Contributors

If you are contributing as part of **GSSoC 2026**:
- Comment on the issue you want to work on to get it assigned.
- Tag your PR with the relevant GSSoC labels.
- Reach out on the project's communication channel if you need guidance.

---

Thank you for helping make Repo-lyzer better!
