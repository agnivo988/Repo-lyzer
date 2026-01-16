# API Reference Documentation

This document provides a comprehensive guide to the internal APIs used within the Repo-lyzer application. It covers the GitHub API client functions, analyzer modules, UI components, and key data structures. This documentation is intended for developers integrating or extending the codebase.

## Table of Contents

- [GitHub API Client](#github-api-client)
- [Analyzer Modules](#analyzer-modules)
- [UI Components](#ui-components)
- [Data Structures](#data-structures)

## GitHub API Client

The GitHub API client (`internal/github`) provides functions to interact with the GitHub API for retrieving repository data.

### NewClient()

Creates a new GitHub API client instance.

**Signature:**
```go
func NewClient() *Client
```

**Parameters:**
- None

**Returns:**
- `*Client`: A pointer to a new Client instance

**Notes:**
- Reads the authentication token from the `GITHUB_TOKEN` environment variable
- Callers must set the `GITHUB_TOKEN` environment variable prior to constructing a client
- Uses the value of `GITHUB_TOKEN` as the API authentication token

**Example:**
```go
client := github.NewClient()
```

### HasToken()

Returns true if a GitHub token is configured.

**Signature:**
```go
func (c *Client) HasToken() bool
```

**Parameters:**
- None

**Returns:**
- `bool`: True if a GitHub token is configured

**Example:**
```go
if client.HasToken() {
    fmt.Println("GitHub token is configured")
}
```

### GetUser()

Fetches the authenticated user.

**Signature:**
```go
func (c *Client) GetUser() (*User, error)
```

**Parameters:**
- None

**Returns:**
- `*User`: User information
- `error`: Error if the request fails

**Example:**
```go
user, err := client.GetUser()
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Authenticated as: %s\n", user.Login)
```

### GetFileContent()

Fetches the content of a file from a repository. Returns the base64 encoded content.

**Signature:**
```go
func (c *Client) GetFileContent(owner, repo, path string) (string, error)
```

**Parameters:**
- `owner` (string): The repository owner's username
- `repo` (string): The repository name
- `path` (string): Path to the file in the repository

**Returns:**
- `string`: Base64 encoded file content
- `error`: Error if the request fails

**Example:**
```go
content, err := client.GetFileContent("octocat", "Hello-World", "README.md")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("File content: %s\n", content)
```

### GetRepo()

Retrieves repository information from GitHub.

**Signature:**
```go
func (c *Client) GetRepo(owner, repo string) (*Repo, error)
```

**Parameters:**
- `owner` (string): The repository owner's username
- `repo` (string): The repository name

**Returns:**
- `*Repo`: Repository information
- `error`: Error if the request fails

**Example:**
```go
repo, err := client.GetRepo("octocat", "Hello-World")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Repository: %s\n", repo.Name)
```

### GetContributors()

Fetches all contributors for a repository, handling pagination automatically.

**Signature:**
```go
func (c *Client) GetContributors(owner, repo string) ([]Contributor, error)
```

**Parameters:**
- `owner` (string): The repository owner's username
- `repo` (string): The repository name

**Returns:**
- `[]Contributor`: Slice of contributors with their commit counts
- `error`: Error if the request fails

**Example:**
```go
contributors, err := client.GetContributors("octocat", "Hello-World")
if err != nil {
    log.Fatal(err)
}
for _, contributor := range contributors {
    fmt.Printf("%s: %d commits\n", contributor.Login, contributor.Commits)
}
```

### GetCommits()

Retrieves commits for a repository within the specified number of days.

**Signature:**
```go
func (c *Client) GetCommits(owner, repo string, days int) ([]Commit, error)
```

**Parameters:**
- `owner` (string): The repository owner's username
- `repo` (string): The repository name
- `days` (int): Number of days to look back for commits

**Returns:**
- `[]Commit`: Slice of commits
- `error`: Error if the request fails

**Example:**
```go
commits, err := client.GetCommits("octocat", "Hello-World", 365)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Found %d commits in the last year\n", len(commits))
```

### GetLanguages()

Retrieves the programming languages used in a repository.

**Signature:**
```go
func (c *Client) GetLanguages(owner, repo string) (map[string]int, error)
```

**Parameters:**
- `owner` (string): The repository owner's username
- `repo` (string): The repository name

**Returns:**
- `map[string]int`: Map of language names to bytes of code
- `error`: Error if the request fails

**Example:**
```go
languages, err := client.GetLanguages("octocat", "Hello-World")
if err != nil {
    log.Fatal(err)
}
for lang, bytes := range languages {
    fmt.Printf("%s: %d bytes\n", lang, bytes)
}
```

### GetFileTree()

Retrieves the file tree for a repository at a specific branch/commit.

**Signature:**
```go
func (c *Client) GetFileTree(owner, repo, branch string) ([]TreeEntry, error)
```

**Parameters:**
- `owner` (string): The repository owner's username
- `repo` (string): The repository name
- `branch` (string): Branch name or commit SHA

**Returns:**
- `[]TreeEntry`: Slice of tree entries (files and directories)
- `error`: Error if the request fails

**Notes:**
- Uses recursive=1 to get the full tree

**Example:**
```go
tree, err := client.GetFileTree("octocat", "Hello-World", "main")
if err != nil {
    log.Fatal(err)
}
for _, entry := range tree {
    fmt.Printf("%s (%s)\n", entry.Path, entry.Type)
}
```

### GetRateLimit()

Fetches current rate limit status from GitHub API.

**Signature:**
```go
func (c *Client) GetRateLimit() (*RateLimit, error)
```

**Parameters:**
- None

**Returns:**
- `*RateLimit`: Rate limit information
- `error`: Error if the request fails

**Example:**
```go
rateLimit, err := client.GetRateLimit()
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Remaining requests: %d/%d\n", rateLimit.Resources.Core.Remaining, rateLimit.Resources.Core.Limit)
```

### GetIssues()

Retrieves issues for a repository with a specific state.

**Signature:**
```go
func (c *Client) GetIssues(owner, repo string, state string) ([]Issue, error)
```

**Parameters:**
- `owner` (string): The repository owner's username
- `repo` (string): The repository name
- `state` (string): Issue state ("open", "closed", or "all")

**Returns:**
- `[]Issue`: Slice of issues
- `error`: Error if the request fails

**Example:**
```go
issues, err := client.GetIssues("octocat", "Hello-World", "open")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Found %d open issues\n", len(issues))
```

## Analyzer Modules

The analyzer modules (`internal/analyzer`) provide functions for analyzing GitHub repository data and computing various metrics.

### CalculateHealth()

Calculates a health score for a repository based on various metrics.

**Signature:**
```go
func CalculateHealth(repo *github.Repo, commits []github.Commit) int
```

**Parameters:**
- `repo` (*github.Repo): Repository information
- `commits` ([]github.Commit): Recent commits

**Returns:**
- `int`: Health score (0-100)

**Notes:**
- Considers factors like description presence, star count, commit count, and open issues

**Example:**
```go
health := analyzer.CalculateHealth(repo, commits)
fmt.Printf("Repository health: %d/100\n", health)
```

### BusFactor()

Calculates the bus factor of a repository based on contributor commit distribution.

**Signature:**
```go
func BusFactor(contributors []github.Contributor) (int, string)
```

**Parameters:**
- `contributors` ([]github.Contributor): Slice of repository contributors with their commit counts

**Returns:**
- `int`: Risk score (1=High Risk, 2=Medium Risk, 3=Low Risk)
- `string`: Risk level description

**Example:**
```go
contributors := []github.Contributor{
    {Login: "alice", Commits: 100},
    {Login: "bob", Commits: 50},
    {Login: "charlie", Commits: 25},
}
score, risk := analyzer.BusFactor(contributors)
// score: 2, risk: "Medium Risk"
fmt.Printf("Bus Factor: %d (%s)\n", score, risk)
```

### AnalyzeCodeQuality()

Performs comprehensive code quality analysis of a repository.

**Signature:**
```go
func AnalyzeCodeQuality(repo *github.Repo, fileTree []github.TreeEntry, languages map[string]int) *CodeQualityMetrics
```

**Parameters:**
- `repo` (*github.Repo): Repository information
- `fileTree` ([]github.TreeEntry): File tree entries
- `languages` (map[string]int): Language usage statistics

**Returns:**
- `*CodeQualityMetrics`: Comprehensive code quality analysis results

**Notes:**
- Analyzes documentation, testing, structure, maintenance, and detects code smells
- Provides recommendations for improvement

**Example:**
```go
metrics := analyzer.AnalyzeCodeQuality(repo, fileTree, languages)
fmt.Printf("Code Quality Grade: %s (%d/100)\n", metrics.Grade, metrics.OverallScore)
for _, rec := range metrics.Recommendations {
    fmt.Printf("Recommendation: %s\n", rec)
}
```

## UI Components

The UI components (`internal/ui`) provide the terminal-based user interface using the Bubble Tea framework.

### Run()

Starts the main Repo-lyzer application with the interactive terminal interface.

**Signature:**
```go
func Run() error
```

**Parameters:**
- None

**Returns:**
- `error`: Error if the application fails to start

**Example:**
```go
err := ui.Run()
if err != nil {
    log.Fatal(err)
}
```

### NewMenuModel()

Creates a new menu model with cursor at the first option.

**Signature:**
```go
func NewMenuModel() MenuModel
```

**Parameters:**
- None

**Returns:**
- `MenuModel`: Initialized menu model

**Example:**
```go
menu := ui.NewMenuModel()
```

### NewTreeModel()

Creates a new tree model for displaying repository file structure.

**Signature:**
```go
func NewTreeModel(result *AnalysisResult) TreeModel
```

**Parameters:**
- `result` (*AnalysisResult): Analysis result containing file tree data

**Returns:**
- `TreeModel`: Initialized tree model

**Example:**
```go
tree := ui.NewTreeModel(analysisResult)
```

### NewResponsiveLayout()

Creates a new responsive layout manager for terminal UI.

**Signature:**
```go
func NewResponsiveLayout(width, height int) *ResponsiveLayout
```

**Parameters:**
- `width` (int): Terminal width
- `height` (int): Terminal height

**Returns:**
- `*ResponsiveLayout`: Responsive layout manager

**Example:**
```go
layout := ui.NewResponsiveLayout(120, 30)
```

### ApplyTheme()

Updates all style variables with the given theme.

**Signature:**
```go
func ApplyTheme(theme Theme)
```

**Parameters:**
- `theme` (Theme): Theme configuration

**Returns:**
- None

**Example:**
```go
ui.ApplyTheme(ui.GetThemeByName("dark"))
```

### CycleTheme()

Switches to the next theme in the available themes list.

**Signature:**
```go
func CycleTheme() Theme
```

**Parameters:**
- None

**Returns:**
- `Theme`: The new current theme

**Example:**
```go
newTheme := ui.CycleTheme()
```

### LoadHistory()

Loads analysis history from file.

**Signature:**
```go
func LoadHistory() (*History, error)
```

**Parameters:**
- None

**Returns:**
- `*History`: History data
- `error`: Error if loading fails

**Example:**
```go
history, err := ui.LoadHistory()
if err != nil {
    log.Fatal(err)
}
```

### NewProgressTracker()

Creates a tracker with default analysis stages.

**Signature:**
```go
func NewProgressTracker() *ProgressTracker
```

**Parameters:**
- None

**Returns:**
- `*ProgressTracker`: Progress tracker instance

**Example:**
```go
tracker := ui.NewProgressTracker()
```

### BuildFileTree()

Creates a file tree from repository content.

**Signature:**
```go
func BuildFileTree(result AnalysisResult) *FileNode
```

**Parameters:**
- `result` (AnalysisResult): Analysis result

**Returns:**
- `*FileNode`: Root of the file tree

**Example:**
```go
root := ui.BuildFileTree(analysisResult)
```

## Data Structures

### Repo

Represents a GitHub repository.

**Fields:**
- `Name` (string): Repository name
- `FullName` (string): Full repository name (owner/repo)
- `Stars` (int): Number of stars
- `Forks` (int): Number of forks
- `OpenIssues` (int): Number of open issues
- `Description` (string): Repository description
- `CreatedAt` (time.Time): Creation date
- `UpdatedAt` (time.Time): Last update date
- `PushedAt` (time.Time): Last push date
- `Language` (string): Primary programming language
- `Fork` (bool): Whether this is a fork
- `Archived` (bool): Whether the repository is archived
- `Private` (bool): Whether the repository is private
- `DefaultBranch` (string): Default branch name
- `HTMLURL` (string): GitHub URL
- `CloneURL` (string): Clone URL

### Contributor

Represents a GitHub repository contributor.

**Fields:**
- `Login` (string): Contributor's GitHub username
- `Commits` (int): Number of commits made by the contributor

### Commit

Represents a GitHub commit.

**Fields:**
- `SHA` (string): Commit SHA hash
- `Commit.Author.Date` (time.Time): Commit author date

### TreeEntry

Represents a file or directory entry in a GitHub repository tree.

**Fields:**
- `Path` (string): Path to the file or directory
- `Mode` (string): File mode/permissions
- `Type` (string): Type of entry ("blob" for files, "tree" for directories)
- `Size` (int): Size of the file in bytes (0 for directories)
- `Sha` (string): SHA hash of the file or tree

### User

Represents a GitHub user.

**Fields:**
- `Login` (string): GitHub username
- `Name` (string): Full name

### Issue

Represents a GitHub issue.

**Fields:**
- `State` (string): Issue state ("open" or "closed")

### RateLimit

Represents GitHub API rate limit information.

**Fields:**
- `Resources.Core.Limit` (int): Maximum requests per hour
- `Resources.Core.Remaining` (int): Remaining requests
- `Resources.Core.Reset` (int): Unix timestamp when limit resets
- `Resources.Core.Used` (int): Requests used

**Methods:**
- `ResetTime()`: Returns reset time as time.Time
- `TimeUntilReset()`: Returns duration until reset
- `IsLimited()`: Returns true if rate limited
- `UsagePercent()`: Returns percentage of limit used
- `FormatResetTime()`: Returns human-readable reset time
- `GetRateLimitStatus()`: Returns formatted status string

### CodeQualityMetrics

Contains comprehensive code quality analysis results.

**Fields:**
- `OverallScore` (int): Overall quality score (0-100)
- `Grade` (string): Letter grade (A, B, C, D, F)
- `DocumentationScore` (int): Documentation quality score
- `TestingScore` (int): Testing quality score
- `StructureScore` (int): Code structure score
- `MaintenanceScore` (int): Maintenance quality score
- `HasReadme` (bool): Whether README exists
- `HasContributing` (bool): Whether CONTRIBUTING guide exists
- `HasLicense` (bool): Whether LICENSE exists
- `HasChangelog` (bool): Whether CHANGELOG exists
- `HasCodeOfConduct` (bool): Whether code of conduct exists
- `HasTests` (bool): Whether tests exist
- `HasCI` (bool): Whether CI/CD is configured
- `HasDocker` (bool): Whether Docker is configured
- `HasEditorConfig` (bool): Whether .editorconfig exists
- `HasGitignore` (bool): Whether .gitignore exists
- `TestFrameworks` ([]string): Detected test frameworks
- `CIProviders` ([]string): Detected CI providers
- `FileStats` (FileStatistics): File statistics
- `CodeSmells` ([]CodeSmell): Detected code quality issues
- `Recommendations` ([]string): Improvement recommendations

### FileStatistics

Contains file-related metrics from code quality analysis.

**Fields:**
- `TotalFiles` (int): Total number of files
- `SourceFiles` (int): Number of source code files
- `TestFiles` (int): Number of test files
- `DocFiles` (int): Number of documentation files
- `ConfigFiles` (int): Number of configuration files
- `TestRatio` (float64): Ratio of test files to source files
- `AvgPathDepth` (float64): Average directory depth
- `FilesByExtension` (map[string]int): Files grouped by extension
- `LargestFiles` ([]string): Files with deepest paths

### CodeSmell

Represents a detected code quality issue.

**Fields:**
- `Type` (string): Type of code smell
- `Severity` (string): Severity level ("Low", "Medium", "High")
- `Description` (string): Description of the issue
- `Location` (string): Location of the issue (optional)

### History

Manages analysis history data.

**Fields:**
- `Entries` ([]HistoryEntry): List of history entries

**Methods:**
- `Save()`: Saves history to file
- `AddEntry(data AnalysisResult)`: Adds new entry
- `GetRecent(count int)`: Returns recent entries
- `Clear()`: Removes all entries
- `Delete(index int)`: Removes specific entry
- `SortByDate()`: Sorts entries by date

### HistoryEntry

Represents a single history entry.

**Fields:**
- `Timestamp` (time.Time): When analysis was performed
- `Repository` (string): Repository name
- `AnalysisResult` (AnalysisResult): Analysis results

**Methods:**
- `Format()`: Returns formatted display string

### ProgressTracker

Tracks analysis progress through stages.

**Fields:**
- `stages` ([]ProgressStage): Analysis stages
- `current` (int): Current stage index
- `startTime` (time.Time): When analysis started

**Methods:**
- `NextStage()`: Advances to next stage
- `GetCurrentStage()`: Returns current stage info
- `GetAllStages()`: Returns all stages
- `GetProgress()`: Returns completion percentage
- `GetProgressBar(width int)`: Returns visual progress bar
- `GetElapsedTime()`: Returns time elapsed

### ResponsiveLayout

Manages responsive terminal UI layout.

**Fields:**
- `Width` (int): Terminal width
- `Height` (int): Terminal height

**Methods:**
- `IsSmallTerminal()`: Returns true if terminal is very small
- `IsMobileTerminal()`: Returns true if terminal is mobile-sized
- `GetMaxContentWidth()`: Returns safe content width
- `GetMaxContentHeight()`: Returns safe content height
- `CenterText(text string)`: Centers text horizontally and vertically
- `CenterContent(content string)`: Centers content with margin
- `WrapText(text string, padding int)`: Wraps text to fit width
- `FormatMenuForDisplay(items []string)`: Formats menu items for display
- `GetMinimumWarning()`: Returns warning for too-small terminal
- `PadContent(content string, horizontal, vertical int)`: Adds padding to content
- `RenderResponsiveBox(title, content string)`: Renders responsive box
- `ShouldShowSidebar()`: Returns whether sidebar should be shown
- `ShouldShowPreview()`: Returns whether preview should be shown
- `GetLayoutMode()`: Returns current layout mode
- `AdjustSpacing()`: Returns appropriate spacing

### MainModel

The main model for the Bubble Tea application, managing the application state and UI components.

**Key Fields:**
- `state` (sessionState): Current application state
- `menu` (MenuModel): Menu component
- `dashboard` (DashboardModel): Dashboard component
- `tree` (TreeModel): File tree component
- `settings` (SettingsModel): Settings component
- `help` (HelpModel): Help component
- `history` (HistoryModel): History component
- `windowWidth` (int): Terminal window width
- `windowHeight` (int): Terminal window height
