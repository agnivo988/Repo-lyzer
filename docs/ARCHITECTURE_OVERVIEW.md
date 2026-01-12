# Architecture Overview

This document provides a high-level overview of Repo-lyzer's modular architecture, including diagrams and explanations of how key components interact.

## High-Level Architecture Diagram

```
┌─────────────────┐
│     main.go     │
│   Entry Point   │
└─────────┬───────┘
          │
          ▼
┌─────────────────┐     ┌─────────────────┐
│      cmd/       │◄────┤   internal/ui/  │
│   CLI Commands  │     │   TUI Interface │
│                 │     │                 │
│ • root.go       │     │ • app.go        │
│ • analyze.go    │     │ • dashboard.go  │
│ • compare.go    │     │ • menu.go       │
└─────────┬───────┘     └─────────────────┘
          │
          ▼
┌─────────────────┐     ┌─────────────────┐
│ internal/github/│◄────┤ internal/analyzer│
│   API Client    │     │   Analysis Logic │
│                 │     │                 │
│ • client.go     │     │ • health.go     │
│ • repo.go       │     │ • bus_factor.go │
│ • commits.go    │     │ • maturity.go   │
└─────────────────┘     └─────────────────┘
          │
          ▼
┌─────────────────┐
│ internal/output/│
│   Data Output   │
│                 │
│ • json.go       │
│ • tables.go     │
│ • charts.go     │
└─────────────────┘
```

## Component Interactions

### 1. Entry Point (main.go)
- **Purpose**: Initializes the application and starts the interactive menu.
- **Interactions**:
  - Calls `cmd.RunMenu()` to start the CLI interface.
  - Serves as the single entry point for the entire application.

### 2. Command Layer (cmd/)
- **Purpose**: Handles CLI command parsing and routing using Cobra.
- **Key Files**:
  - `root.go`: Defines the root command and global flags.
  - `analyze.go`: Implements the analyze command logic.
  - `compare.go`: Implements the compare command logic.
- **Interactions**:
  - Receives user input from the terminal.
  - Orchestrates calls to internal packages for data fetching and analysis.
  - Passes results to the UI layer for display.

### 3. User Interface Layer (internal/ui/)
- **Purpose**: Provides the terminal-based user interface using Bubble Tea and Lipgloss.
- **Key Components**:
  - `app.go`: Main TUI application structure.
  - `dashboard.go`: Renders the analysis dashboard.
  - `menu.go`: Handles interactive menu navigation.
- **Interactions**:
  - Receives analysis results from the cmd layer.
  - Renders interactive dashboards and menus.
  - Handles user input for navigation and selections.

### 4. GitHub API Layer (internal/github/)
- **Purpose**: Manages all interactions with the GitHub REST API.
- **Key Files**:
  - `client.go`: HTTP client setup and authentication.
  - `repo.go`: Repository data fetching.
  - `commits.go`: Commit history and activity data.
  - `contributors.go`: Contributor information.
- **Interactions**:
  - Fetches raw data from GitHub API endpoints.
  - Handles authentication, rate limiting, and error responses.
  - Provides structured data to the analyzer layer.

### 5. Analysis Layer (internal/analyzer/)
- **Purpose**: Contains all business logic for analyzing repository data.
- **Key Modules**:
  - `health.go`: Calculates repository health scores.
  - `bus_factor.go`: Assesses contributor risk.
  - `maturity.go`: Evaluates repository maturity.
  - `security.go`: Performs security analysis.
- **Interactions**:
  - Receives raw data from the GitHub layer.
  - Applies scoring algorithms and metrics.
  - Returns structured analysis results to the command layer.

### 6. Output Layer (internal/output/)
- **Purpose**: Handles data formatting and presentation.
- **Key Files**:
  - `json.go`: JSON export functionality.
  - `tables.go`: Terminal table rendering.
  - `charts.go`: Chart and graph generation.
- **Interactions**:
  - Receives analysis results from various layers.
  - Formats data for different output types (JSON, tables, charts).
  - Provides styled output for the UI layer.

## Data Flow

1. **User Input**: User runs `repo-lyzer` command.
2. **Command Processing**: `cmd/` layer parses arguments and initiates analysis.
3. **Data Fetching**: `internal/github/` retrieves repository data from GitHub API.
4. **Analysis**: `internal/analyzer/` processes data and calculates metrics.
5. **Output Formatting**: `internal/output/` prepares data for display.
6. **UI Rendering**: `internal/ui/` presents results in interactive terminal interface.

## Design Principles

### Modularity
- Each package has a single responsibility.
- Clear separation between API client, analysis logic, and UI rendering.
- Easy to test individual components in isolation.

### Dependency Direction
- Dependencies flow inward: UI depends on business logic, which depends on data access.
- No circular dependencies between packages.

### Interface Segregation
- Components communicate through well-defined interfaces.
- Easy to swap implementations (e.g., different API clients or output formats).

### Scalability
- Modular design allows for easy addition of new analyzers or output formats.
- Separation of concerns enables parallel development by multiple contributors.

## Key Design Decisions

### Why Cobra for CLI?
- Industry-standard Go CLI library.
- Automatic help generation and command parsing.
- Easy to extend with new commands.

### Why Bubble Tea for TUI?
- Declarative UI model fits well with Go's philosophy.
- Excellent performance for terminal applications.
- Rich ecosystem of components and styling options.

### Why Modular Analyzer Structure?
- Allows users to enable/disable specific analyses.
- Easy to add new scoring algorithms.
- Separates data fetching from analysis logic.

### Why Separate Output Layer?
- Supports multiple output formats (JSON, tables, charts).
- Decouples presentation from business logic.
- Enables headless usage for automation.

This architecture ensures Repo-lyzer remains maintainable, extensible, and performant as the project grows.
