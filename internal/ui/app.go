package ui

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/agnivo988/Repo-lyzer/internal/analyzer"
	"github.com/agnivo988/Repo-lyzer/internal/cache"
	"github.com/agnivo988/Repo-lyzer/internal/config"
	"github.com/agnivo988/Repo-lyzer/internal/github"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type sessionState int

const (
	stateMenu sessionState = iota
	stateInput
	stateLoading
	stateDashboard
	stateTree
	stateFileEdit
	stateSettings
	stateHelp
	stateHistory
	stateFavorites
	stateCompareInput
	stateCompareLoading
	stateCompareResult
	stateCloneInput
	stateCloning
)

// Message types for sub-models
type AnalyzeRepoMsg struct {
	repoName string
}

type BackToMenuMsg struct{}

type SwitchToInputMsg struct{}

type ErrorMsg error

// StatusMsg represents a status message with error indication
type StatusMsg struct {
	Message string
	IsError bool
}

func (s StatusMsg) Error() string {
	return s.Message
}

type MainModel struct {
	state sessionState

	// Sub-models for different UI states
	menu           MenuModel
	input          InputModel
	loading        LoadingModel
	compareInput   CompareInputModel
	compareLoading CompareLoadingModel
	compareResult  CompareResultModel
	settings       SettingsModel
	help           HelpModel
	history        HistoryModel
	favorites      *FavoritesModel
	cloneInput     CloneInputModel
	cloning        CloningModel

	// Shared models
	dashboard DashboardModel
	tree      TreeModel
	fileEdit  FileEditModel

	// Shared state
	windowWidth  int
	windowHeight int
	cache        *cache.Cache
	appConfig    *config.AppSettings

	// Additional fields used in Update method
	spinner         spinner.Model
	historyCursor   int
	favoritesCursor int
	animTick        int
	err             interface{}
	analysisType    string
	compareStep     int
	compareInput1   string
	compareInput2   string
	inTokenInput    bool
	tokenInput      string
	settingsOption  string
	helpContent     string
	repoInput       string
	progress        *ProgressTracker
	cacheStatus     string
}

// NewMainModel creates a new MainModel with initialized sub-models
func NewMainModel(cache *cache.Cache, config *config.AppSettings) MainModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return MainModel{
		state:          stateMenu,
		menu:           NewMenuModel(),
		input:          NewInputModel(),
		loading:        NewLoadingModel(),
		compareInput:   NewCompareInputModel(),
		compareLoading: NewCompareLoadingModel(),
		compareResult:  NewCompareResultModel(),
		settings:       NewSettingsModel(),
		help:           NewHelpModel(),
		history:        NewHistoryModel(),
		favorites:      NewFavoritesModel(),
		cloneInput:     NewCloneInputModel(),
		cloning:        NewCloningModel(),
		dashboard:      NewDashboardModel(),
		tree:           NewTreeModel(nil),
		cache:          cache,
		appConfig:      config,
		spinner:        s,
	}
}

// Init initializes the Bubble Tea program
func (m MainModel) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model
func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowWidth = msg.Width
		m.windowHeight = msg.Height
		// Update sub-models with window size
		m.menu.width = msg.Width
		m.menu.height = msg.Height

	case BackToMenuMsg:
		m.state = stateMenu
		return m, nil

	case SwitchToInputMsg:
		m.state = stateInput
		return m, nil

	case AnalyzeRepoMsg:
		m.state = stateLoading
		// TODO: Start analysis command
		return m, nil

	case ErrorMsg:
		m.err = error(msg)
		return m, nil
	}

	// Delegate to current state's sub-model
	switch m.state {
	case stateMenu:
		menuModel, cmd := m.menu.Update(msg)
		m.menu = menuModel.(MenuModel)
		if m.menu.Done {
			switch m.menu.SelectedOption {
			case 0: // Analyze Repository
				m.state = stateInput
			case 1: // Favorites
				m.state = stateFavorites
			case 2: // Compare Repositories
				m.state = stateCompareInput
			case 3: // View History
				m.state = stateHistory
			case 4: // Clone Repository
				m.state = stateCloneInput
			case 5: // Settings
				m.state = stateSettings
			case 6: // Help
				m.state = stateHelp
			case 7: // Exit
				return m, tea.Quit
			}
			m.menu.Done = false
		}
		cmds = append(cmds, cmd)

	case stateInput:
		newInput, cmd := m.input.Update(msg)
		m.input = newInput.(InputModel)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}

		// Handle messages from input model
		switch msg := msg.(type) {
		case AnalyzeRepoMsg:
			m.state = stateLoading
			m.loading.SetRepoName(msg.repoName)
			cmds = append(cmds, m.analyzeRepo(msg.repoName), TickProgressCmd())
		case BackToMenuMsg:
			m.state = stateMenu
		}

	case stateCompareInput:
		newCompareInput, cmd := m.compareInput.Update(msg)
		m.compareInput = newCompareInput.(CompareInputModel)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}

		// Handle messages from compare input model
		switch msg := msg.(type) {
		case CompareReposMsg:
			m.state = stateCompareLoading
			cmds = append(cmds, m.compareRepos(msg.Repo1, msg.Repo2), TickProgressCmd())
		case BackToMenuMsg:
			m.state = stateMenu
		}

	case stateCompareLoading:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)

		switch msg := msg.(type) {
		case CompareResult:
			m.compareResult.result = &msg
			m.state = stateCompareResult
			m.err = nil
		case error:
			m.err = msg
			m.state = stateCompareInput
			m.compareStep = 0
		case tea.KeyMsg:
			if msg.String() == "esc" {
				m.state = stateMenu
				m.compareInput1 = ""
				m.compareInput2 = ""
				m.err = nil
			}
		}

	case stateCompareResult:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "q", "esc":
				m.state = stateMenu
				m.compareResult = CompareResultModel{}
				m.compareInput1 = ""
				m.compareInput2 = ""
			case "j":
				// Export comparison to JSON
				if m.compareResult.result != nil && m.compareResult.result.Repo1.Repo != nil && m.compareResult.result.Repo2.Repo != nil {
					_, err := ExportCompareJSON(*m.compareResult.result)
					if err != nil {
						m.compareResult.err = fmt.Errorf("failed to export JSON: %w", err)
					} else {
						m.compareResult.err = StatusMsg{Message: "✓ Exported comparison to JSON successfully", IsError: false}
					}
				} else {
					m.compareResult.err = fmt.Errorf("no comparison data available for export")
				}
			case "m":
				// Export comparison to Markdown
				if m.compareResult.result != nil && m.compareResult.result.Repo1.Repo != nil && m.compareResult.result.Repo2.Repo != nil {
					_, err := ExportCompareMarkdown(*m.compareResult.result)
					if err != nil {
						m.compareResult.err = fmt.Errorf("failed to export Markdown: %w", err)
					} else {
						m.compareResult.err = StatusMsg{Message: "✓ Exported comparison to Markdown successfully", IsError: false}
					}
				} else {
					m.compareResult.err = fmt.Errorf("no comparison data available for export")
				}
			}
		}

	case stateLoading:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)

		if result, ok := msg.(AnalysisResult); ok {
			m.dashboard.SetData(result)
			m.dashboard.SetCacheStatus("fresh")
			m.state = stateDashboard
			m.progress = nil
			m.cacheStatus = "fresh"
			// Save to history
			history, _ := LoadHistory()
			m.history.Entries = history.Entries
			m.history.AddEntry(result)
			m.history.Save()
		}
		if cachedResult, ok := msg.(CachedAnalysisResult); ok {
			m.dashboard.SetData(cachedResult.Result)
			m.dashboard.SetCacheStatus("cached")
			m.state = stateDashboard
			m.progress = nil
			m.cacheStatus = "cached"
			// Save to history
			history, _ := LoadHistory()
			m.history.Entries = history.Entries
			m.history.AddEntry(cachedResult.Result)
			m.history.Save()
		}
		if err, ok := msg.(error); ok {
			m.err = err
			m.state = stateInput // Go back to input on error
			m.progress = nil
		}

	case stateFavorites:
		switch msg := msg.(type) {
		case StatusMsg:
			m.err = msg
		case tea.KeyMsg:
			switch msg.String() {
			case "up", "k":
				if m.favoritesCursor > 0 {
					m.favoritesCursor--
				}
			case "down", "j":
				if m.favorites.favorites != nil && m.favoritesCursor < len(m.favorites.favorites.Items)-1 {
					m.favoritesCursor++
				}
			case "enter":
				// Analyze selected favorite
				if m.favorites.favorites != nil && len(m.favorites.favorites.Items) > 0 {
					repoName := m.favorites.favorites.Items[m.favoritesCursor].RepoName
					m.favorites.favorites.UpdateUsage(repoName)
					if err := m.favorites.Save(); err != nil {
						log.Printf("Failed to save favorites: %v", err)
						m.err = fmt.Errorf("Failed to save favorites: %v", err)
					} else {
						m.input.input = repoName
						m.state = stateLoading
						cmds = append(cmds, m.analyzeRepo(repoName), TickProgressCmd())
					}
				}
			case "d":
				// Remove from favorites
				if m.favorites.favorites != nil && len(m.favorites.favorites.Items) > 0 {
					m.favorites.favorites.Remove(m.favorites.favorites.Items[m.favoritesCursor].RepoName)
					if err := m.favorites.Save(); err != nil {
						log.Printf("Failed to save favorites: %v", err)
						m.err = fmt.Errorf("Failed to save favorites: %v", err)
					} else {
						if m.favoritesCursor >= len(m.favorites.favorites.Items) && m.favoritesCursor > 0 {
							m.favoritesCursor--
						}
					}
				}
			case "a":
				// Add new favorite (go to input)
				m.state = stateInput
			case "q", "esc":
				m.state = stateMenu
			}
		}

	case stateHistory:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "up", "k":
				if m.historyCursor > 0 {
					m.historyCursor--
				}
			case "down", "j":
				if m.historyCursor < len(m.history.Entries)-1 {
					m.historyCursor++
				}
			case "enter":
				// Re-analyze selected repo
				if len(m.history.Entries) > 0 {
					repoName := m.history.Entries[m.historyCursor].RepoName
					m.input.input = repoName
					m.state = stateLoading
					cmds = append(cmds, m.analyzeRepo(repoName), TickProgressCmd())
				}
			case "d":
				// Delete selected entry
				if len(m.history.Entries) > 0 {
					m.history.Delete(m.historyCursor)
					m.history.Save()
					if m.historyCursor >= len(m.history.Entries) && m.historyCursor > 0 {
						m.historyCursor--
					}
				}
			case "c":
				// Clear all history
				m.history.Clear()
				m.history.Save()
				m.historyCursor = 0
			case "q", "esc":
				m.state = stateMenu
			}
		}

	case stateHelp:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "q", "esc":
				m.state = stateMenu
			}
		}

	case stateSettings:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			// Handle token input mode separately
			if m.inTokenInput {
				switch msg.Type {
				case tea.KeyEnter:
					// Save the token
					if m.appConfig != nil && m.tokenInput != "" {
						m.appConfig.SetGitHubToken(m.tokenInput)
						m.err = fmt.Errorf("✓ GitHub token saved")
					}
					m.inTokenInput = false
					m.tokenInput = ""
				case tea.KeyEsc:
					m.inTokenInput = false
					m.tokenInput = ""
				case tea.KeyBackspace:
					if len(m.tokenInput) > 0 {
						m.tokenInput = m.tokenInput[:len(m.tokenInput)-1]
					}
				case tea.KeyRunes:
					m.tokenInput += string(msg.Runes)
				}
				return m, tea.Batch(cmds...)
			}

			switch msg.String() {
			case "q", "esc":
				m.state = stateMenu
			case "t":
				// Cycle through themes (theme settings)
				if m.settingsOption == "theme" || m.settingsOption == "" {
					theme := CycleTheme()
					if m.appConfig != nil {
						m.appConfig.SetTheme(theme.Name)
					}
					m.err = fmt.Errorf("Theme changed to: %s", theme.Name)
				}
			case "1", "2", "3", "4", "5", "6", "7":
				// Select theme by number (theme settings)
				if m.settingsOption == "theme" || m.settingsOption == "" {
					idx := int(msg.String()[0] - '1')
					if idx >= 0 && idx < len(AvailableThemes) {
						theme := SetThemeByIndex(idx)
						if m.appConfig != nil {
							m.appConfig.SetTheme(theme.Name)
						}
						m.err = fmt.Errorf("Theme: %s", theme.Name)
					}
				}
			case "e":
				// Toggle cache enabled (cache settings)
				if m.settingsOption == "cache" && m.cache != nil {
					cfg := m.cache.GetConfig()
					m.cache.SetEnabled(!cfg.Enabled)
					if cfg.Enabled {
						m.err = fmt.Errorf("Cache disabled")
					} else {
						m.err = fmt.Errorf("Cache enabled")
					}
				}
			case "a":
				// Toggle auto-cache (cache settings)
				if m.settingsOption == "cache" && m.cache != nil {
					cfg := m.cache.GetConfig()
					m.cache.SetAutoCache(!cfg.AutoCache)
					if cfg.AutoCache {
						m.err = fmt.Errorf("Auto-cache disabled")
					} else {
						m.err = fmt.Errorf("Auto-cache enabled")
					}
				}
			case "c":
				// Clear all cache (cache settings) or clear token (token settings)
				if m.settingsOption == "cache" && m.cache != nil {
					m.cache.Clear()
					m.err = fmt.Errorf("Cache cleared")
				} else if m.settingsOption == "token" && m.appConfig != nil {
					m.appConfig.ClearGitHubToken()
					m.err = fmt.Errorf("GitHub token cleared")
				}
			case "x":
				// Clean expired entries (cache settings)
				if m.settingsOption == "cache" && m.cache != nil {
					removed := m.cache.CleanExpired()
					m.err = fmt.Errorf("Removed %d expired entries", removed)
				}
			case "f":
				// Cycle export format (export settings)
				if m.settingsOption == "export" && m.appConfig != nil {
					newFormat := m.appConfig.CycleExportFormat()
					m.err = fmt.Errorf("Export format: %s", newFormat.DisplayName())
				}
			case "i":
				// Enter token input mode (token settings)
				if m.settingsOption == "token" {
					m.inTokenInput = true
					m.tokenInput = ""
				}
			case "y":
				// Confirm reset (reset settings)
				if m.settingsOption == "reset" {
					newSettings, err := config.ResetToDefaults()
					if err == nil {
						m.appConfig = newSettings
						SetThemeByName(newSettings.ThemeName)
						m.err = fmt.Errorf("✓ All settings reset to defaults")
					} else {
						m.err = fmt.Errorf("Failed to reset: %v", err)
					}
					m.state = stateMenu
				}
			}
		}

	case stateCloneInput:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				if m.cloneInput.input != "" {
					m.state = stateCloning
					cmds = append(cmds, m.cloneRepo(m.cloneInput.input))
				}
			case "esc":
				m.state = stateMenu
				m.cloneInput.input = ""
			case "backspace":
				if len(m.cloneInput.input) > 0 {
					m.cloneInput.input = m.cloneInput.input[:len(m.cloneInput.input)-1]
				}
			case "ctrl+u":
				m.cloneInput.input = ""
			default:
				if len(msg.String()) == 1 {
					m.cloneInput.input += msg.String()
				}
			}
		}

	case stateCloning:
		if result, ok := msg.(cloneResult); ok {
			if result.err != nil {
				m.err = result.err
				m.state = stateCloneInput
			} else {
				m.err = fmt.Errorf("✓ Cloned to: %s", result.path)
				m.state = stateMenu
				m.input.input = ""
			}
		}

	case stateDashboard:
		if key, ok := msg.(tea.KeyMsg); ok {
			if key.String() == "." {
				if m.dashboard.data.Repo != nil {
					m.input.input = m.dashboard.data.Repo.FullName
					m.state = stateLoading
					cmds = append(cmds, m.analyzeRepo(m.input.input), TickProgressCmd())
					return m, tea.Batch(cmds...)
				}
			}
		}

		newDash, newCmd := m.dashboard.Update(msg)
		m.dashboard = newDash.(DashboardModel)
		cmds = append(cmds, newCmd)

		if m.dashboard.BackToMenu {
			m.state = stateMenu
			m.dashboard.BackToMenu = false
			m.input.input = ""
		}
	case stateTree:
		newTree, newCmd := m.tree.Update(msg)
		m.tree = newTree.(TreeModel)
		cmds = append(cmds, newCmd)

		if m.tree.Done {
			if m.tree.SelectedPath != "" {
				// Initialize file edit model
				repoName := m.input.input
				if m.dashboard.data.Repo != nil && m.dashboard.data.Repo.FullName != "" {
					repoName = m.dashboard.data.Repo.FullName
				}
				m.fileEdit = NewFileEditModel(m.tree.SelectedPath, repoName)

				// Check ownership
				isOwner := m.checkOwnership()
				m.fileEdit.SetOwnership(isOwner)

				m.state = stateFileEdit
			} else {
				m.state = stateDashboard
			}
			m.tree.Done = false
			m.tree.SelectedPath = ""
		}

	case stateFileEdit:
		newFileEdit, newCmd := m.fileEdit.Update(msg)
		m.fileEdit = newFileEdit.(FileEditModel)
		cmds = append(cmds, newCmd)

		if m.fileEdit.Done {
			m.state = stateTree
			m.fileEdit.Done = false
		}
	}

	return m, tea.Batch(cmds...)
}

// View renders the current UI state
func (m MainModel) View() string {
	switch m.state {
	case stateMenu:
		return m.menu.View()
	case stateInput:
		return m.input.View()
	case stateLoading:
		return m.loading.View(m.windowWidth, m.windowHeight)
	case stateCompareInput:
		return m.compareInput.View()
	case stateCompareLoading:
		return m.compareLoading.View(m.windowWidth, m.windowHeight)
	case stateCompareResult:
		return m.compareResult.View(m.windowWidth, m.windowHeight)
	case stateSettings:
		return m.settings.View()
	case stateHelp:
		return m.help.View(m.windowWidth, m.windowHeight)
	case stateHistory:
		return m.history.View()
	case stateFavorites:
		return m.favorites.View(m.windowWidth, m.windowHeight)
	case stateCloneInput:
		return m.cloneInput.View(m.windowWidth, m.windowHeight)
	case stateCloning:
		return m.cloning.View(m.windowWidth, m.windowHeight)
	case stateDashboard:
		return m.dashboard.View()
	case stateTree:
		return m.tree.View()
	case stateFileEdit:
		return m.fileEdit.View()
	}
	return ""
}

func (m MainModel) inputView() string {
	inputContent :=
		TitleStyle.Render("📥 ENTER REPOSITORY") + "\n\n" +
			InputStyle.Render("> "+m.input.input) + "\n\n" +
			SubtleStyle.Render("Format: owner/repo or GitHub URL  •  Press Enter to analyze")

	if m.err != nil {
		inputContent += "\n\n" + ErrorStyle.Render(fmt.Sprintf("Error: %v", m.err))
	}

	box := BoxStyle.Render(inputContent)

	if m.windowWidth == 0 {
		return box
	}

	return lipgloss.Place(
		m.windowWidth,
		m.windowHeight,
		lipgloss.Center,
		lipgloss.Center,
		box,
	)
}

// cloneResult is the result of a clone operation
type cloneResult struct {
	err  error
	path string
}

// cloneRepo clones a repository to the Desktop folder
func (m MainModel) cloneRepo(repoName string) tea.Cmd {
	return func() tea.Msg {
		parts := strings.Split(repoName, "/")
		if len(parts) != 2 {
			return cloneResult{err: fmt.Errorf("invalid repository URL: must be in owner/repo format or a valid GitHub URL")}
		}

		// Get Desktop path
		home, err := os.UserHomeDir()
		if err != nil {
			return cloneResult{err: err}
		}
		desktopPath := filepath.Join(home, "Desktop")
		clonePath := filepath.Join(desktopPath, parts[1])

		// Check if already exists
		if _, err := os.Stat(clonePath); err == nil {
			return cloneResult{err: fmt.Errorf("folder already exists: %s", clonePath)}
		}

		// Clone the repository
		repoURL := fmt.Sprintf("https://github.com/%s/%s.git", parts[0], parts[1])
		cmd := exec.Command("git", "clone", repoURL, clonePath)

		if err := cmd.Run(); err != nil {
			return cloneResult{err: fmt.Errorf("clone failed: %w", err)}
		}

		// Open file manager to show the cloned folder
		openFileManager(clonePath)

		return cloneResult{path: clonePath}
	}
}

func (m MainModel) analyzeRepo(repoName string) tea.Cmd {
	return func() tea.Msg {
		parts := strings.Split(repoName, "/")
		if len(parts) != 2 {
			return fmt.Errorf("invalid repository URL: must be in owner/repo format or a valid GitHub URL")
		}

		// Check cache first
		if m.cache != nil {
			if entry, found := m.cache.Get(repoName); found {
				// Unmarshal cached analysis
				var result AnalysisResult
				if err := json.Unmarshal(entry.Analysis, &result); err == nil {
					// Return cached result with status
					return CachedAnalysisResult{
						Result:   result,
						IsCached: true,
						CachedAt: entry.CachedAt,
					}
				}
			}
		}

		tracker := NewProgressTracker()

		// Stage 1: Fetch repository
		client := github.NewClient()
		repo, err := client.GetRepo(parts[0], parts[1])
		if err != nil {
			return err
		}
		tracker.NextStage()

		// Stage 2: Analyze commits
		commits, err := client.GetCommits(parts[0], parts[1], 365)
		if err != nil {
			return fmt.Errorf("failed to get commits: %w", err)
		}
		tracker.NextStage()

		// Stage 3: Analyze contributors
		contributors, err := client.GetContributorsWithAvatars(parts[0], parts[1], 15)
		if err != nil {
			return fmt.Errorf("failed to get contributors: %w", err)
		}
		tracker.NextStage()

		// Stage 4: Analyze languages
		languages, err := client.GetLanguages(parts[0], parts[1])
		if err != nil {
			return fmt.Errorf("failed to get languages: %w", err)
		}
		fileTree, err := client.GetFileTree(parts[0], parts[1], repo.DefaultBranch)
		if err != nil {
			return fmt.Errorf("failed to get file tree: %w", err)
		}
		tracker.NextStage()

		// Stage 5: Compute metrics
		score := analyzer.CalculateHealth(repo, commits)
		busFactor, busRisk := analyzer.BusFactor(contributors)
		maturityScore, maturityLevel := analyzer.RepoMaturityScore(repo, len(commits), len(contributors), false)

		// Stage 6: Analyze dependencies and contributor insights
		deps, _ := analyzer.AnalyzeDependencies(client, parts[0], parts[1], repo.DefaultBranch, fileTree)
		contributorInsights := analyzer.AnalyzeContributors(contributors)

		// Stage 7: Security vulnerability scan
		security, _ := analyzer.ScanDependencies(deps)
		tracker.NextStage()

		// Mark complete
		tracker.NextStage()
		commitsLast90Days := 0
		cutoff := time.Now().AddDate(0, 0, -90)

		for _, c := range commits {
			if c.Commit.Author.Date.After(cutoff) {
				commitsLast90Days++
			}
		}
		riskAlerts := analyzer.AnalyzeRiskAlerts(
			busFactor,
			score,
			commitsLast90Days,
			security != nil && security.CriticalCount > 0,
		)

		// Analyze Hotspots
		hotspots, _ := analyzer.AnalyzeHotspots(repo, commits, fileTree, client)

		// Generate quality dashboard
		qualityDashboard := analyzer.GenerateQualityDashboard(
			repo,
			commits,
			contributors,
			score,
			busFactor,
			maturityLevel,
			maturityScore,
			security,
			nil, // codeQuality - not implemented yet
			deps,
			hotspots,
		)

		result := AnalysisResult{
			Repo:                repo,
			Commits:             commits,
			Contributors:        contributors,
			FileTree:            fileTree,
			Languages:           languages,
			HealthScore:         score,
			BusFactor:           busFactor,
			BusRisk:             busRisk,
			MaturityScore:       maturityScore,
			MaturityLevel:       maturityLevel,
			Dependencies:        deps,
			ContributorInsights: contributorInsights,
			Security:            security,
			ContributorActivity: analyzer.AnalyzeContributorActivity(commits),
			RiskAlerts:          riskAlerts,
			QualityDashboard:    qualityDashboard,
		}

		// Save to cache
		if m.cache != nil {
			m.cache.Set(repoName, result)
		}

		return result
	}
}

func (m MainModel) checkOwnership() bool {
	client := github.NewClient()
	user, err := client.GetUser()
	if err != nil {
		return false // If we can't get user, assume not owner
	}

	expectedOwner := m.fileEdit.repoOwner
	return user.Login == expectedOwner
}

func (m MainModel) compareInputView() string {
	var currentInput string
	var prompt string

	if m.compareStep == 0 {
		prompt = "📥 ENTER FIRST REPOSITORY"
		currentInput = m.compareInput1
	} else {
		prompt = "📥 ENTER SECOND REPOSITORY"
		currentInput = m.compareInput2
	}

	inputContent := TitleStyle.Render(prompt) + "\n\n"

	if m.compareStep == 1 {
		inputContent += SubtleStyle.Render("First: "+m.compareInput1) + "\n\n"
	}

	inputContent += InputStyle.Render("> "+currentInput) + "\n\n"
	inputContent += SubtleStyle.Render("Format: owner/repo  •  Press Enter to continue  •  ESC to go back")

	if m.err != nil {
		inputContent += "\n\n" + ErrorStyle.Render(fmt.Sprintf("Error: %v", m.err))
	}

	box := BoxStyle.Render(inputContent)

	if m.windowWidth == 0 {
		return box
	}

	return lipgloss.Place(
		m.windowWidth,
		m.windowHeight,
		lipgloss.Center,
		lipgloss.Center,
		box,
	)
}

func (m MainModel) compareResultView() string {
	if m.compareResult.result == nil || m.compareResult.result.Repo1.Repo == nil || m.compareResult.result.Repo2.Repo == nil {
		return "No comparison data"
	}

	r1 := m.compareResult.result.Repo1
	r2 := m.compareResult.result.Repo2

	header := TitleStyle.Render(fmt.Sprintf("📊 Comparison: %s vs %s", r1.Repo.FullName, r2.Repo.FullName))

	// Check if repositories are identical
	if r1.Repo.Stars == r2.Repo.Stars &&
		r1.Repo.Forks == r2.Repo.Forks &&
		len(r1.Commits) == len(r2.Commits) &&
		len(r1.Contributors) == len(r2.Contributors) &&
		r1.BusFactor == r2.BusFactor &&
		r1.MaturityScore == r2.MaturityScore {

		noDiffBox := BoxStyle.Render("✅ No differences found between the two repositories.\nBoth repositories have identical metrics.")
		footer := SubtleStyle.Render("j: export JSON • m: export Markdown • q/ESC: back to menu")

		content := lipgloss.JoinVertical(
			lipgloss.Left,
			header,
			noDiffBox,
			footer,
		)

		if m.windowWidth == 0 {
			return content
		}

		return lipgloss.Place(
			m.windowWidth,
			m.windowHeight,
			lipgloss.Center,
			lipgloss.Center,
			content,
		)
	}

	// Build comparison table
	rows := []string{
		fmt.Sprintf("%-20s │ %-25s │ %-25s", "Metric", r1.Repo.FullName, r2.Repo.FullName),
		strings.Repeat("─", 75),
		fmt.Sprintf("%-20s │ %-25d │ %-25d", "⭐ Stars", r1.Repo.Stars, r2.Repo.Stars),
		fmt.Sprintf("%-20s │ %-25d │ %-25d", "🍴 Forks", r1.Repo.Forks, r2.Repo.Forks),
		fmt.Sprintf("%-20s │ %-25d │ %-25d", "📦 Commits (1y)", len(r1.Commits), len(r2.Commits)),
		fmt.Sprintf("%-20s │ %-25d │ %-25d", "👥 Contributors", len(r1.Contributors), len(r2.Contributors)),
		fmt.Sprintf("%-20s │ %-25s │ %-25s", "💚 Health Score", fmt.Sprintf("%d", r1.HealthScore), fmt.Sprintf("%d", r2.HealthScore)),
		fmt.Sprintf("%-20s │ %-25s │ %-25s", "⚠️ Bus Factor", fmt.Sprintf("%d (%s)", r1.BusFactor, r1.BusRisk), fmt.Sprintf("%d (%s)", r2.BusFactor, r2.BusRisk)),
		fmt.Sprintf("%-20s │ %-25s │ %-25s", "🏗️ Maturity", fmt.Sprintf("%s (%d)", r1.MaturityLevel, r1.MaturityScore), fmt.Sprintf("%s (%d)", r2.MaturityLevel, r2.MaturityScore)),
	}

	tableContent := strings.Join(rows, "\n")
	tableBox := BoxStyle.Render(tableContent)

	// Verdict
	var verdict string
	if r1.MaturityScore > r2.MaturityScore {
		verdict = fmt.Sprintf("➡️ %s appears more mature and stable.", r1.Repo.FullName)
	} else if r2.MaturityScore > r1.MaturityScore {
		verdict = fmt.Sprintf("➡️ %s appears more mature and stable.", r2.Repo.FullName)
	} else {
		verdict = "➡️ Both repositories are similarly mature."
	}
	verdictBox := BoxStyle.Render("📌 Verdict\n" + verdict)

	footer := SubtleStyle.Render("j: export JSON • m: export Markdown • q/ESC: back to menu")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		tableBox,
		verdictBox,
		footer,
	)

	// Add error/status message if present
	if m.err != nil {
		content = lipgloss.JoinVertical(
			lipgloss.Left,
			content,
			"\n"+ErrorStyle.Render(fmt.Sprintf("Status: %v", m.err)),
		)
	}

	if m.windowWidth == 0 {
		return content
	}

	return lipgloss.Place(
		m.windowWidth,
		m.windowHeight,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}

func (m MainModel) compareRepos(repo1Name, repo2Name string) tea.Cmd {
	return func() tea.Msg {
		parts1 := strings.Split(repo1Name, "/")
		parts2 := strings.Split(repo2Name, "/")

		if len(parts1) != 2 {
			return fmt.Errorf("invalid repository URL: first repository must be in owner/repo format or a valid GitHub URL")
		}
		if len(parts2) != 2 {
			return fmt.Errorf("invalid repository URL: second repository must be in owner/repo format or a valid GitHub URL")
		}

		client := github.NewClient()

		// Analyze first repo
		repo1, err := client.GetRepo(parts1[0], parts1[1])
		if err != nil {
			return fmt.Errorf("failed to fetch %s: %w", repo1Name, err)
		}
		commits1, _ := client.GetCommits(parts1[0], parts1[1], 365)
		contributors1, _ := client.GetContributorsWithAvatars(parts1[0], parts1[1], 15)
		languages1, _ := client.GetLanguages(parts1[0], parts1[1])
		fileTree1, _ := client.GetFileTree(parts1[0], parts1[1], repo1.DefaultBranch)
		score1 := analyzer.CalculateHealth(repo1, commits1)
		busFactor1, busRisk1 := analyzer.BusFactor(contributors1)
		maturityScore1, maturityLevel1 := analyzer.RepoMaturityScore(repo1, len(commits1), len(contributors1), false)

		result1 := AnalysisResult{
			Repo:          repo1,
			Commits:       commits1,
			Contributors:  contributors1,
			FileTree:      fileTree1,
			Languages:     languages1,
			HealthScore:   score1,
			BusFactor:     busFactor1,
			BusRisk:       busRisk1,
			MaturityScore: maturityScore1,
			MaturityLevel: maturityLevel1,
		}

		// Analyze second repo
		repo2, err := client.GetRepo(parts2[0], parts2[1])
		if err != nil {
			return fmt.Errorf("failed to fetch %s: %w", repo2Name, err)
		}
		commits2, _ := client.GetCommits(parts2[0], parts2[1], 365)
		contributors2, _ := client.GetContributorsWithAvatars(parts2[0], parts2[1], 15)
		languages2, _ := client.GetLanguages(parts2[0], parts2[1])
		fileTree2, _ := client.GetFileTree(parts2[0], parts2[1], repo2.DefaultBranch)
		score2 := analyzer.CalculateHealth(repo2, commits2)
		busFactor2, busRisk2 := analyzer.BusFactor(contributors2)
		maturityScore2, maturityLevel2 := analyzer.RepoMaturityScore(repo2, len(commits2), len(contributors2), false)

		result2 := AnalysisResult{
			Repo:          repo2,
			Commits:       commits2,
			Contributors:  contributors2,
			FileTree:      fileTree2,
			Languages:     languages2,
			HealthScore:   score2,
			BusFactor:     busFactor2,
			BusRisk:       busRisk2,
			MaturityScore: maturityScore2,
			MaturityLevel: maturityLevel2,
		}

		return CompareResult{
			Repo1: result1,
			Repo2: result2,
		}
	}
}

// Run starts the Bubble Tea program
func Run(cache *cache.Cache, config *config.AppSettings) error {
	model := NewMainModel(cache, config)
	p := tea.NewProgram(model, tea.WithAltScreen())
	_, err := p.Run()
	return err
}
