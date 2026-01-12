<h1 align="center">Repo-lyzer</h1>

<p align="center">
  <img src="https://res.cloudinary.com/dhyii4oiw/image/upload/v1767324445/Screenshot_2026-01-02_085503_ros5gz.png" alt="Repo-lyzer Logo" width="300">
</p>

**Repo-lyzer** is a modern, terminal-based CLI tool written in **Golang** that analyzes GitHub repositories and presents insights in a beautifully formatted, interactive dashboard.  
It is designed for **developers, recruiters, and open-source enthusiasts** to quickly evaluate a repository’s health, activity, and contributor statistics.

---

## 👥 Who is Repo-lyzer for?

- 👨‍💻 **Developers** evaluating open-source projects  
- 🧑‍💼 **Recruiters** assessing repository health and activity  
- 🌱 **Contributors** exploring project structure and engagement  

---

## 🌟 Features

- **Repository Overview** – Stars, forks, open issues, general info  
- **Language Breakdown** – Percentage of languages used with colored bars  
- **Commit Activity** – Horizontal graph showing commit frequency over the past year  
- **Health Score** – Activity & contributor-based scoring  
- **Bus Factor** – Measures critical contributors to assess project risk  
- **Repo Maturity Score** – Evaluates age, activity, and structure  
- **Recruiter Summary** – Quick snapshot for hiring evaluation  
- **File Tree Viewer** – Browse repository structure in-dashboard  
- **Export Options** – Export analysis as JSON, Markdown, CSV, or HTML  
- **Compare Mode** – Side-by-side repository comparison  
- **Interactive CLI Menu** – Fully navigable TUI (keyboard driven)  
- **Colorized Output** – Neon-style colors and ASCII styling  
- **Settings Persistence** – Theme, export preferences, and GitHub token saved locally  

---

## 🛠 Tech Stack & Libraries

- **[Golang](https://golang.org/)** – Core language  
- **[Cobra](https://github.com/spf13/cobra)** – CLI command management  
- **[Bubble Tea](https://github.com/charmbracelet/bubbletea)** – Interactive TUI  
- **[Bubbles](https://github.com/charmbracelet/bubbles)** – UI components  
- **[Lipgloss](https://github.com/charmbracelet/lipgloss)** – Styling & layout  
- **[Tablewriter](https://github.com/olekukonko/tablewriter)** – Terminal tables  
- **[x/term](https://pkg.go.dev/golang.org/x/term)** – Terminal size detection  
- **GitHub REST API** – Repository, commits, issues, contributors  

---

## 📝 Project Overview

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

### 🔄 Workflow

1. User launches `repo-lyzer`  
2. Interactive menu → **Analyze** or **Compare**  
3. GitHub API fetch (repos, commits, contributors, languages)  
4. Metrics computed (health, bus factor, maturity)  
5. Displayed in **centered, styled terminal dashboard**

---

## 🏗️ Architecture Overview

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

## 📚 Documentation

### For Contributors
- [ARCHITECTURE.md](docs/ARCHITECTURE.md) – Complete architecture guide  
- [ANALYZER_INTEGRATION.md](docs/ANALYZER_INTEGRATION.md) – Adding new analyzers  
- [IMPLEMENTATION_DETAILS.md](docs/IMPLEMENTATION_DETAILS.md) – Technical deep dive  

### Reference
- [DOCUMENTATION_INDEX.md](docs/DOCUMENTATION_INDEX.md) – Master index  
- [QUICK_REFERENCE.md](docs/QUICK_REFERENCE.md) – Quick usage guide  
- [CHANGE_LOG.md](docs/CHANGE_LOG.md) – Version history  

---

## 🚧 Challenges Faced

- Centering multi-section layouts in terminal  
- Handling high-activity repositories  
- GitHub API rate limits  
- Horizontal commit graphs  
- Unified responsive dashboard  

---

## 🛠 Installation (From Source)

```bash
git clone https://github.com/agnivo988/Repo-lyzer.git
cd Repo-lyzer
go mod tidy
go run main.go
```

---

## ▶️ Usage

### Analyze Repository
```bash
repo-lyzer analyze golang/go
```

### Compare Repositories
Available via interactive menu.

### Export Results
Export to JSON or Markdown from dashboard.

---
  <img src="https://res.cloudinary.com/dhyii4oiw/image/upload/v1767290545/Screenshot_2026-01-01_224310_c0hhr8.png" width="90%">
  <img src="https://res.cloudinary.com/dhyii4oiw/image/upload/v1767324721/Screenshot_2026-01-02_090050_u6xweq.png" width="90%">
  <img src="https://res.cloudinary.com/dhyii4oiw/image/upload/v1767324721/Screenshot_2026-01-02_090043_keqfs4.png" width="90%">
  <img src="https://res.cloudinary.com/dhyii4oiw/image/upload/v1767324721/Screenshot_2026-01-02_090104_dm7bgk.png" width="90%">
  <img src="https://res.cloudinary.com/dhyii4oiw/image/upload/v1767324829/Screenshot_2026-01-02_090335_acms5i.png" width="90%">
</p>

---

## 📦 Installation (For Users)

```bash
go install github.com/agnivo988/Repo-lyzer@v1.0.5
repo-lyzer
```

---

## 📄 License

MIT License © 2026 Agniva Mukherjee
yaml
Copy code
