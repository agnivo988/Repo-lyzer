package ui

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/agnivo988/Repo-lyzer/internal/analyzer"
	"github.com/agnivo988/Repo-lyzer/internal/cache"
	"github.com/agnivo988/Repo-lyzer/internal/config"
	"github.com/agnivo988/Repo-lyzer/internal/github"
	"github.com/charmbracelet/bubbles/help"
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

type MainModel struct {
	state           sessionState
	menu            MenuModel
	input           string // Repository input
	compareInput1   string // First repo for comparison
	compareInput2   string // Second repo for comparison
	compareStep     int    // 0 = entering first repo, 1 = entering second repo
	spinner         spinner.Model
	dashboard       DashboardModel
	tree            TreeModel
	fileEdit        FileEditModel
	help            help.Model
	progress        *ProgressTracker
	animTick        int
	err             error
	windowWidth     int
	windowHeight    int
	analysisType    string // quick, detailed, custom
	appSettings     tea.LogOptionsSetter
	compareResult   *CompareResult // Holds comparison data
	history         *History       // Analysis history
	historyCursor   int            // Current selection in history
	helpContent     string         // Content for help screen
	settingsOption  string         // Selected settings option
	cache           *cache.Cache       // Offline cache for analysis results
	cacheStatus     string             // Cache status: "fresh", "cached", "expired", ""
	favorites       *Favorites         // Favorite repositories
	favoritesCursor int                // Current selection in favorites
	appConfig       *config.AppSettings // Application settings
	tokenInput      string             // Buffer for token input
	inTokenInput    bool               // Whether currently inputting token
}

// NewMainModel creates a new main application model with default settings.
// It initializes all sub-models (menu, dashboard, tree, etc.) and sets up
// the spinner with appropriate styling for the loading state.
// Returns the initialized MainModel with state set to menu.
func NewMainModel() MainModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	// Initialize cache
	repoCache, _ := cache.NewCache()

	// Load application settings
	appConfig, _ := config.LoadSettings()

	// Apply saved theme
	if appConfig != nil && appConfig.ThemeName != "" {
		SetThemeByName(appConfig.ThemeName)
	}

	return MainModel{
		state:       stateMenu,
		menu:        NewMenuModel(),
		spinner:     s,
		dashboard:   NewDashboardModel(),
		tree:        NewTreeModel(nil),
		cache:       repoCache,
		appConfig:   appConfig,
	}
}

func (m MainModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowWidth = msg.Width
		m.windowHeight = msg.Height
		// Handle terminal resize
		if m.windowWidth != msg.Width || m.windowHeight != msg.Height {
			// Adapt layout accordingly
			m.windowWidth = msg.Width
			m.windowHeight = msg.Height
		}
		// Propagate to children
		m.menu.Update(msg)
		m.dashboard.Update(msg)
		m.help.Update(msg)
		newTree, _ := m.tree.Update(msg)
		m.tree = newTree.(TreeModel)

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		// Global shortcuts
		if msg.String() == "q" && m.state == stateMenu {
			return m, tea.Quit
		}
		// Ctrl+H → Open History from anywhere
		if msg.String() == "ctrl+h" {
			m.state = stateHistory
			m.historyCursor = 0
			history, _ := LoadHistory()
			m.history = history
			return m, nil
		}
	
	case struct{}:
		if m.state == stateLoading || m.state == stateCompareLoading {
			m.animTick++
			return m, TickProgressCmd()
		}

	case string:
		if msg == "switch_to_tree" {
			m.state = stateTree
			// Update tree with current analysis data
			if m.dashboard.data.Repo != nil {
				m.tree = NewTreeModel(&m.dashboard.data)
				// Initialize tree with current window size
				var cmd tea.Cmd
				var tm tea.Model
				tm, cmd = m.tree.Update(tea.WindowSizeMsg{Width: m.windowWidth, Height: m.windowHeight})
				m.tree = tm.(TreeModel)
				cmds = append(cmds, cmd)
			}
		}
		if msg == "refresh_data" {
			// Re-analyze the current repo
			if m.dashboard.data.Repo != nil {
				m.state = stateLoading
				cmds = append(cmds, m.analyzeRepo(m.dashboard.data.Repo.FullName), TickProgressCmd()) // Add TickProgressCmd
			}
		}
		if msg == "add_to_favorites" {
			// Add current repo to favorites
			if m.dashboard.data.Repo != nil {
				if m.favorites == nil {
					m.favorites, _ = LoadFavorites()
				}
				m.favorites.Add(m.dashboard.data.Repo.FullName)
				m.favorites.Save()
				m.err = fmt.Errorf("⭐ Added to favorites: %s", m.dashboard.data.Repo.FullName)
			}
		}
	}

	switch m.state {
	case stateMenu:
		newMenu, newCmd := m.menu.Update(msg)
		m.menu = newMenu.(MenuModel)
		cmds = append(cmds, newCmd)

		if m.menu.Done {
			switch m.menu.SelectedOption {
			case 0: // Analyze
				if m.menu.submenuType == "analyze" {
					// Analysis type selection
					analysisTypes := []string{"quick", "detailed", "custom"}
					if m.menu.submenuCursor < len(analysisTypes) {
						m.analysisType = analysisTypes[m.menu.submenuCursor]
					}
					m.state = stateInput
				}
				m.menu.Done = false
			case 1: // Favorites
				m.state = stateFavorites
				m.favoritesCursor = 0
				favs, _ := LoadFavorites()
				m.favorites = favs
				m.menu.Done = false
			case 2: // Compare
				m.state = stateCompareInput
				m.compareStep = 0
				m.compareInput1 = ""
				m.compareInput2 = ""
				m.menu.Done = false
			case 3: // History
				m.state = stateHistory
				m.historyCursor = 0
				history, _ := LoadHistory()
				m.history = history
				m.menu.Done = false
			case 4: // Clone Repository
				m.state = stateCloneInput
				m.input = ""
				m.menu.Done = false
			case 5: // Settings
				if m.menu.submenuType == "settings" {
					// Settings option selection
					settingsOptions := []string{"theme", "cache", "export", "token", "reset"}
					if m.menu.submenuCursor < len(settingsOptions) {
						m.settingsOption = settingsOptions[m.menu.submenuCursor]
					}
					m.state = stateSettings
				}
				m.menu.Done = false
			case 6: // Help
				if m.menu.submenuType == "help" {
					// Help option selection
					helpOptions := []string{"shortcuts", "getting-started", "features", "troubleshooting"}
					if m.menu.submenuCursor < len(helpOptions) {
						m.helpContent = helpOptions[m.menu.submenuCursor]
					}
					m.state = stateHelp
				}
				m.menu.Done = false
			case 7: // Exit
				return m, tea.Quit
			}
		}

	case stateInput:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyEnter:
				cleanInput := sanitizeRepoInput(m.input)

				if cleanInput != "" {
					m.input = cleanInput
					m.err = nil
					m.state = stateLoading
					cmds = append(cmds, m.analyzeRepo(cleanInput), TickProgressCmd())
				} else {
					m.err = fmt.Errorf("please enter a valid repository (owner/repo or GitHub URL)")
				}

			case tea.KeyBackspace:
				if len(m.input) > 0 {
					m.input = m.input[:len(m.input)-1]
				}
			case tea.KeyRunes:
				m.input += string(msg.Runes)
			case tea.KeyEsc:
				m.state = stateMenu
			case tea.KeyCtrlU:
				m.input = "" // Clear entire line
			case tea.KeyCtrlA:
				// Move to start - for TUI we just clear (no cursor)
				// In a real implementation, you'd track cursor position
			case tea.KeyCtrlE:
				// Move to end - already at end in this simple impl
			case tea.KeyCtrlW:
				// Delete word backward
				m.input = strings.TrimRight(m.input, " ")
				if idx := strings.LastIndex(m.input, " "); idx >= 0 {
					m.input = m.input[:idx+1]
				} else {
					m.input = ""
				}
			}
		}

	case stateCompareInput:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyEnter:
				if m.compareStep == 0 && m.compareInput1 != "" {
					// Sanitize first repo
					m.compareInput1 = sanitizeRepoInput(m.compareInput1)
					m.compareStep = 1

				} else if m.compareStep == 1 && m.compareInput2 != "" {
					// Sanitize both repos before comparison
					m.compareInput1 = sanitizeRepoInput(m.compareInput1)
					m.compareInput2 = sanitizeRepoInput(m.compareInput2)

					m.err = nil
					m.state = stateCompareLoading
					cmds = append(cmds, m.compareRepos(m.compareInput1, m.compareInput2), TickProgressCmd())
				}

			case tea.KeyBackspace:
				if m.compareStep == 0 && len(m.compareInput1) > 0 {
					m.compareInput1 = m.compareInput1[:len(m.compareInput1)-1]
				} else if m.compareStep == 1 && len(m.compareInput2) > 0 {
					m.compareInput2 = m.compareInput2[:len(m.compareInput2)-1]
				}
			case tea.KeyRunes:
				if m.compareStep == 0 {
					m.compareInput1 += string(msg.Runes)
				} else {
					m.compareInput2 += string(msg.Runes)
				}
			case tea.KeyEsc:
				if m.compareStep == 1 {
					// Go back to first repo input
					m.compareStep = 0
				} else {
					m.state = stateMenu
					m.menu.Done = false
					m.compareInput1 = ""
					m.compareInput2 = ""
				}
			case tea.KeyCtrlU:
				// Clear current input
				if m.compareStep == 0 {
					m.compareInput1 = ""
				} else {
					m.compareInput2 = ""
				}
			case tea.KeyCtrlW:
				// Delete word backward
				if m.compareStep == 0 {
					m.compareInput1 = strings.TrimRight(m.compareInput1, " ")
					if idx := strings.LastIndex(m.compareInput1, " "); idx >= 0 {
						m.compareInput1 = m.compareInput1[:idx+1]
					} else {
						m.compareInput1 = ""
					}
				} else {
					m.compareInput2 = strings.TrimRight(m.compareInput2, " ")
					if idx := strings.LastIndex(m.compareInput2, " "); idx >= 0 {
						m.compareInput2 = m.compareInput2[:idx+1]
					} else {
						m.compareInput2 = ""
					}
				}
			}
		}

	case stateCompareLoading:
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)

		switch msg := msg.(type) {
		case CompareResult:
			m.compareResult = &msg
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
				m.compareResult = nil
				m.compareInput1 = ""
				m.compareInput2 = ""
			case "j":
				// Export comparison to JSON
				if m.compareResult != nil {
					filename, err := ExportCompareJSON(*m.compareResult)
					if err != nil {
						m.err = err
					} else {
						m.err = nil
						// Show success message briefly (will need status in model)
						_ = filename // TODO: show status message
					}
				}
			case "m":
				// Export comparison to Markdown
				if m.compareResult != nil {
					filename, err := ExportCompareMarkdown(*m.compareResult)
					if err != nil {
						m.err = err
					} else {
						m.err = nil
						_ = filename
					}
				}
			}
		}

	case stateLoading:
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)

		if result, ok := msg.(AnalysisResult); ok {
			m.dashboard.SetData(result)
			m.dashboard.SetCacheStatus("fresh")
			m.state = stateDashboard
			m.progress = nil
			m.cacheStatus = "fresh"
			// Save to history
			if m.history == nil {
				m.history, _ = LoadHistory()
			}
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
			if m.history == nil {
				m.history, _ = LoadHistory()
			}
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
		case tea.KeyMsg:
			switch msg.String() {
			case "up", "k":
				if m.favoritesCursor > 0 {
					m.favoritesCursor--
				}
			case "down", "j":
				if m.favorites != nil && m.favoritesCursor < len(m.favorites.Items)-1 {
					m.favoritesCursor++
				}
			case "enter":
				// Analyze selected favorite
				if m.favorites != nil && len(m.favorites.Items) > 0 {
					repoName := m.favorites.Items[m.favoritesCursor].RepoName
					m.favorites.UpdateUsage(repoName)
					m.favorites.Save()
					m.input = repoName
					m.state = stateLoading
					cmds = append(cmds, m.analyzeRepo(repoName), TickProgressCmd())
				}
			case "d":
				// Remove from favorites
				if m.favorites != nil && len(m.favorites.Items) > 0 {
					m.favorites.Remove(m.favorites.Items[m.favoritesCursor].RepoName)
					m.favorites.Save()
					if m.favoritesCursor >= len(m.favorites.Items) && m.favoritesCursor > 0 {
						m.favoritesCursor--
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
				if m.history != nil && m.historyCursor < len(m.history.Entries)-1 {
					m.historyCursor++
				}
			case "enter":
				// Re-analyze selected repo
				if m.history != nil && len(m.history.Entries) > 0 {
					repoName := m.history.Entries[m.historyCursor].RepoName
					m.input = repoName
					m.state = stateLoading
					cmds = append(cmds, m.analyzeRepo(repoName), TickProgressCmd())
				}
			case "d":
				// Delete selected entry
				if m.history != nil && len(m.history.Entries) > 0 {
					m.history.Delete(m.historyCursor)
					m.history.Save()
					if m.historyCursor >= len(m.history.Entries) && m.historyCursor > 0 {
						m.historyCursor--
					}
				}
			case "c":
				// Clear all history
				if m.history != nil {
					m.history.Clear()
					m.history.Save()
					m.historyCursor = 0
				}
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
				if m.input != "" {
					m.state = stateCloning
					cmds = append(cmds, m.cloneRepo(m.input))
				}
			case "esc":
				m.state = stateMenu
				m.input = ""
			case "backspace":
				if len(m.input) > 0 {
					m.input = m.input[:len(m.input)-1]
				}
			case "ctrl+u":
				m.input = ""
			default:
				if len(msg.String()) == 1 {
					m.input += msg.String()
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
				m.input = ""
			}
		}

	case stateDashboard:
		if key, ok := msg.(tea.KeyMsg); ok {
			if key.String() == "." {
				if m.dashboard.data.Repo != nil {
					m.input = m.dashboard.data.Repo.FullName
					m.state = stateLoading
					cmds = append(cmds, m.analyzeRepo(m.input), TickProgressCmd())
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
			m.input = ""
		}
	case stateTree:
		newTree, newCmd := m.tree.Update(msg)
		m.tree = newTree.(TreeModel)
		cmds = append(cmds, newCmd)

		if m.tree.Done {
			if m.tree.SelectedPath != "" {
				// Initialize file edit model
				repoName := m.input
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

func (m MainModel) View() string {
	switch m.state {
	case stateMenu:
		return m.menu.View()
	case stateInput:
		return m.inputView()
	case stateCompareInput:
		return m.compareInputView()
	case stateFavorites:
		return m.favoritesView()
	case stateHistory:
		return m.historyView()
	case stateCloneInput:
		return m.cloneInputView()
	case stateCloning:
		return m.cloningView()
	case stateLoading:
		loadMsg := fmt.Sprintf("📊 Analyzing %s", m.input)
		if m.analysisType != "" {
			loadMsg += fmt.Sprintf(" (%s mode)", strings.ToUpper(m.analysisType))
		}

		statusView := fmt.Sprintf("%s %s...", m.spinner.View(), loadMsg)
		
		if len(SatelliteFrames) > 0 {
			frame := SatelliteFrames[m.animTick%len(SatelliteFrames)]
			statusView += "\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("#00E5FF")).Render(frame)
		}

		// Show progress stages if available
		if m.progress != nil {
			stages := m.progress.GetAllStages()
			statusView += "\n\n"
			for _, stage := range stages {
				prefix := "⏳ "
				if stage.IsComplete {
					prefix = "✅ "
				} else if stage.IsActive {
					prefix = "⚙️  "
				}
				statusView += prefix + stage.Name + "\n"
			}

			// Add elapsed time
			elapsed := m.progress.GetElapsedTime()
			statusView += fmt.Sprintf("\n⏱️  %ds elapsed", int(elapsed.Seconds()))
		}

		statusView += "\n\n" + SubtleStyle.Render("Press ESC to cancel")

		return lipgloss.Place(
			m.windowWidth, m.windowHeight,
			lipgloss.Center, lipgloss.Center,
			statusView,
		)
	case stateCompareLoading:
		loadMsg := fmt.Sprintf("📊 Comparing %s vs %s", m.compareInput1, m.compareInput2)
		statusView := fmt.Sprintf("%s %s...", m.spinner.View(), loadMsg)
		
		if len(SatelliteFrames) > 0 {
			frame := SatelliteFrames[m.animTick%len(SatelliteFrames)]
			statusView += "\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("#00E5FF")).Render(frame)
		}
		
		statusView += "\n\n" + SubtleStyle.Render("Press ESC to cancel")

		return lipgloss.Place(
			m.windowWidth, m.windowHeight,
			lipgloss.Center, lipgloss.Center,
			statusView,
		)
	case stateCompareResult:
		return m.compareResultView()
	case stateTree:
		return m.tree.View()
	case stateFileEdit:
		return m.fileEdit.View()
	case stateHelp:
		return m.helpView()
	case stateSettings:
		return m.settingsView()
	case stateDashboard:
		return m.dashboard.View()
	}
	return ""
}

func (m MainModel) inputView() string {
	inputContent :=
		TitleStyle.Render("📥 ENTER REPOSITORY") + "\n\n" +
			InputStyle.Render("> "+m.input) + "\n\n" +
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
			return cloneResult{err: fmt.Errorf("repository must be in owner/repo format")}
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
			return fmt.Errorf("repository must be in owner/repo format")
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
		contributors, err := client.GetContributors(parts[0], parts[1])
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
	if m.compareResult == nil || m.compareResult.Repo1.Repo == nil || m.compareResult.Repo2.Repo == nil {
		return "No comparison data"
	}

	r1 := m.compareResult.Repo1
	r2 := m.compareResult.Repo2

	header := TitleStyle.Render(fmt.Sprintf("📊 Comparison: %s vs %s", r1.Repo.FullName, r2.Repo.FullName))

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
			return fmt.Errorf("first repository must be in owner/repo format")
		}
		if len(parts2) != 2 {
			return fmt.Errorf("second repository must be in owner/repo format")
		}

		client := github.NewClient()

		// Analyze first repo
		repo1, err := client.GetRepo(parts1[0], parts1[1])
		if err != nil {
			return fmt.Errorf("failed to fetch %s: %w", repo1Name, err)
		}
		commits1, _ := client.GetCommits(parts1[0], parts1[1], 365)
		contributors1, _ := client.GetContributors(parts1[0], parts1[1])
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
		contributors2, _ := client.GetContributors(parts2[0], parts2[1])
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

func Run() error {
	p := tea.NewProgram(NewMainModel(), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
func sanitizeRepoInput(input string) string {
	// Remove null bytes and trim spaces
	clean := strings.ReplaceAll(input, "\x00", "")
	clean = strings.TrimSpace(clean)

	// Allow full GitHub URLs
	if strings.Contains(clean, "github.com/") {
		parts := strings.Split(clean, "github.com/")
		if len(parts) == 2 {
			clean = parts[1]
		}
	}

	// Remove trailing slash if present
	clean = strings.TrimSuffix(clean, "/")

	return clean
}

func (m MainModel) favoritesView() string {
	header := TitleStyle.Render("⭐ Favorite Repositories")

	if m.favorites == nil || len(m.favorites.Items) == 0 {
		content := lipgloss.JoinVertical(
			lipgloss.Left,
			header,
			BoxStyle.Render("No favorites yet!\n\nAnalyze a repository and press 'b' to bookmark it."),
			SubtleStyle.Render("a: add new • q/ESC: back to menu"),
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

	// Build favorites list
	var lines []string
	lines = append(lines, fmt.Sprintf("%-35s │ %-10s │ %s", "Repository", "Uses", "Last Used"))
	lines = append(lines, strings.Repeat("─", 65))

	for i, fav := range m.favorites.Items {
		prefix := "  "
		if i == m.favoritesCursor {
			prefix = "▶ "
		}
		line := fmt.Sprintf("%s%-33s │ %-10d │ %s",
			prefix,
			fav.RepoName,
			fav.UseCount,
			fav.LastUsed.Format("2006-01-02"),
		)
		if i == m.favoritesCursor {
			lines = append(lines, SelectedStyle.Render(line))
		} else {
			lines = append(lines, line)
		}
	}

	tableBox := BoxStyle.Render(strings.Join(lines, "\n"))
	footer := SubtleStyle.Render("↑↓: navigate • Enter: analyze • d: remove • a: add new • q/ESC: back")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		tableBox,
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

func (m MainModel) historyView() string {
	header := TitleStyle.Render("📜 Analysis History")

	if m.history == nil || len(m.history.Entries) == 0 {
		content := lipgloss.JoinVertical(
			lipgloss.Left,
			header,
			BoxStyle.Render("No history yet. Analyze a repository to get started!"),
			SubtleStyle.Render("q/ESC: back to menu"),
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

	// Build history list
	var lines []string
	lines = append(lines, fmt.Sprintf("%-30s │ %-8s │ %-5s │ %-12s │ %s", "Repository", "Stars", "Health", "Maturity", "Analyzed"))
	lines = append(lines, strings.Repeat("─", 85))

	for i, entry := range m.history.Entries {
		prefix := "  "
		if i == m.historyCursor {
			prefix = "▶ "
		}
		line := fmt.Sprintf("%s%-28s │ ⭐%-6d │ 💚%-3d │ %-12s │ %s",
			prefix,
			entry.RepoName,
			entry.Stars,
			entry.HealthScore,
			entry.MaturityLevel,
			entry.AnalyzedAt.Format("2006-01-02 15:04"),
		)
		if i == m.historyCursor {
			lines = append(lines, SelectedStyle.Render(line))
		} else {
			lines = append(lines, line)
		}
	}

	tableBox := BoxStyle.Render(strings.Join(lines, "\n"))

	footer := SubtleStyle.Render("↑↓: navigate • Enter: re-analyze • d: delete • c: clear all • q/ESC: back")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		tableBox,
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

func (m MainModel) cloneInputView() string {
	header := TitleStyle.Render("📥 CLONE REPOSITORY")

	inputContent := fmt.Sprintf(
		"Enter repository to clone (owner/repo):\n\n> %s█\n\n"+
			"The repository will be cloned to your Desktop folder.",
		m.input,
	)

	var errMsg string
	if m.err != nil {
		errMsg = "\n" + ErrorStyle.Render(m.err.Error())
	}

	footer := SubtleStyle.Render("Enter: clone • ESC: back • Ctrl+U: clear")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		BoxStyle.Render(inputContent),
		errMsg,
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

func (m MainModel) cloningView() string {
	header := TitleStyle.Render("📥 CLONING REPOSITORY")

	content := fmt.Sprintf(
		"%s Cloning %s to Desktop...\n\n"+
			"Please wait while the repository is being cloned.",
		m.spinner.View(),
		m.input,
	)

	return lipgloss.Place(
		m.windowWidth,
		m.windowHeight,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Left,
			header,
			BoxStyle.Render(content),
		),
	)
}

func (m MainModel) helpView() string {
	var title string
	var content string

	switch m.helpContent {
	case "shortcuts":
		title = "❓ Keyboard Shortcuts"
		content = `
Global:
  Ctrl+H        Open analysis history
Main Menu:
  ↑↓/jk         Navigate menu
  Enter         Select option
  q             Quit application

Repository Input:
  Enter         Start analysis
  ESC           Back to menu
  Ctrl+U        Clear input
  Ctrl+W        Delete word
  Ctrl+A        Move to start
  Ctrl+E        Move to end

Dashboard Navigation:
  ←→/hl         Switch between views
  1-7           Jump to specific view
  e             Toggle export menu
  f             Open file tree
  r             Refresh data
  ?/h           Toggle help
  q/ESC         Go back
   .            Re-analyze current repository

File Tree:
  ↑↓/jk         Navigate files
  Enter         Open file details
  ESC           Back to dashboard

History:
  ↑↓/jk         Navigate entries
  Enter         Re-analyze repository
  d             Delete entry
  c             Clear all history
  q/ESC         Back to menu
`
	case "getting-started":
		title = "🚀 Getting Started"
		content = `
Welcome to Repo-lyzer!

1. Choose "Analyze Repository" from the main menu
2. Enter a repository in the format: owner/repo
   Example: microsoft/vscode
3. Select analysis type:
   - Quick: Fast overview
   - Detailed: Comprehensive analysis
   - Custom: Advanced options
4. Wait for analysis to complete
5. Navigate through the dashboard views
6. Export results if needed

For GitHub API access:
- Set GITHUB_TOKEN environment variable for higher rate limits
- Private repositories require authentication
`
	case "features":
		title = "✨ Features Guide"
		content = `
Repository Analysis:
  • Health Score: Overall repository health
  • Bus Factor: Risk of losing key contributors
  • Maturity Level: Project maturity assessment
  • Language Breakdown: Programming languages used
  • Commit Activity: Development activity over time
  • Top Contributors: Most active contributors
  • Recruiter Summary: Key insights for hiring

Export Options:
  • JSON: Structured data for further processing
  • Markdown: Human-readable reports

Additional Features:
  • Repository Comparison: Compare multiple repos
  • Analysis History: Re-analyze previous repos
  • File Tree: Explore repository structure
  • GitHub API Status: Monitor rate limit usage
`
	case "troubleshooting":
		title = "🔧 Troubleshooting"
		content = `
Common Issues:

Repository Not Found:
  • Check spelling: owner/repo format
  • Ensure repository is public or you have access
  • GitHub API might be rate limited

Analysis Fails:
  • Check internet connection
  • Verify GitHub API status
  • Try again later if rate limited

High Rate Limits:
  • Set GITHUB_TOKEN environment variable
  • Authenticated requests: 5000/hour
  • Unauthenticated: 60/hour

Private Repositories:
  • Require GITHUB_TOKEN with repo scope
  • Token must have access to the repository

Performance:
  • Detailed analysis takes longer
  • Large repositories may take several minutes
  • Use Quick analysis for fast results
`
	default:
		title = "❓ Help"
		content = `
Select a help topic from the menu above.
`
	}

	helpContent := TitleStyle.Render(title) + "\n\n" + content + "\n\n" + SubtleStyle.Render("Press ESC or q to go back")

	box := BoxStyle.Render(helpContent)

	if m.windowWidth == 0 {
		return box
	}

	return lipgloss.Place(
		m.windowWidth, m.windowHeight,
		lipgloss.Center, lipgloss.Center,
		box,
	)
}

func (m MainModel) settingsView() string {
	var title string
	var content string

	switch m.settingsOption {
	case "theme":
		title = "🎨 Theme Settings"

		// Build theme list with current indicator
		themeList := ""
		for i, theme := range AvailableThemes {
			indicator := "  "
			if i == CurrentThemeIndex {
				indicator = "▶ "
			}
			themeList += fmt.Sprintf("  %s[%d] %s\n", indicator, i+1, theme.Name)
		}

		content = fmt.Sprintf(`
Current theme: %s

Available themes:
%s
Keybindings:
  • Press 1-7 to select a theme
  • Press 't' to cycle through themes

Theme changes are applied immediately!
`, CurrentTheme.Name, themeList)
	case "cache":
		title = "💾 Cache Settings"

		// Get cache stats
		cacheInfo := "Cache not initialized"
		if m.cache != nil {
			stats := m.cache.GetStats()
			cfg := m.cache.GetConfig()

			enabledStr := "Disabled"
			if cfg.Enabled {
				enabledStr = "Enabled"
			}
			autoStr := "Off"
			if cfg.AutoCache {
				autoStr = "On"
			}

			cacheInfo = fmt.Sprintf(`
Status: %s
Auto-cache: %s
TTL: %s
Max Size: %d MB

Statistics:
  • Total repos cached: %d
  • Valid (not expired): %d
  • Expired: %d
  • Cache size: %.2f MB
  • Location: %s

Keybindings:
  • Press 'e' to toggle caching
  • Press 'a' to toggle auto-cache
  • Press 'c' to clear all cache
  • Press 'x' to clean expired entries
`, enabledStr, autoStr, cache.FormatTTL(cfg.TTL), cfg.MaxSize,
				stats.TotalRepos, stats.ValidRepos, stats.ExpiredRepos,
				stats.TotalSizeMB, stats.CacheDir)
		}
		content = cacheInfo
	case "export":
		title = "📤 Export Options"

		// Get current export settings
		currentFormat := "JSON"
		exportDir := "~/Downloads/"
		if m.appConfig != nil {
			currentFormat = m.appConfig.DefaultExportFormat.DisplayName()
			exportDir = m.appConfig.ExportDirectory
		}

		// Build format list with indicator
		formatList := ""
		formats := []string{"JSON", "Markdown", "CSV", "HTML"}
		for _, f := range formats {
			indicator := "  "
			if f == currentFormat {
				indicator = "▶ "
			}
			formatList += fmt.Sprintf("  %s%s\n", indicator, f)
		}

		content = fmt.Sprintf(`
Current export format: %s
Export directory: %s

Available formats:
%s
Keybindings:
  • Press 'f' to cycle through formats

Export formats are saved automatically!
`, currentFormat, exportDir, formatList)
	case "token":
		title = "🔑 GitHub Token"

		// Check if in token input mode
		if m.inTokenInput {
			content = fmt.Sprintf(`
Enter GitHub Personal Access Token:

> %s█

Press Enter to save, ESC to cancel.
`, m.tokenInput)
		} else {
			// Show current token status
			tokenStatus := "❌ Not configured"
			tokenDisplay := ""
			if m.appConfig != nil && m.appConfig.HasGitHubToken() {
				tokenStatus = "✅ Configured"
				tokenDisplay = fmt.Sprintf("\nToken: %s", m.appConfig.GetMaskedToken())
			}

			// Check environment variable
			envToken := os.Getenv("GITHUB_TOKEN")
			envStatus := "Not set"
			if envToken != "" {
				envStatus = "Set (will be used if app token not configured)"
			}

			content = fmt.Sprintf(`
GitHub API Token Configuration:

Status: %s%s
Environment: %s

Keyindings:
  • Press 'i' to input a new token
  • Press 'c' to clear saved token

Benefits of using a token:
  • Higher API rate limits (5000 vs 60 requests/hour)
  • Access to private repositories
  • More detailed analysis
`, tokenStatus, tokenDisplay, envStatus)
		}
	case "reset":
		title = "🔄 Reset to Defaults"
		content = `
Reset all settings to default values:

This will:
  • Clear all saved settings
  • Reset theme to default
  • Clear export preferences
  • Remove custom configurations

Warning: This action cannot be undone.

Press 'y' to confirm reset, or ESC to cancel.
`
	default:
		title = "⚙️ Settings"
		content = `
Select a settings option from the menu.
`
	}

	settingsContent := TitleStyle.Render(title) + "\n\n" + content + "\n\n" + SubtleStyle.Render("Press ESC or q to go back")

	box := BoxStyle.Render(settingsContent)

	if m.windowWidth == 0 {
		return box
	}

	return lipgloss.Place(
		m.windowWidth, m.windowHeight,
		lipgloss.Center, lipgloss.Center,
		box,
	)
}
