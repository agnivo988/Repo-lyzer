# Repo-lyzer

<img src="https://res.cloudinary.com/dhyii4oiw/image/upload/v1767324445/Screenshot_2026-01-02_085503_ros5gz.png" alt="This is the logo" width="auto" height="auto">

**Repo-lyzer** is a modern, terminal-based CLI tool written in **Golang** that analyzes GitHub repositories and presents insights in a beautifully formatted, interactive dashboard. It is designed for developers, recruiters, and open-source enthusiasts to quickly evaluate a repository’s health, activity, and contributor statistics.

---
## 👥 Who is Repo-lyzer for?

- 👨‍💻 **Developers** evaluating open-source projects  
- 🧑‍💼 **Recruiters** assessing repository health and activity  
- 🌱 **Contributors** exploring project structure and engagement  

---


## 🌟 Features

- **Repository Overview:** Shows stars, forks, open issues, and general info.
- **Language Breakdown:** Displays percentage of languages used with colored bars.
- **Commit Activity:** Horizontal graph showing commit frequency over the past year.
- **Health Score:** Calculates repository health based on activity and contributor stats.
- **Bus Factor:** Measures critical contributors to assess project risk.
- **Repo Maturity Score:** Evaluates repository age, activity, and structure.
- **Recruiter Summary:** Quick summary highlighting key metrics for recruitment evaluation.
- **File Tree Viewer:** Explore the repository's file structure directly in the dashboard.
- **Export Options:** Export analysis results to JSON or Markdown.
- **Compare Mode:** Compare two repositories side by side.
- **Interactive CLI Menu:** Fully navigable TUI with keyboard arrows, input prompts, and instant feedback.
- **Colorized Output:** Uses neon-style colors and ASCII styling for a modern CLI experience.

---

## 🛠 Tech Stack & Libraries Used

- **[Golang](https://golang.org/)** – Core language for CLI development
- **[Cobra](https://github.com/spf13/cobra)** – CLI command management
- **[Bubble Tea](https://github.com/charmbracelet/bubbletea)** – Terminal-based interactive UI
- **[Bubbles](https://github.com/charmbracelet/bubbles)** – TUI components (Spinners, etc.)
- **[Lipgloss](https://github.com/charmbracelet/lipgloss)** – Styling terminal output, colors, borders, alignment
- **[Tablewriter](https://github.com/olekukonko/tablewriter)** – Beautiful tables in the terminal
- **[x/term](https://pkg.go.dev/golang.org/x/term)** – Terminal size detection
- **[GitHub REST API](https://docs.github.com/en/rest)** – Fetching repo, commits, issues, and contributors

---

## 📝 Project Overview

Repo-lyzer is structured in a **modular architecture** for scalability and maintainability:
```text

repo-analyzer/
│
├── cmd/
│ ├── root.go
│ ├── analyze.go
│ └── compare.go
│
├── internal/
│ ├── github/
│ ├── analyzer/
│ └── output/
│
├── docs/
│ ├── DOCUMENTATION_INDEX.md
│ ├── QUICK_REFERENCE.md
│ ├── IMPLEMENTATION_DETAILS.md
│ ├── ANALYZER_INTEGRATION.md
│ ├── CHANGE_LOG.md
│ ├── PHASE2_README.md
│ ├── PHASE2_DOCUMENTATION_INDEX.md
│ ├── PHASE2_QUICK_REFERENCE.md
│ ├── PHASE2_COMPLETION.md
│ ├── PHASE2_FINAL_SUMMARY.md
│ ├── PHASE2_FILE_CHANGES.md
│ ├── FEATURE_INVENTORY_PHASE2.md
│ ├── UI_FIXES_SUMMARY.md
│ ├── RESOLUTION_REPORT.md
│ └── TODO.md
│
├── config/ 
├── main.go 
├── go.mod 
└── README.md
```

**Workflow:**

1. User launches `repo-lyzer`.
2. Interactive TUI menu guides the user to select **Analyze** or **Compare** mode.
3. CLI fetches data from GitHub API (repos, commits, contributors, languages).
4. Computes health, bus factor, maturity score, and recruiter summary.
5. Displays results in a **beautifully styled, centered terminal dashboard**.

---

## 📚 Documentation

For detailed information about the project, please refer to the documentation in the `docs/` folder:

- **[DOCUMENTATION_INDEX.md](docs/DOCUMENTATION_INDEX.md)** – Main index for all documentation
- **[QUICK_REFERENCE.md](docs/QUICK_REFERENCE.md)** – Quick reference guide
- **[IMPLEMENTATION_DETAILS.md](docs/IMPLEMENTATION_DETAILS.md)** – Technical implementation details
- **[ANALYZER_INTEGRATION.md](docs/ANALYZER_INTEGRATION.md)** – Analyzer integration guide
- **[CHANGE_LOG.md](docs/CHANGE_LOG.md)** – Changelog and version history
- **[PHASE2_README.md](docs/PHASE2_README.md)** – Phase 2 development overview

---

## 🚧 Challenges Faced

- **Terminal layout:** Centering multiple sections with minimal gaps using Lipgloss and Bubble Tea.
- **Dynamic data rendering:** Handling repositories with very high activity or many contributors without breaking the TUI.
- **GitHub API rate limits:** Required careful handling of requests and optional support for personal access tokens.
- **Horizontal commit graphs:** Creating readable, interactive graphs for commit activity over a year.
- **Unified dashboard:** Combining multiple sections while maintaining responsive vertical and horizontal alignment.

---

## 🛠 Installation (From Source)

1. Clone the repo:

```bash
git clone https://github.com/agnivo988/Repo-lyzer.git
cd Repo-lyzer
```
---

## 🚀 Quick Start (For Contributors)

If you want to run Repo-lyzer locally or contribute to the project, follow the steps below.

### Prerequisites
- Go **v1.20+**
- Git
- Internet connection (for GitHub API access)

### Run Locally

```bash
go mod tidy
go run main.go
```
This will launch the interactive Repo-lyzer terminal dashboard.

---

## ▶️ Usage Examples

Analyze a repository:
```bash
repo-lyzer analyze golang/go
```
**🔄 Compare two repositories**
Repository comparison is available through the interactive menu.  
Launch the application and select **Compare Repositories** from the dashboard.

**📝 Export Analysis Results**
Export functionality is available from within the interactive dashboard.  
After analysis, choose the export option from the menu to save results.

---

## ⌨️ Keyboard Shortcuts

| Key | Action |
|----|-------|
| `Enter` | Confirm input / run analysis |
| `ESC` | Go back / cancel current action |
| `q` | Quit application (from main menu) |
| `H` | Open analysis history |
| `↑ / ↓` | Navigate menu or history |
| `j ` | Navigate lists (vim-style) |
| `m` | Export comparison as Markdown |
| `j` | Export comparison as JSON |

These shortcuts make Repo-lyzer feel fast and efficient for terminal power users.
---

## 🔐 GitHub API Configuration (Optional)

Repo-lyzer uses the GitHub REST API.
To get higher rate limits and smoother analysis, you can provide a GitHub Personal Access Token.

## ⚙️ Setup

- 🔑 Generate a token from: 
  [GitHub Personal Access Tokens](https://github.com/settings/tokens)

- 🌍 Set it as an environment variable:
```bash
export GITHUB_TOKEN=your_token_here   # macOS/Linux
setx GITHUB_TOKEN your_token_here     # Windows
```

ℹ️ If no token is provided, Repo-lyzer will use GitHub’s public rate limits.
---

## How it looks
<img src="https://res.cloudinary.com/dhyii4oiw/image/upload/v1767290545/Screenshot_2026-01-01_224310_c0hhr8.png" alt="This is the Opening Screen" width="auto" height="auto">
<img src="https://res.cloudinary.com/dhyii4oiw/image/upload/v1767324721/Screenshot_2026-01-02_090050_u6xweq.png" alt="" width="auto" height="auto">
<img src="https://res.cloudinary.com/dhyii4oiw/image/upload/v1767324721/Screenshot_2026-01-02_090043_keqfs4.png" alt="" width="auto" height="auto">
<img src="https://res.cloudinary.com/dhyii4oiw/image/upload/v1767324721/Screenshot_2026-01-02_090104_dm7bgk.png" alt="" width="auto" height="auto">
<img src="https://res.cloudinary.com/dhyii4oiw/image/upload/v1767324829/Screenshot_2026-01-02_090335_acms5i.png" alt="" width="auto" height="auto">

---

## Installation (For Users)
1. Go to the terminal and run
```bash
go install github.com/agnivo988/Repo-lyzer@v1.0.0
```

2. Then run 
```bash
repo-lyzer
```
---


## 🛠 Troubleshooting

**Invalid repository input error**
- Ensure the format is `owner/repo` (example: `golang/go`)
- Avoid pasting full GitHub URLs

**GitHub API rate limit issues**
- Set a `GITHUB_TOKEN` for higher rate limits

**Application doesn’t start**
- Verify Go version is `v1.20+`
- Run `go mod tidy` before `go run main.go`

---

## License
MIT License © 2026 Agniva Mukherjee


