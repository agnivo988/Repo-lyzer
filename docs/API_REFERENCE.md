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

## Analyzer Modules

The analyzer modules (`internal/analyzer`) provide functions for analyzing GitHub repository data and computing various metrics.

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
