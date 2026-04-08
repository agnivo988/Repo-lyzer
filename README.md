## Maintainers
- @agnivo988

## Contributors

<a href="https://github.com/Aamod007"><img src="https://github.com/Aamod007.png" width="50" height="50" alt="Aamod007" title="Contributor"></a> <a href="https://github.com/Aditya8369"><img src="https://github.com/Aditya8369.png" width="50" height="50" alt="Aditya8369" title="Contributor"></a> <a href="https://github.com/agnivo988"><img src="https://github.com/agnivo988.png" width="50" height="50" alt="agnivo988" title="Project Maintainer"></a> <a href="https://github.com/Gupta-02"><img src="https://github.com/Gupta-02.png" width="50" height="50" alt="Gupta-02" title="Contributor"></a> <a href="https://github.com/GauravKarakoti"><img src="https://github.com/GauravKarakoti.png" width="50" height="50" alt="GauravKarakoti" title="Contributor"></a> <a href="https://github.com/Sappymukherjee214"><img src="https://github.com/Sappymukherjee214.png" width="50" height="50" alt="Sappymukherjee214" title="Contributor"></a> <a href="https://github.com/ItsMeArm00n"><img src="https://github.com/ItsMeArm00n.png" width="50" height="50" alt="ItsMeArm00n" title="Contributor"></a> <a href="https://github.com/MuktaRedij"><img src="https://github.com/MuktaRedij.png" width="50" height="50" alt="MuktaRedij" title="Contributor"></a> <a href="https://github.com/Kiran95021"><img src="https://github.com/Kiran95021.png" width="50" height="50" alt="Kiran95021" title="Contributor"></a> <a href="https://github.com/Shriii19"><img src="https://github.com/Shriii19.png" width="50" height="50" alt="Shriii19" title="Contributor"></a> <a href="https://github.com/KUMARI-SONALIUPADHYAY"><img src="https://github.com/KUMARI-SONALIUPADHYAY.png" width="50" height="50" alt="KUMARI-SONALIUPADHYAY" title="Contributor"></a> <a href="https://github.com/magic-peach"><img src="https://github.com/magic-peach.png" width="50" height="50" alt="magic-peach" title="Contributor"></a> <a href="https://github.com/coderabbitai"><img src="https://github.com/coderabbitai.png" width="50" height="50" alt="coderabbitai[bot]" title="Bot Contributor"></a> <a href="https://github.com/sahoo-tech"><img src="https://github.com/sahoo-tech.png" width="50" height="50" alt="sahoo-tech" title="Contributor"></a> <a href="https://github.com/Abhijeet-980"><img src="https://github.com/Abhijeet-980.png" width="50" height="50" alt="Abhijeet-980" title="Contributor"></a> <a href="https://github.com/Diksha78-bot"><img src="https://github.com/Diksha78-bot.png" width="50" height="50" alt="Diksha78-bot" title="Bot Contributor"></a> <a href="https://github.com/Shivani-Meena07"><img src="https://github.com/Shivani-Meena07.png" width="50" height="50" alt="Shivani-Meena07" title="Contributor"></a> <a href="https://github.com/ShashankSaga"><img src="https://github.com/ShashankSaga.png" width="50" height="50" alt="ShashankSaga" title="Contributor"></a>


<h1 align="center">Repo-lyzer</h1>
<p align="center">
  <img src="https://res.cloudinary.com/dhyii4oiw/image/upload/v1767324445/Screenshot_2026-01-02_085503_ros5gz.png" alt="Repo-lyzer Logo" width="300">
</p>

**Repo-lyzer** is a modern, terminal-based CLI tool written in **Golang** that analyzes GitHub repositories and presents insights in a beautifully formatted, interactive dashboard.  
It is designed for **developers, recruiters, and open-source enthusiasts** to quickly evaluate a repository’s health, activity, and contributor statistics.

---

## Who is Repo-lyzer for?

-  **Developers** evaluating open-source projects  
-  **Recruiters** assessing repository health and activity  
-  **Contributors** exploring project structure and engagement  

---

## Features

- **Repository Overview** – Stars, forks, open issues, general info  
- **Language Breakdown** – Percentage of languages used with colored bars  
- **Commit Activity** – Horizontal graph showing commit frequency over the past year  
- **Health Score** – Activity & contributor-based scoring  
- **Bus Factor** – Measures critical contributors to assess project risk  
- **Repo Maturity Score** – Evaluates age, activity, and structure  
- **Recruiter Summary** – Quick snapshot for hiring evaluation  
- **Quick Summary** – Fast 5-line overview with key metrics (commits 30d, top language, contributors, health, last commit)  
- **File Tree Viewer** – Browse repository structure in-dashboard  
- **Export Options** – Export analysis as JSON, Markdown, CSV, or HTML  
- **Compare Mode** – Side-by-side repository comparison  
- **Interactive CLI Menu** – Fully navigable TUI (keyboard driven)  
- **Colorized Output** – Neon-style colors and ASCII styling  
- **Settings Persistence** – Theme, export preferences, and GitHub token saved locally  

---

## Tech Stack & Libraries

- **[Golang](https://golang.org/)** – Core language  
- **[Cobra](https://github.com/spf13/cobra)** – CLI command management  
- **[Bubble Tea](https://github.com/charmbracelet/bubbletea)** – Interactive TUI  
- **[Bubbles](https://github.com/charmbracelet/bubbles)** – UI components  
- **[Lipgloss](https://github.com/charmbracelet/lipgloss)** – Styling & layout  
- **[Tablewriter](https://github.com/olekukonko/tablewriter)** – Terminal tables  
- **[x/term](https://pkg.go.dev/golang.org/x/term)** – Terminal size detection  
- **GitHub REST API** – Repository, commits, issues, contributors  

---

## Project Overview

Repo-lyzer follows a **modular architecture** for scalability and maintainability.

```
repo-analyzer/
│
├── cmd/
│   ├── root.go
│   ├── analyze.go
│   └── compare.go
│
├── internal/
│   ├── github/       # GitHub API client
│   ├── analyzer/     # Metric computations
│   ├── cache/        # Offline caching
│   ├── config/       # Settings persistence
│   └── ui/           # TUI components
│
├── docs/
│   ├── DOCUMENTATION_INDEX.md
│   ├── QUICK_REFERENCE.md
│   ├── IMPLEMENTATION_DETAILS.md
│   ├── ANALYZER_INTEGRATION.md
│   └── CHANGE_LOG.md
│
├── main.go
├── go.mod
└── README.md
```

### Workflow

1. User launches `repo-lyzer`  
2. Interactive menu → **Analyze** or **Compare**  
3. GitHub API fetch (repos, commits, contributors, languages)  
4. Metrics computed (health, bus factor, maturity)  
5. Displayed in **centered, styled terminal dashboard**

---

## Architecture Overview

```
┌────────────────────────────────────────────┐
│               main.go                      │
└────────────────────────────────────────────┘
                    │
                    ▼
┌────────────────────────────────────────────┐
│                 cmd/                       │
└────────────────────────────────────────────┘
                    │
                    ▼
┌────────────────────────────────────────────┐
│             internal/ui/                   │
└────────────────────────────────────────────┘
          │           │           │
          ▼           ▼           ▼
┌──────────────┐ ┌──────────────┐ ┌──────────────┐
│   github     │ │   analyzer   │ │   output     │
└──────────────┘ └──────────────┘ └──────────────┘
```

### Key Directories

| Directory | Purpose |
|---------|---------|
| `cmd/` | CLI commands |
| `internal/github/` | GitHub API client |
| `internal/analyzer/` | Metric computations |
| `internal/cache/` | Offline caching |
| `internal/config/` | Settings persistence |
| `internal/ui/` | TUI components |
| `internal/output/` | Formatting & rendering |
| `docs/` | Documentation |

---

## Documentation

### For Contributors
- [ARCHITECTURE.md](docs/ARCHITECTURE.md) – Complete architecture guide  
- [ANALYZER_INTEGRATION.md](docs/ANALYZER_INTEGRATION.md) – Adding new analyzers  
- [IMPLEMENTATION_DETAILS.md](docs/IMPLEMENTATION_DETAILS.md) – Technical deep dive  

### Reference
- [DOCUMENTATION_INDEX.md](docs/DOCUMENTATION_INDEX.md) – Master index  
- [QUICK_REFERENCE.md](docs/QUICK_REFERENCE.md) – Quick usage guide  
- [CHANGE_LOG.md](docs/CHANGE_LOG.md) – Version history  

---

## Challenges Faced

- Centering multi-section layouts in terminal  
- Handling high-activity repositories  
- GitHub API rate limits  
- Horizontal commit graphs  
- Unified responsive dashboard  

---

## Installation (From Source)

```bash
git clone https://github.com/agnivo988/Repo-lyzer.git
cd Repo-lyzer
go mod tidy
go run main.go
```

---

## Usage

### Quick Summary (Fast Overview)
Get a quick 5-line summary of any repository:
```bash
repo-lyzer summary golang/go
```
Or use the flag with analyze:
```bash
repo-lyzer analyze --summary microsoft/vscode
```

**Example Output:**
```
📊 Repository Summary: golang/go
   Commits (30d): 30
   Top Language: Go
   Contributors: 381
   Health Score: 90/100
   Last Commit: 4 hours ago
```

### Analyze Repository
```bash
repo-lyzer analyze golang/go
```

### Save Analysis to File
Save the full analysis JSON to a file:
```bash
repo-lyzer analyze golang/go --save analysis.json
```

This will create/overwrite the file with valid JSON and print the file path:
```
✅ Analysis saved to: analysis.json
```

Combine with `--compact` to save to file AND print to stdout:
```bash
repo-lyzer analyze golang/go --save analysis.json --compact
```

### Compare Repositories
Available via interactive menu.

### Export Results
Export to JSON or Markdown from dashboard.

### Interactive CLI Menu
Repo-lyzer features a hierarchical, keyboard-driven menu system for easy navigation and access to all features.

#### Main Menu Options
- **📊 Analyze Repository** – Enter analysis submenu to choose analysis type
- **🔄 Compare Repositories** – Start side-by-side repository comparison
- **📜 View History** – Browse previously analyzed repositories
- **⚙️ Settings** – Access settings submenu for configuration
- **❓ Help** – Open help submenu for guidance
- **🚪 Exit** – Quit the application

#### Submenus

**Analysis Types** (accessed via 📊 Analyze Repository):
- **⚡ Quick Analysis** – Fast overview with key metrics
- **🔍 Detailed Analysis** – Comprehensive analysis with all metrics
- **⚙️ Custom Analysis** – User-configurable analysis options

**Settings** (accessed via ⚙️ Settings):
- **Theme Settings** – Customize UI themes and colors
- **Export Options** – Configure default export formats
- **GitHub Token** – Set personal access token for higher rate limits
- **Reset to Defaults** – Restore all settings to default values

**Help** (accessed via ❓ Help):
- **Keyboard Shortcuts** – View all available shortcuts
- **Getting Started** – Quick start guide
- **Features Guide** – Detailed feature explanations
- **Troubleshooting** – Common issues and solutions

#### Navigation
- **Arrow Keys** (↑/↓) – Navigate menu options
- **Enter** – Select option or enter submenu
- **ESC** – Return to previous menu level
- **Tab** – Quick navigation between sections

---
  <img src="https://res.cloudinary.com/dhyii4oiw/image/upload/v1767290545/Screenshot_2026-01-01_224310_c0hhr8.png" width="90%">
  <img src="https://res.cloudinary.com/dhyii4oiw/image/upload/v1767324721/Screenshot_2026-01-02_090050_u6xweq.png" width="90%">
  <img src="https://res.cloudinary.com/dhyii4oiw/image/upload/v1767324721/Screenshot_2026-01-02_090043_keqfs4.png" width="90%">
  <img src="https://res.cloudinary.com/dhyii4oiw/image/upload/v1767324721/Screenshot_2026-01-02_090104_dm7bgk.png" width="90%">
  <img src="https://res.cloudinary.com/dhyii4oiw/image/upload/v1767324829/Screenshot_2026-01-02_090335_acms5i.png" width="90%">
</p>

---

## Installation (For Users)

```bash
go install github.com/agnivo988/Repo-lyzer@v1.0.6
```
```bash
repo-lyzer
```

---

## Contributor Section

Contributions are welcome. If you want to improve Repo-lyzer, start by reviewing the architecture and analyzer documentation, then make focused changes with clear testable behavior.

### Local Setup for Contributors

```bash
git clone https://github.com/agnivo988/Repo-lyzer.git
cd Repo-lyzer
go mod tidy
go run main.go
```

### How to Contribute

- Fork the repository and create a feature branch
- Keep changes scoped to one feature, fix, or documentation update
- Update documentation when behavior or commands change
- Test your changes locally before opening a pull request
- Open a pull request with a short description of the problem and solution

---

## License

MIT License © 2026 Agniva Mukherjee
