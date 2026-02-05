# 🤝 Contributing to Repo-lyzer

Thanks for your interest in contributing to **Repo-lyzer** 🎉  
Repo-lyzer is a modern, terminal-based GitHub repository analyzer built with **Go**, **Bubble Tea**, and **Lipgloss**. Contributions of all sizes are welcome!

This document explains how to set up the project, make changes, and submit them properly.

---

## 🚀 Ways to Contribute

You can help Repo-lyzer by:

- 🐛 Reporting bugs  
- ✨ Proposing or implementing new features  
- 🧠 Improving scoring algorithms (health, bus factor, maturity)  
- 🎨 Enhancing TUI layout, styling, or UX  
- 🧪 Adding tests or improving reliability  
- 📚 Improving documentation  

---

## 🛠 Getting Started

### Prerequisites

- **Star the Repo-lyzer Repository**

Make sure you have the following installed:

- **Go 1.21+**  
- **Git**  
- **A GitHub account**  
- *(Optional)* GitHub Personal Access Token to avoid rate limits  

---

## 📦 Project Setup

### 1. Fork the Repository

Click **Fork** on GitHub and clone your fork:

```bash
git clone https://github.com/YOUR_USERNAME/Repo-lyzer.git
cd Repo-lyzer
```

### 2. Install Dependencies

```bash
go mod tidy
```

### 3. Run the CLI Locally

```bash
go run main.go
```

Or Install It Locally:

```bash
go install
repo-lyzer
```

---

## 🗂 Project Structure

```
repo-analyzer/
├── cmd/               # Cobra commands (analyze, compare, root)
├── internal/
│   ├── github/        # GitHub API interactions
│   ├── analyzer/      # Scoring logic (health, bus factor, maturity)
│   └── output/        # Bubble Tea + Lipgloss TUI rendering
├── config/            # Configuration & token handling
├── main.go
├── go.mod
└── README.md
```

---

## 📐 Guidelines

- Keep GitHub API logic inside `internal/github`
- Keep scoring logic isolated in `internal/analyzer`
- Keep UI rendering inside `internal/output`
- Avoid mixing API, logic, and UI code

---

## 🔐 GitHub API & Tokens
Repo-lyzer uses the GitHub REST API and may hit rate limits.

Optional: Use a Personal Access Token.

Create a token and export it:

```bash
export GITHUB_TOKEN=your_token_here
```

Repo-lyzer will automatically detect and use it.

⚠️ Never commit tokens or secrets.

---

## 🧑‍💻 Coding Guidelines

### Go Style

- Follow standard Go formatting (`gofmt`)
- Keep functions small and readable
- Prefer clear logic over clever tricks

### TUI & UX

- Keep layouts responsive to terminal size
- Avoid hardcoded widths when possible
- Ensure output remains readable on small terminals
- Test changes with repositories of different sizes and activity levels

### ⚡ Performance

- Minimize unnecessary API calls
- Handle large repositories gracefully
- Cache or batch requests when possible

---

## ➕ Adding Features

Before implementing large features:

- Open an issue describing the idea
- Explain why it's useful
- Discuss the approach if it affects architecture or scoring

This avoids duplicated work and design conflicts.

---

## 🧪 Testing

While automated tests are limited due to API usage and TUI complexity, manually test with:

- Highly active repositories
- Old but inactive repositories
- Both Analyze and Compare modes
- With and without a GitHub token

---

## 🔀 Submitting a Pull Request

### 1. Create a Branch

```bash
git checkout -b feature/your-feature-name
```

### 2. Make Your Changes

Ensure:

- Code builds successfully
- TUI renders correctly
- No secrets are committed

### 3. Commit Your Changes

```bash
git commit -m "Add: short description of change"
```

### 4. Push and Open a PR

```bash
git push -u origin feature/your-feature-name
```

---

## 📝 Pull Request Guidelines

When opening a PR, include:

- What the change does
- Why it's needed
- Screenshots (for UI changes)

### Commit Message Examples

- `Add: recruiter summary panel`
- `Fix: handle empty contributor list`
- `Improve: commit activity graph rendering`
- `Refactor: separate GitHub client logic`

---

## 📜 Code of Conduct

This project follows a [Code of Conduct](CODE_OF_CONDUCT.md).  
By contributing, you agree to uphold respectful and inclusive behavior in all interactions.

---

## ❓ Need Help?

- Open an **Issue** for bugs or feature requests
- Use **Discussions** for questions or ideas
- Check existing issues before submitting new ones

---

## 🙏 Thank You!

Your contributions help make Repo-lyzer better for everyone. Happy coding! 🚀






