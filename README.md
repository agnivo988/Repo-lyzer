<h1 align="center">Repo-lyzer</h1>
<p align="center">
  <img src="https://res.cloudinary.com/dhyii4oiw/image/upload/v1767324445/Screenshot_2026-01-02_085503_ros5gz.png" alt="Repo-lyzer Logo" width="300">
</p>

**Repo-lyzer** is a modern, terminal-based CLI tool written in **Golang** that analyzes GitHub repositories and presents insights in a beautifully formatted, interactive dashboard.

---

## Features
- **Deep Analytics** – Repository health, maturity scores, and bus factor.
- **Contribution Friendliness** – Assess how easy it is to contribute to a repository using weighted metrics (`--contribute` flag).
- **Interactive TUI** – Fully navigable keyboard-driven menu system.
- **Visual Data** – Language breakdown bars and horizontal commit graphs.
- **File Explorer** – Browse repository structures directly in the dashboard.
- **Multi-Format Export** – Save reports as JSON, Markdown, CSV, or HTML.

---

## Quick Start

### Installation
```bash
go install [github.com/agnivo988/Repo-lyzer@v1.0.6](https://github.com/agnivo988/Repo-lyzer@v1.0.6)
repo-lyzer
```

### Basic Usage
```bash
# Get a 5-line quick summary
repo-lyzer summary golang/go

# Run full interactive analysis
repo-lyzer analyze microsoft/vscode

# Run analysis with contribution scoring enabled
# Run analysis with contribution scoring enabled
repo-lyzer analyze microsoft/vscode --contribute
```

### Docker Usage

You can run `repo-lyzer` using Docker without installing Go. The Docker image uses a non-root user and is optimized for production.

```bash
# Build the image
docker build -t repo-lyzer .

# Run the CLI interactively
docker run -it --rm repo-lyzer
```

### Docker Compose (Daemon Mode)

For continuous monitoring and scheduling, you can run the daemon mode using `docker-compose`:

```bash
# Start the daemon
docker compose up -d

# View logs
docker compose logs -f
```

#### Environment Variables
You can configure `repo-lyzer` using environment variables. These will override file-based settings:
- `REPO_LYZER_GITHUB_TOKEN`: Your GitHub Personal Access Token
- `REPO_LYZER_INTERVAL`: Scheduler polling interval (e.g., `30s`, `5m`, `1h`)
- `REPO_LYZER_LOG_LEVEL`: Logging level (`debug`, `info`, `warn`, `error`)
- `REPO_LYZER_CONFIG_PATH`: Override the config file path (defaults to `/app/data/settings.json` in the container)

The `docker-compose.yml` mounts a local `./data` directory to persist settings and reports.

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

---

## Documentation

### For Contributors
- [ARCHITECTURE.md](docs/ARCHITECTURE.md) – Complete architecture guide  
- [ANALYZER_INTEGRATION.md](docs/ANALYZER_INTEGRATION.md) – Adding new analyzers  
- [IMPLEMENTATION_DETAILS.md](docs/IMPLEMENTATION_DETAILS.md) – Technical deep dive
- [PROJECT STRUCTURE.md](docs/PROJECT_STRUCTURE.md) - Project Structure and Workflow

### Reference
- [DOCUMENTATION_INDEX.md](docs/DOCUMENTATION_INDEX.md) – Master index  
- [QUICK_REFERENCE.md](docs/QUICK_REFERENCE.md) – Quick usage guide  
- [CHANGE_LOG.md](docs/CHANGE_LOG.md) – Version history  

---

## Maintainers & Contributors
### Maintainer: @agnivo988

<a href="https://github.com/Aamod007"><img src="https://github.com/Aamod007.png" width="50" height="50" alt="Aamod007" title="Contributor"></a> <a href="https://github.com/Aditya8369"><img src="https://github.com/Aditya8369.png" width="50" height="50" alt="Aditya8369" title="Contributor"></a> <a href="https://github.com/agnivo988"><img src="https://github.com/agnivo988.png" width="50" height="50" alt="agnivo988" title="Project Maintainer"></a> <a href="https://github.com/Gupta-02"><img src="https://github.com/Gupta-02.png" width="50" height="50" alt="Gupta-02" title="Contributor"></a> <a href="https://github.com/GauravKarakoti"><img src="https://github.com/GauravKarakoti.png" width="50" height="50" alt="GauravKarakoti" title="Contributor"></a> <a href="https://github.com/Sappymukherjee214"><img src="https://github.com/Sappymukherjee214.png" width="50" height="50" alt="Sappymukherjee214" title="Contributor"></a> <a href="https://github.com/ItsMeArm00n"><img src="https://github.com/ItsMeArm00n.png" width="50" height="50" alt="ItsMeArm00n" title="Contributor"></a> <a href="https://github.com/MuktaRedij"><img src="https://github.com/MuktaRedij.png" width="50" height="50" alt="MuktaRedij" title="Contributor"></a> <a href="https://github.com/Kiran95021"><img src="https://github.com/Kiran95021.png" width="50" height="50" alt="Kiran95021" title="Contributor"></a> <a href="https://github.com/Shriii19"><img src="https://github.com/Shriii19.png" width="50" height="50" alt="Shriii19" title="Contributor"></a> <a href="https://github.com/KUMARI-SONALIUPADHYAY"><img src="https://github.com/KUMARI-SONALIUPADHYAY.png" width="50" height="50" alt="KUMARI-SONALIUPADHYAY" title="Contributor"></a> <a href="https://github.com/magic-peach"><img src="https://github.com/magic-peach.png" width="50" height="50" alt="magic-peach" title="Contributor"></a> <a href="https://github.com/coderabbitai"><img src="https://github.com/coderabbitai.png" width="50" height="50" alt="coderabbitai[bot]" title="Bot Contributor"></a> <a href="https://github.com/sahoo-tech"><img src="https://github.com/sahoo-tech.png" width="50" height="50" alt="sahoo-tech" title="Contributor"></a> <a href="https://github.com/Abhijeet-980"><img src="https://github.com/Abhijeet-980.png" width="50" height="50" alt="Abhijeet-980" title="Contributor"></a> <a href="https://github.com/Diksha78-bot"><img src="https://github.com/Diksha78-bot.png" width="50" height="50" alt="Diksha78-bot" title="Bot Contributor"></a> <a href="https://github.com/Shivani-Meena07"><img src="https://github.com/Shivani-Meena07.png" width="50" height="50" alt="Shivani-Meena07" title="Contributor"></a> <a href="https://github.com/ShashankSaga"><img src="https://github.com/ShashankSaga.png" width="50" height="50" alt="ShashankSaga" title="Contributor"></a>



---

## License
**MIT License © 2026 Agniva Mukherjee**
