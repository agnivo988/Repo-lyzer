package ui

import (
	"fmt"
	"strings"

	"github.com/agnivo988/Repo-lyzer/internal/analyzer"
	"github.com/agnivo988/Repo-lyzer/internal/github"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type sessionState int

const (
	stateMenu sessionState = iota
	stateInput
	stateLoading
	stateDashboard
	stateTree
	stateSettings
	stateHelp
	stateHistory
	stateCompareInput
	stateCompareLoading
	stateCompareResult
)

type MainModel struct {
	state         sessionState
	menu          MenuModel
	input         string // Repository input
	compareInput1 string // First repo for comparison
	compareInput2 string // Second repo for comparison
	compareStep   int    // 0 = entering first repo, 1 = entering second repo
	spinner       spinner.Model
	dashboard     DashboardModel
	tree          TreeModel
	help          help.Model
	progress      *ProgressTracker
	err           error
	windowWidth   int
	windowHeight  int
	analysisType  string // quick, detailed, custom
	appSettings    tea.LogOptionsSetter
	compareResult *CompareResult // Holds comparison data
	history       *History       // Analysis history
	historyCursor int            // Current selection in history
}

func NewMainModel() MainModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return MainModel{
		state:        stateMenu,
		menu:         NewMenuModel(),
		spinner:      s,
		dashboard:    NewDashboardModel(),
		tree:         NewTreeModel(nil),
		appSettings:  nil, 
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
				cmds = append(cmds, m.analyzeRepo(m.dashboard.data.Repo.FullName))
			}
		}
	}

	switch m.state {
	case stateMenu:
		newMenu, newCmd := m.menu.Update(msg)
		m.menu = newMenu.(MenuModel)
		cmds = append(cmds, newCmd)

		if m.menu.SelectedOption == 0 && m.menu.Done { // Analyze
			m.state = stateInput
			m.menu.Done = false // Reset for back navigation
		} else if m.menu.SelectedOption == 1 && m.menu.Done { // Compare
			m.state = stateCompareInput
			m.compareStep = 0
			m.compareInput1 = ""
			m.compareInput2 = ""
			m.menu.Done = false
		} else if m.menu.SelectedOption == 2 && m.menu.Done { // History
			m.state = stateHistory
			m.historyCursor = 0
			// Load history
			history, _ := LoadHistory()
			m.history = history
			m.menu.Done = false
		} else if m.menu.SelectedOption == 3 && m.menu.Done { // Exit
			return m, tea.Quit
		}

	case stateInput:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyEnter:
				if m.input != "" {
					m.state = stateLoading
					cmds = append(cmds, m.analyzeRepo(m.input))
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
					// Move to second repo input
					m.compareStep = 1
				} else if m.compareStep == 1 && m.compareInput2 != "" {
					// Both repos entered, start comparison
					m.state = stateCompareLoading
					cmds = append(cmds, m.compareRepos(m.compareInput1, m.compareInput2))
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
			m.state = stateDashboard
			m.progress = nil
			// Save to history
			if m.history == nil {
				m.history, _ = LoadHistory()
			}
			m.history.AddEntry(result)
			m.history.Save()
		}
		if err, ok := msg.(error); ok {
			m.err = err
			m.state = stateInput // Go back to input on error
			m.progress = nil
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
					cmds = append(cmds, m.analyzeRepo(repoName))
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

	case stateDashboard:
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
			m.state = stateDashboard
			m.tree.Done = false
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
	case stateHistory:
		return m.historyView()
	case stateLoading:
		loadMsg := fmt.Sprintf("📊 Analyzing %s", m.input)
		if m.analysisType != "" {
			loadMsg += fmt.Sprintf(" (%s mode)", strings.ToUpper(m.analysisType))
		}

		statusView := fmt.Sprintf("%s %s...", m.spinner.View(), loadMsg)

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
	case stateDashboard:
		return m.dashboard.View()
	}
	return ""
}

func (m MainModel) inputView() string {
	inputContent :=
		TitleStyle.Render("📥 ENTER REPOSITORY") + "\n\n" +
			InputStyle.Render("> "+m.input) + "\n\n" +
			SubtleStyle.Render("Format: owner/repo  •  Press Enter to run")

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

func (m MainModel) analyzeRepo(repoName string) tea.Cmd {
	return func() tea.Msg {
		parts := strings.Split(repoName, "/")
		if len(parts) != 2 {
			return fmt.Errorf("repository must be in owner/repo format")
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
		commits, _ := client.GetCommits(parts[0], parts[1], 365)
		tracker.NextStage()

		// Stage 3: Analyze contributors
		contributors, _ := client.GetContributors(parts[0], parts[1])
		tracker.NextStage()

		// Stage 4: Analyze languages
		languages, _ := client.GetLanguages(parts[0], parts[1])
		fileTree, _ := client.GetFileTree(parts[0], parts[1], repo.DefaultBranch)
		tracker.NextStage()

		// Stage 5: Compute metrics
		score := analyzer.CalculateHealth(repo, commits)
		busFactor, busRisk := analyzer.BusFactor(contributors)
		maturityScore, maturityLevel := analyzer.RepoMaturityScore(repo, len(commits), len(contributors), false)
		tracker.NextStage()

		// Mark complete
		tracker.NextStage()

		return AnalysisResult{
			Repo:          repo,
			Commits:       commits,
			Contributors:  contributors,
			FileTree:      fileTree,
			Languages:     languages,
			HealthScore:   score,
			BusFactor:     busFactor,
			BusRisk:       busRisk,
			MaturityScore: maturityScore,
			MaturityLevel: maturityLevel,
		}
	}
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
